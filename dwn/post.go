package dwn

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// PostHandler gets, creates, and updates posts
func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
	case "PUT":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("could not read put request json body:", err)
			http.Error(w, "could not read put request json body", http.StatusInternalServerError)
			return
		}
		var req Post
		err = json.Unmarshal(body, &req)
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("could not read post request json body:", err)
			http.Error(w, "could not read post request json body", http.StatusInternalServerError)
			return
		}
		var req Post
		err = json.Unmarshal(body, &req)
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
