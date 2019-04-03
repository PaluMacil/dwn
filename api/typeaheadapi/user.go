package typeaheadapi

import (
	"encoding/json"
	"github.com/PaluMacil/dwn/module/core"
	"net/url"
)

// api/typeahead/user?query=searchstring
func (rt *TypeaheadRoute) handleUser() {
	switch rt.R.Method {
	case "GET":
		if rt.API().ServeCannot(core.PermissionViewUsers) {
			return
		}
		qry, err := url.QueryUnescape(rt.R.URL.Query().Get("query"))
		if len(qry) < 2 || err != nil {
			rt.API().ServeBadRequest()
			return
		}
		//TODO: only allow logged in users to search
		user, err := rt.Db.Users.CompletionSuggestions(qry) //TODO: check for absence of @ and search by username
		if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
		if err := json.NewEncoder(rt.W).Encode(user); err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
	default:
		rt.API().ServeMethodNotAllowed()
		return
	}
}
