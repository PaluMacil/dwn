package dwn

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// PostHandler gets, creates, and updates posts
func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post Post
	switch r.Method {
	case "GET":
		var id int
		qID := r.URL.Query().Get("ID")
		id, err := strconv.Atoi(qID)
		if err != nil {
			http.Error(w, "Post ID not valid", http.StatusBadRequest)
			return
		}
		if err := Db.One("ID", id, &post); err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
	case "PUT":
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			log.Println("PostHandler: could not unmarshal put request json body:", err)
			http.Error(w, "could not unmarshal put request json body", http.StatusInternalServerError)
			return
		}
		err = Db.Save(&post)
		if err != nil {
			log.Println("PostHandler: could not save put request json body:", err)
			http.Error(w, "could not save put request json body", http.StatusInternalServerError)
			return
		}
	case "POST":
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			log.Println("PostHandler: could not unmarshal post request json body:", err)
			http.Error(w, "could not unmarshal post request json body", http.StatusInternalServerError)
			return
		}
		err = Db.Save(&post)
		if err != nil {
			log.Println("PostHandler: could not save post request json body:", err)
			http.Error(w, "could not save post request json body", http.StatusInternalServerError)
			return
		}
	}
}

// Post stores a blog post
type Post struct {
	ID       int `storm:"id,increment"`
	Topic    Topic
	Title    string
	Slug     string
	Author   UserInfo
	Body     string
	Format   PostFormat
	IsStub   bool
	Modified time.Time
	Created  time.Time `storm:"index"`
}

// Topic indicates the topic of a post
type Topic int

const (
	// TopicAll indicates all topics
	TopicAll Topic = iota
	// TopicPersonal indicates personal topics
	TopicPersonal
	// TopicTech indicates tech and programming topics
	TopicTech
	// TopicFood indicates food and dining topics
	TopicFood
	// TopicFun indicates game and vacation topics
	TopicFun
)

// PostFormat indicates the format of a post
type PostFormat int

const (
	// PostFormatPlain indicates plaintext
	PostFormatPlain PostFormat = iota
	// PostFormatMarkdown indicates markdown
	PostFormatMarkdown
)
