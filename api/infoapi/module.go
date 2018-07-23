package infoapi

import (
	"net/http"

	"github.com/PaluMacil/dwn/api"
	"github.com/PaluMacil/dwn/app"
)

type Module struct {
	*app.App
}

func New(app *app.App) *Module {
	return &Module{
		App: app,
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
		route.handleServerInfo(mod.App)
	}
}
