package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {
	var (
		waitSeconds int
		limitFlag   int
	)
	flag.IntVar(&waitSeconds, "wait", 5, "wait sets the default wait between web requests")
	flag.IntVar(&limitFlag, "limit", 0, "limit sets the max number of pages to check (0 is no limit)")
	flag.Parse()

	var totalPages int
	var wg *sync.WaitGroup = new(sync.WaitGroup)
	const geekURL string = "https://boardgamegeek.com/browse/boardgame/page/%v"
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	if waitSeconds > 0 {
		c.Limit(&colly.LimitRule{
			Delay: time.Duration(waitSeconds) * time.Second,
		})
	}

	c.OnHTML("td.collection_objectname", func(e *colly.HTMLElement) {
		linkElem := e.DOM.Find("a[href]")
		link, _ := linkElem.Attr("href")
		num := strings.Split(link, "/")[2]
		name := linkElem.Text()
		wg.Add(1)
		game := fmt.Sprintf("%s: %s\n", num, name)
		f, err := os.OpenFile("games.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Couldn't open file: %s", err)
		}
		go WriteToFile(game, f, wg)
	})

	c.OnHTML("#main_content", func(e *colly.HTMLElement) {
		if totalPages == 0 {
			lastPageText := e.DOM.Find(`div.infobox a[title="last page"]`).Text()
			totalPages, err := strconv.Atoi(strings.Trim(lastPageText, "[]"))
			if err != nil {
				log.Fatalf("Couldn't parse page total: %s", err)
			}
			for i := 2; i <= coalesceInt(limitFlag, totalPages); i++ {
				c.Visit(fmt.Sprintf(geekURL, i))
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			e := HttpResponseError{
				Code:     r.StatusCode,
				Filename: r.FileName(),
				Body:     string(r.Body),
			}
			log.Println("Error", e.Code, e.Body)
			AddError(e)
		}
	})

	log.Println("First visit", fmt.Sprintf(geekURL, 1))
	c.Visit(fmt.Sprintf(geekURL, 1))

	wg.Wait()
	log.Println("Complete,", totalPages, "games saved.")
	errorData, err := json.Marshal(errorList)
	if err != nil {
		log.Fatalf("marshalling errorList: %s", err)
	}
	err = ioutil.WriteFile("errors.json", errorData, 0666)
	if err != nil {
		log.Fatalf("writing error.json: %s", err)
	}
	log.Println("Errors logged:", len(errorList))
}

var fileLock sync.Mutex

func WriteToFile(game string, f *os.File, w *sync.WaitGroup) {
	fileLock.Lock()
	defer fileLock.Unlock()
	f.Write([]byte(game))
}

var errorList []HttpResponseError

func AddError(e HttpResponseError) {
	errorLock.Lock()
	defer errorLock.Unlock()
	errorList = append(errorList, e)
}

var errorLock sync.Mutex

type HttpResponseError struct {
	Code     int
	Filename string
	Body     string
}

func coalesceInt(nums ...int) int {
	for _, num := range nums {
		if num != 0 {
			return num
		}
	}
	return 0
}
