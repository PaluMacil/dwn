package echo

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type EchoModule struct {
	history *History
}

func (mod EchoModule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
	}
	request := string(requestDump)
	mod.history.add(request)
	_, _ = fmt.Fprintf(w, request)
}
