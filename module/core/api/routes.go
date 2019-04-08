package api

import (
	"github.com/PaluMacil/dwn/module"
	"github.com/gorilla/mux"
)

func RegisterRoutes(rt *mux.Router, factory module.HandlerFactory) {
	rt.PathPrefix("/groups").
		Handler(factory.Handler(groupsHandler)).
		Methods("GET", "POST")
}
