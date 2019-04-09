package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module"
	"github.com/PaluMacil/dwn/module/core"
)

// GET /api/core/groups/{group}
func groupsHandler(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionViewGroups); err != nil {
		return err
	}
	group, err := db.Groups.Get(vars["group"])
	if db.IsKeyNotFoundErr(err) {
		return module.StatusNotFound
	} else if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(group); err != nil {
		return err
	}

	return nil
}

// POST /api/core/groups
func createGroupHandler(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionEditGroups); err != nil {
		return err
	}
	var request core.GroupCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}
	if exists, err := db.Groups.Exists(request.Name); exists {
		return module.StatusError{http.StatusBadRequest, errors.New("group already exists")}
	} else if err != nil {
		return err
	}
	group := request.Group(cur.User.Email)
	if err := db.Groups.Set(group); err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(group); err != nil {
		return err
	}

	return nil
}
