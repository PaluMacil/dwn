package api

import (
	"github.com/PaluMacil/dwn/webserver/handler"
	"github.com/gorilla/mux"
)

// RegisterBaseRoutes defines /oauth/...
func RegisterBaseRoutes(rt *mux.Router, factory handler.Factory) {
	rt.Path("/google/{step}").
		Handler(factory.Handler(flowHandler, handler.OptionAllowAnonymous)).
		Methods("GET", "POST")
}

// RegisterAPIRoutes defines /api/oauth/...
func RegisterAPIRoutes(rt *mux.Router, factory handler.Factory) {

}
