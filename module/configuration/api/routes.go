package api

import (
	"github.com/PaluMacil/dwn/webserver/handler"
	"github.com/gorilla/mux"
)

// RegisterRoutes defines /api/configuration/...
func RegisterRoutes(rt *mux.Router, factory handler.Factory) {
	rt.Path("/credentials").
		Handler(factory.Handler(getCredentialHandler)).
		Methods("GET").
		Queries("name", "{credName}", "type", "{credType}")
	rt.Path("/credentials").
		Handler(factory.Handler(getCredentialHandler)).
		Methods("GET").
		Queries("type", "{credType}")
	rt.Path("/credentials").
		Handler(factory.Handler(postCredentialHandler)).
		Methods("POST")
}
