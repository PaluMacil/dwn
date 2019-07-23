package api

import (
	"github.com/PaluMacil/dwn/webserver/handler"
	"github.com/gorilla/mux"
)

// RegisterRoutes defines /api/registration/...
func RegisterRoutes(rt *mux.Router, factory handler.Factory) {
	rt.Path("/user").
		Handler(factory.Handler(userHandler, handler.OptionAllowAnonymous)).
		Methods("GET")
}
