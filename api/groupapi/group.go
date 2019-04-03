package groupapi

import (
	"encoding/json"
	"github.com/PaluMacil/dwn/module/core"
	"net/http"
)

// api/group/{groupname}
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
	case "POST":
		if rt.API().ServeCannot(core.PermissionEditGroups) {
			return
		}
		var request core.GroupCreationRequest
		if err := json.NewDecoder(rt.R.Body).Decode(&request); err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
		if exists, err := rt.Db.Groups.Exists(request.Name); exists {
			rt.API().ServeBadRequest()
			return
		} else if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
		group := request.Group(rt.Current.User.Email)
		if err := rt.Db.Groups.Set(group); err != nil {
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
