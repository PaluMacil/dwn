package main

import (
	"flag"
	"fmt"
	"log"
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
	totalPages       *uint32
	MainCollector    *colly.Collector
	InitialCollector *colly.Collector
}

func (c *Config) GetTotalPages() int {
	return int(atomic.LoadUint32(c.totalPages))
}

func (c *Config) SetTotalPages(pages int) {
	atomic.StoreUint32(c.totalPages, uint32(pages))
}

func (c *Config) GeekURL(num int) string {
	return fmt.Sprintf("https://boardgamegeek.com/browse/boardgame/page/%v", num)
}

func NewConfig() *Config {
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
		GameList:    &FileManager{filename: "games.txt"},
		MainCollector: colly.NewCollector(
			colly.Async(true),
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
	}

	// set error handling
	responseHandler := func(r *colly.Response) {
		if r.StatusCode != 200 {
			e := HttpResponseError{
				Code: r.StatusCode,
				Body: string(r.Body),
			}
			log.Println("Error", e.Code, e.Body)
			config.Errors.Add(e)
		}
	}
	config.InitialCollector.OnResponse(responseHandler)
	config.MainCollector.OnResponse(responseHandler)
	errHandler := func(r *colly.Response, err error) {
		e := HttpResponseError{
			Code:  r.StatusCode,
			Body:  string(r.Body),
			Error: err.Error(),
		}
		log.Println("Error", e.Code, e.Error)
		config.Errors.Add(e)
	}
	config.MainCollector.OnError(errHandler)
	config.InitialCollector.OnError(errHandler)

	return config
}
