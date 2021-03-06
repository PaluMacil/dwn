package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func scrapeGames(config *Config) {
	config.MainCollector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	config.MainCollector.OnHTML("td.collection_objectname", func(e *colly.HTMLElement) {
		linkElem := e.DOM.Find("a[href]")
		link, _ := linkElem.Attr("href")
		num := strings.Split(link, "/")[2]
		name := linkElem.Text()
		game := fmt.Sprintf("%s:%s\n", num, name)
		config.GameList.Add(game)
	})
}
