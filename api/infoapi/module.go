package infoapi

import (
	"net/http"

	"github.com/PaluMacil/dwn/configuration"

	"github.com/PaluMacil/dwn/api"
	"github.com/PaluMacil/dwn/database"
)

type Module struct {
	Config configuration.Configuration
	Db     *database.Database
}

func New(db *database.Database, config configuration.Configuration) *Module {
	return &Module{
		Db:     db,
		Config: config,
	}
}

type InfoRoute api.Route

func (rt InfoRoute) API() api.Route {
	return api.Route(rt)
}

// api/info/...
func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := InfoRoute(api.GetRoute(w, r, mod.Db))
	switch route.Endpoint {
	case "server":
		route.handleServerInfo(mod.Config)
	case "permissions":
		route.handlePermissions()
	}
}
