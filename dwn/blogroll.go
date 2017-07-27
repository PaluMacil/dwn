package dwn

import (
	"encoding/json"
	"log"
	"net/http"
)

// BlogRollHandler gets, creates, and updates posts
func BlogRollHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//TODO: Get topic-based blog rolls vias query := r.URL.Query()
	var posts []Post
	err := Db.All(&posts)
	if err != nil {
		log.Println("could not get all posts:", err)
		http.Error(w, "ould not get all posts", http.StatusInternalServerError)
		return
	}
	blogroll := BlogRoll{
		Posts:          posts,
		StartIndex:     1,
		PageSize:       len(posts),
		TotalAvailable: len(posts),
	}
	json.NewEncoder(w).Encode(blogroll)
}

// BlogRoll formats a slice of returned posts
type BlogRoll struct {
	Posts          []Post
	StartIndex     int
	PageSize       int
	TotalAvailable int
}
