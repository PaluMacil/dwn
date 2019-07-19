package api

import (
	"github.com/PaluMacil/dwn/webserver/handler"
	"github.com/gorilla/mux"
)

// RegisterRoutes defines /api/core/...
func RegisterRoutes(rt *mux.Router, factory handler.Factory) {
	rt.Path("/groups/{group}").
		Handler(factory.Handler(groupHandler)).
		Methods("GET")
	rt.Path("/groups").
		Handler(factory.Handler(groupsHandler)).
		Methods("GET")
	rt.Path("/groups").
		Handler(factory.Handler(createGroupHandler)).
		Methods("POST")
	rt.Path("/users").
		Handler(factory.Handler(usersHandler)).
		Methods("GET")
	rt.Path("/users/displayname").
		Handler(factory.Handler(userDisplayNameMapHandler, handler.OptionAllowAnonymous)).
		Methods("GET").
		Queries("ids", "{ids}")
	rt.Path("/usergroups/members-of/{group}").
		Handler(factory.Handler(membersOfHandler)).
		Methods("GET")
	rt.Path("/usergroups/groups-for/{userID}").
		Handler(factory.Handler(groupsForHandler)).
		Methods("GET")
	rt.Path("/usergroups").
		Handler(factory.Handler(addUserHandler)).
		Methods("POST")
	rt.Path("/usergroups").
		Handler(factory.Handler(removeUserHandler)).
		Methods("DELETE").
		Queries("userID", "{userID}").
		Queries("group", "{group}")
	rt.Path("/permissions").
		Handler(factory.Handler(permissionsHandler)).
		Methods("GET")
	rt.Path("/permissions").
		Handler(factory.Handler(addPermissionHandler)).
		Methods("PUT").
		Queries("permission", "{permission}").
		Queries("group", "{group}")
	rt.Path("/permissions").
		Handler(factory.Handler(removePermissionHandler)).
		Methods("DELETE").
		Queries("permission", "{permission}").
		Queries("group", "{group}")
	rt.Path("/sessions/login").
		Handler(factory.Handler(loginHandler, handler.OptionAllowAnonymous)).
		Methods("POST")
	rt.Path("/sessions/logout/{token}").
		Handler(factory.Handler(logoutHandler)).
		Methods("DELETE")
	rt.Path("/sessions/me").
		Handler(factory.Handler(meHandler)).
		Methods("GET")
	rt.Path("/sessions").
		Handler(factory.Handler(sessionsHandler)).
		Methods("GET")
}
