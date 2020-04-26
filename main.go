package main

import (
	"log"
	"os"

	"github.com/PaluMacil/dwn/application"
)

func main() {
	prod := len(os.Args) == 2 && os.Args[1] == "prod"
	if prod {
		log.Printf("starting (prod mode)\n")
	} else {
		log.Printf("starting (dev mode)\n")
	}
	app, err := application.New(prod)
	if err != nil {
		log.Fatalln("could not start app:", err)
	}
	if err := app.Db.Setup.Initialization.EnsureDatabase(); err != nil {
		app.Db.Close()
		log.Fatalln("could not ensure app setup:", err)
	}

	app.Web.Serve()

	app.Db.Close()
	log.Printf("database stopped\n")
}
