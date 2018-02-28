package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	// check for required layout and home files
	if _, err := os.Stat("content/home.gohtml"); os.IsNotExist(err) {
		log.Fatalln("home.gohtml does not exist")
	}
	if _, err := os.Stat("content/layout.gohtml"); os.IsNotExist(err) {
		log.Fatalln("layout.gohtml does not exist")
	}

	// Basic non-rotating log
	os.MkdirAll("log", os.ModePerm)
	f, err := os.OpenFile(filepath.Join("log", "log.txt"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stderr, f))

	// if not dev mode, parse templates during program load (otherwise, parse at each request)
	var dev = len(os.Args) > 1 && strings.ToUpper(os.Args[1]) == "DEV"
	var pages Pages //TODO: Race condition exists in DEV mode; eliminate this with a mutex
	if !dev {
		var err error
		pages, err = parsePages(dev)
		if err != nil {
			log.Fatalln("problem during production startup:", err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if dev {
			var err error
			pages, err = parsePages(dev)
			if err != nil {
				log.Println("problem during development request:", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		app := App{
			Brand: getBrand(),
			Nav:   buildNav(pages),
		}
		// check for optional components and as template.HTML if necessary
		// TODO: If I add optional components, split this into a func that takes a pointer to template.HTML
		// and sets it based upon the component name
		bannerPath := filepath.Join("content", "banner.gohtml")
		if _, err := os.Stat("content/banner.gohtml"); err == nil {
			var buf bytes.Buffer
			tmpl, err := template.ParseFiles(bannerPath)
			if err != nil {
				log.Println("could not parse banner:", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			err = tmpl.ExecuteTemplate(&buf, "banner.gohtml", app)
			if err != nil {
				log.Println("could not execute banner:", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			html := template.HTML(buf.String())
			app.Banner = &html
		}

		switch {
		case r.URL.Path == "/":
			pages["home"].Execute(w, app)
		case isContentFile(r.URL):
			fs := http.FileServer(http.Dir("content"))
			fs.ServeHTTP(w, r)
		default:
			if val, ok := pages[lastPart(r.URL)]; ok {
				val.Execute(w, app)
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	})
	srv := &http.Server{
		Addr:    ":3035",
		Handler: mux,
		//time from when the connection is accepted to when the request body is fully read
		ReadTimeout: 5 * time.Second,
		//time from the end of the request header read to the end of the response write
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Now serving on port 3035")
	log.Println(srv.ListenAndServe())
}

func lastPart(url *url.URL) string {
	pathParts := strings.Split(url.Path, `/`)
	return pathParts[len(pathParts)-1]
}

func isContentFile(url *url.URL) bool {
	part := lastPart(url)
	// Does last portion of URL contain a '.'? Make sure it isn't gohtml?
	return strings.Contains(part, ".") && !strings.Contains(part, ".gohtml")
}

func notComponent(filename string) bool {
	var componentPaths = []string{
		filepath.Join("content", "layout.gohtml"),
		filepath.Join("content", "banner.gohtml"),
	}
	for _, f := range componentPaths {
		if f == filename {
			return false
		}
	}
	return true
}

func parsePages(isDev bool) (Pages, error) {
	pages := make(map[string]*template.Template)
	pageFiles, err := filepath.Glob("content/*.gohtml")
	if err != nil {
		return pages, fmt.Errorf("could not parse pages: %s", err)
	}
	for _, f := range pageFiles {
		if notComponent(f) {
			tmpl, err := template.ParseFiles("content/layout.gohtml", f)
			if err != nil {
				return pages, fmt.Errorf("could not parse template: %s", err)
			}
			pageName := strings.Split(f, ".")[0]
			pages[filepath.Base(pageName)] = tmpl
		}
	}
	return pages, nil
}

// Pages is a map of kebabcase names to templates
type Pages map[string]*template.Template

// App specifies Brand and Nav
type App struct {
	Brand  string
	Banner *template.HTML
	Nav    Nav
}

// NavItem specifies the text and URL for a link on the top nav bar
type NavItem struct {
	Text string
	URL  template.URL
}

// Nav is a slice of NavItem
type Nav []NavItem

func buildNav(pages Pages) Nav {
	menu := make([]NavItem, 0, len(pages))
	// add Home nav item first
	homeItem := NavItem{
		Text: "Home",
		URL:  template.URL("/"),
	}
	menu = append(menu, homeItem)
	// Get sorted keys of page map
	pageNames := make([]string, len(pages))
	i := 0
	for k := range pages {
		pageNames[i] = k
		i++
	}
	sort.Strings(pageNames)
	for _, n := range pageNames {
		if n != "home" {
			item := NavItem{
				Text: kebabToTitle(n),
				URL:  template.URL("/" + n),
			}
			menu = append(menu, item)
		}
	}
	return menu
}

func kebabToTitle(kebab string) string {
	kebab = strings.Replace(kebab, "-", " ", -1)
	return strings.Title(kebab)
}

func getBrand() string {
	wd, _ := os.Getwd()
	parts := strings.Split(wd, string(os.PathSeparator))
	dir := parts[len(parts)-1]
	return strings.ToTitle(dir)
}
