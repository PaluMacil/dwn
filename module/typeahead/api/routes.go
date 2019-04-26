package api

import (
	"github.com/PaluMacil/dwn/webserver/handler"
	"github.com/gorilla/mux"
)

// RegisterRoutes defines /api/typeahead/...
func RegisterRoutes(rt *mux.Router, factory handler.Factory) {
	rt.Path("/users").
		Handler(factory.Handler(usersHandler)).
		Methods("GET").
		Queries("query", "{query}")
}
