package shoppingapi

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

type ShoppingRoute api.Route

func (rt ShoppingRoute) API() api.Route {
	return api.Route(rt)
}

// api/shopping/...
func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := ShoppingRoute(api.GetRoute(w, r, mod.Db))
	switch route.Endpoint {
	case "all":
		route.handleAll()
	case "remove":
		route.handleRemove()
	default:
		http.NotFound(w, r)
	}
}
