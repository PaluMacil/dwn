package groupapi

import (
	"encoding/json"
	"net/http"
)

// api/user/{groupname}
func (rt *GroupRoute) handleGroup() {
	if rt.Current == nil {
		http.Error(rt.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	switch rt.R.Method {
	case "GET":
		group, err := rt.Db.Groups.Get(rt.ID)
		if err != nil { //TODO: first check not exists
			rt.API().ServeInternalServerError(err)
			return
		}
		if err := json.NewEncoder(rt.W).Encode(group); err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
	default:
		rt.API().ServeMethodNotAllowed()
	}
}
