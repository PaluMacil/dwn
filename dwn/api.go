package dwn

import (
	"fmt"
	"net/http"
)

// APIHandler is a filler example as I work on the API
func APIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Url contains %s!", r.URL.Path[1:])
}
