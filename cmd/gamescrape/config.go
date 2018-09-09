package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

type Config struct {
	WaitSeconds      int
	LimitFlag        int
	GameList         *FileManager
	Errors           ErrorList
	totalPages       uint32
	MainCollector    *colly.Collector
	InitialCollector *colly.Collector
}

func (c *Config) GetTotalPages() int {

	return int(atomic.LoadUint32(&c.totalPages))
}

func (c *Config) SetTotalPages(pages int) {
	atomic.StoreUint32(&c.totalPages, uint32(pages))
}

func (c *Config) GeekURL(num int) string {
	return fmt.Sprintf("https://boardgamegeek.com/browse/boardgame/page/%v", num)
}

func NewConfig() *Config {
	os.Mkdir("data", 0777)
	var (
		waitSeconds int
		limitFlag   int
	)
	flag.IntVar(&waitSeconds, "wait", 5, "wait sets the default wait between web requests")
	flag.IntVar(&limitFlag, "limit", 0, "limit sets the max number of pages to check (0 is no limit)")
	flag.Parse()
	config := &Config{
		WaitSeconds: waitSeconds,
		LimitFlag:   limitFlag,
		GameList:    &FileManager{filename: filepath.Join("data", "games.txt")},
		Errors:      ErrorList{},
		totalPages:  0,
		MainCollector: colly.NewCollector(
			colly.Async(waitSeconds == 0),
		),
		InitialCollector: colly.NewCollector(),
	}

	extensions.RandomUserAgent(config.InitialCollector)
	extensions.RandomUserAgent(config.MainCollector)
	if waitSeconds > 0 {
		config.MainCollector.Limit(&colly.LimitRule{
			Delay: time.Duration(waitSeconds) * time.Second,
		})
	} else {
		config.MainCollector.Limit(&colly.LimitRule{
			Delay:       time.Duration(0) * time.Second,
			Parallelism: 20,
			DomainGlob:  "*",
		})
		// Some hints on transport tweaks for scraping
		// http://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/
		defaultRoundTripper := http.DefaultTransport
		defaultRoundTripperPointer := defaultRoundTripper.(*http.Transport)
		t := *defaultRoundTripperPointer // deref to get copy
		t.MaxIdleConns = 0
		t.MaxIdleConnsPerHost = 1000
		config.MainCollector.WithTransport(&t)
	}

	// set error handling
	re := regexp.MustCompile("page/([0-9]*)(?:|$)")
	responseHandler := func(r *colly.Response) {
		if r.StatusCode != 200 {
			pageString := re.FindStringSubmatch(r.Request.URL.String())[1]
			page, err := strconv.Atoi(pageString)
			if err != nil {
				log.Println("parsing page string", err)
			}
			e := HttpResponseError{
				Code: r.StatusCode,
				Body: string(r.Body),
				Page: page,
			}
			log.Println("Error", e.Code, e.Body)
			config.Errors.Add(e)
		}
	}
	config.InitialCollector.OnResponse(responseHandler)
	config.MainCollector.OnResponse(responseHandler)
	errHandler := func(r *colly.Response, responseErr error) {
		pageString := re.FindStringSubmatch(r.Request.URL.String())[1]
		page, err := strconv.Atoi(pageString)
		if err != nil {
			fmt.Println("error parsing '", pageString, "' from url", r.Request.URL.String())
			return
		}
		e := HttpResponseError{
			Code:  r.StatusCode,
			Body:  string(r.Body),
			Error: responseErr.Error(),
			Page:  page,
		}
		log.Println("Error", e.Code, e.Error)
		config.Errors.Add(e)
	}
	config.MainCollector.OnError(errHandler)
	config.InitialCollector.OnError(errHandler)

	return config
}
