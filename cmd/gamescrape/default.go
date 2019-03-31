package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func defaultCmd(config *Config) {
	scrapeGames(config)

	config.InitialCollector.OnHTML("#main_content", func(e *colly.HTMLElement) {
		pages := config.GetTotalPages()
		if pages == 0 {
			lastPageText := e.DOM.Find(`div.infobox a[title="last page"]`).Text()
			pages, err := strconv.Atoi(strings.Trim(lastPageText, "[]"))
			if err != nil {
				log.Printf("Couldn't parse page total: %s", err)
				return
			}
			config.SetTotalPages(pages)
			for i := 1; i <= coalesceInt(config.LimitFlag, pages); i++ {
				config.MainCollector.Visit(config.GeekURL(i))
			}
		}
	})

	log.Println("First visit", config.GeekURL(1))
	config.InitialCollector.Visit(config.GeekURL(1))
}

func coalesceInt(nums ...int) int {
	for _, num := range nums {
		if num != 0 {
			return num
		}
	}
	return 0
}
