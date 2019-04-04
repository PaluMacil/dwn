package main

import (
	"log"

	"github.com/PaluMacil/dwn/application"
)

func main() {
	log.Printf("starting\n")
	app, err := application.New()
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
