package infoapi

import (
	"encoding/json"
	"github.com/PaluMacil/dwn/db"
	"net/http"
)

// api/info/permissions
func (rt *InfoRoute) handlePermissions() {
	if rt.Current == nil {
		http.Error(rt.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	switch rt.R.Method {
	case "GET":
		if err := json.NewEncoder(rt.W).Encode(db.Permissions); err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
	default:
		rt.API().ServeMethodNotAllowed()
	}
}
