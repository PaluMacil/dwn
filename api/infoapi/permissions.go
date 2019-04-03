package infoapi

import (
	"encoding/json"
	"net/http"

	"github.com/PaluMacil/dwn/core"
)

// api/info/permissions
func (rt *InfoRoute) handlePermissions() {
	if rt.Current == nil {
		http.Error(rt.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	switch rt.R.Method {
	case "GET":
		if err := json.NewEncoder(rt.W).Encode(core.Permissions); err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
	default:
		rt.API().ServeMethodNotAllowed()
	}
}
