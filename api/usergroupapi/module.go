package usergroupapi

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

type UserGroupRoute api.Route

func (rt UserGroupRoute) API() api.Route {
	return api.Route(rt)
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
