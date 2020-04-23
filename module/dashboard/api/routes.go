package api

import (
	"github.com/PaluMacil/dwn/webserver/handler"
	"github.com/gorilla/mux"
)

// RegisterRoutes defines /api/dashboard/...
func RegisterRoutes(rt *mux.Router, factory handler.Factory) {
	rt.Path("/board").
		Handler(factory.Handler(boardHandler)).
		Methods("GET")
}
