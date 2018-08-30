package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	config := NewConfig()
	if len(os.Args) == 1 ||
		(len(os.Args) >= 2 && strings.ToLower(os.Args[1]) == "default") ||
		(len(os.Args) >= 2 && os.Args[1][:1] == "-") {
		defaultCmd(config)
	} else if len(os.Args) >= 2 && strings.ToLower(os.Args[1]) == "retry" {
		retryCmd(config)
	} else {
		log.Fatalln("Invalid command:", os.Args[1:])
	}

	config.MainCollector.Wait()

	log.Println("Complete,", config.GetTotalPages()-len(config.Errors), "games saved.")
	err := config.Errors.WriteFile()
	if err != nil {
		log.Fatalf("writing error.json: %s", err)
	}
	log.Println("Errors logged:", len(config.Errors))
}
