package infoapi

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"runtime"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/core"
)

// api/info/server
func (rt *InfoRoute) handleServerInfo(config configuration.Configuration) {
	if rt.Current == nil {
		http.Error(rt.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if rt.API().ServeCannot(core.PermissionViewAppSettings) {
		return
	}
	switch rt.R.Method {
	case "GET":
		info, err := rt.Db.SetupInfo.Get()
		if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		dataSize, _ := dirSize(config.Database.DataDir)

		resp := InfoResponse{
			Config:          config,
			SetupInfo:       info,
			GoVersion:       runtime.Version(),
			NumCPUs:         runtime.NumCPU(),
			AllocatedMemory: m.Alloc,
			DataDirSize:     dataSize,
		}

		if err := json.NewEncoder(rt.W).Encode(resp); err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
	default:
		rt.API().ServeMethodNotAllowed()
	}
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

type InfoResponse struct {
	Config configuration.Configuration `json:"config"`
	core.SetupInfo
	GoVersion       string `json:"goVersion"`
	NumCPUs         int    `json:"numCPUs"`
	AllocatedMemory uint64 `json:"allocatedMemory"`
	DataDirSize     int64  `json:"dataDirSize"`
}
