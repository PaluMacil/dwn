package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PaluMacil/dwn/dwn"
)

func main() {
	f, err := os.OpenFile("dwn.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stderr, f))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	})
	fs := http.FileServer(http.Dir("dist"))
	mux.Handle("/app/", http.StripPrefix("/app/", fs))
	mux.HandleFunc("/api/", dwn.APIHandler)
	srv := &http.Server{
		Addr:    ":1337",
		Handler: mux,
		//time from when the connection is accepted to when the request body is fully read
		ReadTimeout: 5 * time.Second,
		//time from the end of the request header read to the end of the response write
		WriteTimeout: 10 * time.Second,
	}

	log.Println(srv.ListenAndServe())
}
