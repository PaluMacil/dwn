package api

import (
	"github.com/PaluMacil/dwn/webserver/handler"
	"github.com/gorilla/mux"
)

// RegisterRoutes defines /api/server/info
func RegisterRoutes(rt *mux.Router, factory handler.Factory) {
	rt.Path("/info").
		Handler(factory.Handler(handleInfo)).
		Methods("GET")
}
