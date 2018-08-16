package main

import (
	"log"

	"github.com/PaluMacil/dwn/app"
	"github.com/PaluMacil/dwn/setup"
	"github.com/PaluMacil/dwn/webserver"
)

func main() {
	appModule, err := app.New()
	if err != nil {
		log.Fatalln("could not start app:", err)
	}
	setupModule := setup.New(appModule)
	if err := setupModule.Ensure(); err != nil {
		appModule.Db.Close()
		log.Fatalln("could not ensure app setup:", err)
	}

	web := webserver.New(appModule)
	web.Serve()

	appModule.Db.Close()
	log.Printf("Badger: database stopped\n")
}
