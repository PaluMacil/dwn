package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/webserver/errs"
)

// /api/core/permissions
func permissionsHandler(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := json.NewEncoder(w).Encode(core.Permissions); err != nil {
		return err
	}

	return nil
}

// PUT /api/core/permissions
func addPermissionHandler(
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

	groupName := vars["group"]
	group, err := db.Groups.Get(groupName)
	if err != nil {
		return err
	}
	permission, err := url.QueryUnescape(vars["permission"])
	if permission == "" || err != nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("invalid permission")}
	}
	if !group.HasPermission(permission) {
		group.Permissions = append(group.Permissions, permission)
		group.ModifiedBy = cur.User.Email
		group.ModifiedDate = time.Now()
		err := db.Groups.Set(group)
		if err != nil {
			return err
		}
	}

	return nil
}

// DELETE /api/core/permissions
func removePermissionHandler(
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
	groupName := vars["group"]
	group, err := db.Groups.Get(groupName)
	if err != nil {
		return err
	}
	permission, err := url.QueryUnescape(vars["permission"])
	if permission == "" || err != nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("invalid permission")}
	}
	if group.HasPermission(permission) {
		group.Permissions = remove(group.Permissions, permission)
		group.ModifiedBy = cur.User.Email
		group.ModifiedDate = time.Now()
		err := db.Groups.Set(group)
		if err != nil {
			return err
		}
	}

	return nil
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
