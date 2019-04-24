package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/module/setup"
)

// api/server/info
func handleInfo(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionViewAppSettings); err != nil {
		return err
	}

	info, err := db.Setup.Initialization.Get()
	if err != nil {
		return err
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	dataSize, _ := dirSize(config.Database.DataDir)

	resp := InfoResponse{
		Config:          config,
		Initialization:  info,
		GoVersion:       runtime.Version(),
		NumCPUs:         runtime.NumCPU(),
		AllocatedMemory: m.Alloc,
		DataDirSize:     dataSize,
	}

	return json.NewEncoder(w).Encode(resp)
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
	setup.Initialization
	GoVersion       string `json:"goVersion"`
	NumCPUs         int    `json:"numCPUs"`
	AllocatedMemory uint64 `json:"allocatedMemory"`
	DataDirSize     int64  `json:"dataDirSize"`
}
