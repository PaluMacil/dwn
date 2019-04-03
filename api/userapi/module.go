package userapi

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

type UserRoute api.Route

func (rt UserRoute) API() api.Route {
	return api.Route(rt)
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
	case "logout":
		route.handleLogout()
	case "groups-for":
		route.handleGroupsForUser()
	default: // TODO: don't use any bare routes for post / put (slash routing redirects)
		// or simply don't use default on the base api/module url
		route.handleUser()
	}
}
