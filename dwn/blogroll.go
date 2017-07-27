package dwn

import "net/http"

// PostHandler gets, creates, and updates posts
func BlogRollHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

type BlogRoll struct {
	Posts          []Post
	StartIndex     int
	TotalAvailable int
}
