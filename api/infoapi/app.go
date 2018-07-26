package infoapi

import (
	"encoding/json"
	"net/http"

	"github.com/PaluMacil/dwn/app"
	"github.com/PaluMacil/dwn/db"
	"runtime"
)

// api/info/server
func (rt *InfoRoute) handleServerInfo(app *app.App) {
	if rt.Current == nil {
		http.Error(rt.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if rt.API().ServeCannot(db.PermissionViewAppSettings) {
		return
	}
	switch rt.R.Method {
	case "GET":
		info, err := rt.Db.SetupInfo.Get()
		if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}

		resp := InfoResponse{app, info, runtime.Version()}

		if err := json.NewEncoder(rt.W).Encode(resp); err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
	default:
		rt.API().ServeMethodNotAllowed()
	}
}

type InfoResponse struct {
	*app.App
	db.SetupInfo
	GoVersion string `json:"goVersion"`
}
