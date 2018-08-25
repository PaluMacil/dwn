package searchapi

import (
	"net/http"

	"github.com/PaluMacil/dwn/api"
	"github.com/PaluMacil/dwn/database"
)

type Module struct {
	Db *database.Database
}

func New(db *database.Database) *Module {
	return &Module{
		Db: db,
	}
}

type SearchRoute api.Route

func (rt SearchRoute) API() api.Route {
	return api.Route(rt)
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
