package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PaluMacil/dwn/api/groupapi"
	"github.com/PaluMacil/dwn/api/userapi"
	"github.com/PaluMacil/dwn/app"
	"github.com/PaluMacil/dwn/auth"
	"github.com/PaluMacil/dwn/spa"
)

func main() {
	appModule := app.New()
	spaModule := spa.New(&appModule)
	authModule := auth.New(&appModule)

	userapi := userapi.New(&appModule)
	groupapi := groupapi.New(&appModule)

	mux := http.NewServeMux()
	mux.Handle("/", spaModule)
	mux.Handle("/oauth/", authModule)
	mux.Handle("/api/user/", userapi)
	mux.Handle("/api/group/", groupapi)

	srv := &http.Server{
		Addr:    ":" + appModule.Port,
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
		log.Println("Now serving on port", appModule.Port)
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
