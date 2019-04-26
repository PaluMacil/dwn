package api

import (
	"github.com/PaluMacil/dwn/webserver/handler"
	"github.com/gorilla/mux"
)

// RegisterRoutes defines /api/shopping/...
func RegisterRoutes(rt *mux.Router, factory handler.Factory) {
	rt.Path("/items").
		Handler(factory.Handler(listHandler)).
		Methods("GET")
	rt.Path("/items").
		Handler(factory.Handler(addHandler)).
		Methods("POST")
	rt.Path("/items").
		Handler(factory.Handler(removeHandler)).
		Methods("DELETE").
		Queries("name", "{name}")
}
