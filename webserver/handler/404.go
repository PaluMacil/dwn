package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type Status404Handler string

var status404HandlerLookup = map[string]func(w http.ResponseWriter, req *http.Request) {
	"DEFAULT": func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Not Found")
	},
	"DETAILED": func(w http.ResponseWriter, req *http.Request) {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Printf("dumping request in 404 handler: %v", err)
			fmt.Fprintln(w, "Not Found")
			return
		}
		fmt.Fprintf(w, `<html>
					<head><title>Not Found</title></head>
					<body>
						<h3>Not Found</h3>
						<pre>%s</pre>
					</body>
					</html>`, string(requestDump))
	},
}

func (h Status404Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	selectedHandler, exists := status404HandlerLookup[string(h)]
	if !exists {
		log.Println("no 404 handler found with name", h)
		selectedHandler = status404HandlerLookup["DEFAULT"]
	}
	selectedHandler(w, req)
}