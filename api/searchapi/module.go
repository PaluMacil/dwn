package searchapi

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

type SearchRoute api.Route

func (ur SearchRoute) API() api.Route {
	return api.Route(ur)
}

// api/search/...
func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := SearchRoute(api.GetRoute(w, r, mod.Db))
	switch route.Endpoint {
	case "user":
		route.handleUser()
	case "index":
		route.handleIndex()
	}
}
