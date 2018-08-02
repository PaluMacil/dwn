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

func (rt GroupRoute) API() api.Route {
	return api.Route(rt)
}

// api/group/...
func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := GroupRoute(api.GetRoute(w, r, mod.Db))
	switch route.Endpoint {
	case "all":
		route.handleAll()
	case "permission":
		route.handlePermission()
	case "users-for":
		route.handleUsersForGroup()
	default:
		route.handleGroup()
	}
}
