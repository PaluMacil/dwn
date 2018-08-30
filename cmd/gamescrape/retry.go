package main

import "log"

func retryCmd(config *Config) {
	scrapeGames(config)
	errorQueue := config.Errors.FromFile()
	pages := len(errorQueue)
	config.SetTotalPages(pages)
	log.Println("Loaded", pages, "errors for retry...")
	for i := 0; i < pages; i++ {
		pageNum := errorQueue[i].Page
		config.MainCollector.Visit(config.GeekURL(pageNum))
	}
}
