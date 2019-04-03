package groupapi

import (
	"encoding/json"
	"github.com/PaluMacil/dwn/module/core"
	"net/http"
)

// api/group/all
func (rt *GroupRoute) handleAll() {
	if rt.Current == nil {
		http.Error(rt.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	switch rt.R.Method {
	case "GET":
		if rt.API().ServeCannot(core.PermissionViewGroups) {
			return
		}
		groups, err := rt.Db.Groups.All()
		if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}

		if err := json.NewEncoder(rt.W).Encode(groups); err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
	default:
		rt.API().ServeMethodNotAllowed()
	}
}
