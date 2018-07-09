package main

import (
	"log"
	"net/http"
	"time"

	"github.com/PaluMacil/dwn/app"
	"github.com/PaluMacil/dwn/auth"
	"github.com/PaluMacil/dwn/spa"
)

func main() {
	appModule := app.NewApp()
	spaModule := spa.New(&appModule)
	authModule := auth.New(&appModule)

	mux := http.NewServeMux()
	mux.Handle("/", spaModule)

	mux.Handle("/test", appModule)
	mux.Handle("/oauth/", authModule)
	//mux.HandleFunc("/oauth/google/login", app.HandleGoogleLogin)
	//mux.HandleFunc("/oauth/google/callback", app.HandleGoogleCallback)

	srv := &http.Server{
		Addr:    ":" + appModule.Port,
		Handler: mux,
		//time from when the connection is accepted to when the request body is fully read
		ReadTimeout: 5 * time.Second,
		//time from the end of the request header read to the end of the response write
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println("Now serving on port", appModule.Port)
	log.Println(srv.ListenAndServe())
}
