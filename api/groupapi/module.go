package groupapi

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

type GroupRoute api.Route

func (r GroupRoute) API() api.Route {
	return api.Route(r)
}

// api/group/...
func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := GroupRoute(api.GetRoute(w, r, mod.Db))
	switch route.Endpoint {
	case "all":
		route.handleAll()
	default:
		route.handleGroup()
	}
}
