package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/store"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/webserver/errs"
)

// POST /api/core/usergroups
func addUserHandler(
	db *database.Database,
	config configuration.Configuration,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionEditGroups); err != nil {
		return err
	}
	if r.Body == nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("no request body")}
	}
	var ug core.UserGroup
	err := json.NewDecoder(r.Body).Decode(&ug)
	if err != nil {
		return err
	}
	// check both user and group exist
	userExists, err := db.Users.Exists(ug.UserID)
	if err != nil {
		return err
	}
	groupExists, err := db.Groups.Exists(ug.GroupName)
	if err != nil {
		return err
	}
	if !userExists || !groupExists {
		return errs.StatusError{http.StatusBadRequest, errors.New("user or group doesn't exist")}
	}
	ug.CreatedDate = time.Now()
	err = db.UserGroups.Set(ug)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(ug)
}

// DELETE /api/core/usergroups
func removeUserHandler(
	db *database.Database,
	config configuration.Configuration,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionEditGroups); err != nil {
		return err
	}
	if r.Body == nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("no request body")}
	}

	userId, err := store.StringToIdentity(vars["userID"])
	if err != nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("invalid userId")}
	}
	group := vars["group"]
	ug := core.UserGroup{
		UserID:    userId,
		GroupName: group,
	}

	err = db.UserGroups.Delete(ug.UserID, ug.GroupName)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(ug)
}

// GET /api/core/usergroups/members-of/{group}
func membersOfHandler(
	db *database.Database,
	config configuration.Configuration,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionViewUsers); err != nil {
		return err
	}
	users, err := db.Users.UsersFor(vars["group"])
	if err != nil {
		return err
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].DisplayName < users[j].DisplayName
	})
	userInfo := core.Users(users).Info()
	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		return err
	}

	return nil
}

// GET /api/core/usergroups/groups-for/{userID}
func groupsForHandler(
	db *database.Database,
	config configuration.Configuration,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionViewGroups); err != nil {
		return err
	}
	userId, err := store.StringToIdentity(vars["userID"])
	if err != nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("invalid userId")}
	}
	groups, err := db.Groups.GroupsFor(userId)
	if err != nil {
		return err
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return json.NewEncoder(w).Encode(groups)
}
