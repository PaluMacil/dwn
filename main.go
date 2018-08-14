package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PaluMacil/dwn/api/groupapi"
	"github.com/PaluMacil/dwn/api/infoapi"
	"github.com/PaluMacil/dwn/api/userapi"
	"github.com/PaluMacil/dwn/api/usergroupapi"
	"github.com/PaluMacil/dwn/app"
	"github.com/PaluMacil/dwn/auth"
	"github.com/PaluMacil/dwn/setup"
	"github.com/PaluMacil/dwn/spa"
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

	spaModule := spa.New(appModule)
	authModule := auth.New(appModule)

	userModule := userapi.New(appModule)
	groupModule := groupapi.New(appModule)
	usergroupModule := usergroupapi.New(appModule)
	infoModule := infoapi.New(appModule)

	mux := http.NewServeMux()
	mux.Handle("/", spaModule)
	mux.Handle("/oauth/", authModule)
	mux.Handle("/api/user/", userModule)
	mux.Handle("/api/group/", groupModule)
	mux.Handle("/api/usergroup/", usergroupModule)
	mux.Handle("/api/info/", infoModule)

	srv := &http.Server{
		Addr:    ":" + appModule.Config.WebServer.Port,
		Handler: mux,
		//time from when the connection is accepted to when the request body is fully read
		ReadTimeout: 5 * time.Second,
		//time from the end of the request header read to the end of the response write
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	go func() {
		log.Println("Now serving on port", appModule.Config.WebServer.Port)
		log.Println(srv.ListenAndServe())
	}()
	<-stop

	log.Printf("shutting down ...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	appModule.Db.Close()
	log.Printf("Badger: database stopped\n")
}
