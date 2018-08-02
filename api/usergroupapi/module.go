package usergroupapi

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

type UserGroupRoute api.Route

func (ur UserGroupRoute) API() api.Route {
	return api.Route(ur)
}

// api/usergroup/...
func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := UserGroupRoute(api.GetRoute(w, r, mod.Db))
	switch route.Endpoint {
	case "add":
		route.handleAdd()
	case "remove":
		route.handleRemove()
	default:
		http.NotFound(w, r)
	}
}
