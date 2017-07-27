package dwn

import "net/http"

// PostHandler gets, creates, and updates posts
func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

type Post struct {
	ID     int
	Title  string
	Author UserInfo
	Body   string
	Format PostFormat
	IsStub bool
}

type PostFormat int

const (
	PostFormatPlain PostFormat = iota
	PostFormatMarkdown
)
