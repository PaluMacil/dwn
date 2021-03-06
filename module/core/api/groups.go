package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/webserver/errs"
)

// GET /api/core/groups
func groupsHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionViewGroups); err != nil {
		return err
	}
	groups, err := db.Groups.All()
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(groups)
}

// GET /api/core/groups/{group}
func groupHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionViewGroups); err != nil {
		return err
	}
	group, err := db.Groups.Get(vars["group"])
	if db.IsKeyNotFoundErr(err) {
		return errs.StatusNotFound
	} else if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(group)
}

// POST /api/core/groups
func createGroupHandler(
	db *database.Database,
	cur core.Current,
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
		return errs.StatusError{http.StatusBadRequest, errors.New("group already exists")}
	} else if err != nil {
		return err
	}
	group := request.Group(cur.User.ID)
	if err := db.Groups.Set(group); err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(group)
}
