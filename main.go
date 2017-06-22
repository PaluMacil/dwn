package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"fmt"

	"github.com/PaluMacil/dwn/dwn"
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/gob"
)

func main() {
	f, err := os.OpenFile("dwn.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stderr, f))

	dwn.Db, err = storm.Open("dwn.db", storm.Codec(gob.Codec))
	if err != nil {
		log.Panic(err)
	}
	defer dwn.Db.Close()

	mux := http.NewServeMux()
	var admins []dwn.User
	err = dwn.Db.Find("Role", dwn.RoleAdmin, &admins)
	if err == storm.ErrNotFound {
		fmt.Println("No admins detected.")
	} else if err != nil {
		log.Print(err)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	})
	fs := http.FileServer(http.Dir("dist"))
	mux.Handle("/app/", http.StripPrefix("/app/", fs))
	mux.HandleFunc("/api/", dwn.APIHandler)
	mux.HandleFunc("/api/auth/token/", dwn.APIHandler)
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
