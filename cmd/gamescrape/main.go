package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	c.Limit(&colly.LimitRule{
		Delay: 5 * time.Second,
	})

	c.OnHTML("td.collection_objectname", func(e *colly.HTMLElement) {
		linkElem := e.DOM.Find("a[href]")
		link, _ := linkElem.Attr("href")
		num := strings.Split(link, "/")[2]
		name := linkElem.Text()
		fmt.Printf("%s: %s\n", num, name)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(`https://boardgamegeek.com/browse/boardgame/page/1`)
}
