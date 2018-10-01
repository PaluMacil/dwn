package main

import (
	"log"

	"github.com/PaluMacil/dwn/application"
	"github.com/PaluMacil/dwn/setup"
)

func main() {
	log.Printf("starting\n")
	app, err := application.New()
	if err != nil {
		log.Fatalln("could not start app:", err)
	}
	setupModule := setup.New(app.Db)
	if err := setupModule.Ensure(); err != nil {
		app.Db.Close()
		log.Fatalln("could not ensure app setup:", err)
	}

	app.Web.Serve()

	app.Db.Close()
	log.Printf("database stopped\n")
}
