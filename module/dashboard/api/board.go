package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
)

// /api/dashboard/board
func boardHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionListProjects); err != nil {
		return err
	}
	board, err := db.Dashboard.Board.Get()
	if err != nil {
		return fmt.Errorf("getting dashboard: %w", err)
	}

	return json.NewEncoder(w).Encode(board)
}
