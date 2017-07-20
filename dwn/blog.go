package dwn

import (
	"fmt"
	"net/http"
)

// PostHandler gets, creates, and updates posts
func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Url contains %s!", r.URL.Path[1:])
}
