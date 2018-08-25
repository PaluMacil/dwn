package typeaheadapi

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

type TypeaheadRoute api.Route

func (rt TypeaheadRoute) API() api.Route {
	return api.Route(rt)
}

// api/typeahead/...
func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := TypeaheadRoute(api.GetRoute(w, r, mod.Db))
	switch route.Endpoint {
	case "user":
		route.handleUser()
	}
}
