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

	err = dwn.Db.Init(&dwn.User{})
	if err != nil {
		log.Panic(err)
	}
	err = dwn.Db.Init(&dwn.Session{})
	if err != nil {
		log.Panic(err)
	}

	var admins []dwn.User
	err = dwn.Db.Find("Role", dwn.RoleAdmin, &admins)

	if err == storm.ErrNotFound {
		fmt.Println("No admin users detected!")
		var user = dwn.User{
			Role:      dwn.RoleAdmin,
			CreatedAt: time.Now(),
		}
		fmt.Println("You must create an admin before continuing...")
		fmt.Println("Name:")
		fmt.Scanln(&user.Name)
		fmt.Println("Email:")
		fmt.Scanln(&user.Email)
		fmt.Println("Password:")
		var plainPassword string
		fmt.Scanln(&plainPassword)
		user.Password, err = dwn.HashPassword(plainPassword)
		if err != nil {
			log.Println("Error hashing password:", err)
		}
		dwn.Db.Save(&user)
		if err != nil {
			log.Println("Could not save user:", err)
		} else {
			log.Println("Saved!", user)
		}
	} else if err != nil {
		log.Print(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	})
	fs := http.FileServer(http.Dir("dist"))
	mux.Handle("/app/", http.StripPrefix("/app/", fs))
	mux.HandleFunc("/api/", dwn.APIHandler)
	//no trailing slash in pattern for exact match
	mux.HandleFunc("/api/account/token", dwn.TokenHandler)
	//TODO: mux.HandleFunc("/api/auth/user", dwn.UserHandler)
	//TODO: mux.HandleFunc("/api/blog/post", dwn.PostHandler)
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
