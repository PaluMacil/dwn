package api

import (
	"encoding/json"
	"errors"
	"github.com/PaluMacil/dwn/database/store"
	"github.com/PaluMacil/dwn/webserver/errs"
	"net/http"
	"sort"
	"strings"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
)

// GET /api/core/users
func usersHandler(
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
	users, err := db.Users.All()
	if err != nil {
		return err
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].DisplayName < users[j].DisplayName
	})
	userInfo := core.Users(users).Info()

	return json.NewEncoder(w).Encode(userInfo)
}

// DELETE /api/core/users
func deleteUserHandler(
	db *database.Database,
	config configuration.Configuration,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionEditUserInfo); err != nil {
		return err
	}
	id, err := store.StringToIdentity(vars["id"])
	if err != nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("invalid userID")}
	}
	if err := db.Users.Delete(id); err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(true)
}

// GET /api/core/users/displayname?ids={2,3,4}
func userDisplayNameMapHandler(
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
	idStrings := strings.Split(vars["ids"], ",")
	idDisplayNameMap := make(map[store.Identity]core.DisplayName)
	for _, idString := range idStrings {
		id, err := store.StringToIdentity(idString)
		if err != nil {
			return err
		}
		user, err := db.Users.Get(id)
		if err != nil {
			return err
		}
		idDisplayNameMap[id] = user.DisplayName
	}

	return json.NewEncoder(w).Encode(idDisplayNameMap)
}
