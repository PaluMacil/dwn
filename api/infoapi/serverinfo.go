package infoapi

import (
	"encoding/json"
	"net/http"

	"runtime"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/dwn"
)

// api/info/server
func (rt *InfoRoute) handleServerInfo(config configuration.Configuration) {
	if rt.Current == nil {
		http.Error(rt.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if rt.API().ServeCannot(dwn.PermissionViewAppSettings) {
		return
	}
	switch rt.R.Method {
	case "GET":
		info, err := rt.Db.SetupInfo.Get()
		if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}

		resp := InfoResponse{config, info, runtime.Version(), runtime.NumCPU()}

		if err := json.NewEncoder(rt.W).Encode(resp); err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
	default:
		rt.API().ServeMethodNotAllowed()
	}
}

type InfoResponse struct {
	Config configuration.Configuration `json:"config"`
	dwn.SetupInfo
	GoVersion string `json:"goVersion"`
	NumCPUs   int    `json:"numCPUs"`
}
