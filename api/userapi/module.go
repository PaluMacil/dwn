package userapi

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

type UserRoute api.Route

func (ur UserRoute) API() api.Route {
	return api.Route(ur)
}

// api/user/...
func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := UserRoute(api.GetRoute(w, r, mod.Db))
	switch route.Endpoint {
	case "me":
		route.handleMe()
	case "all":
		route.handleAll()
	case "sessions":
		route.handleSessionDetails()
	default:
		route.handleUser()
	}
}
