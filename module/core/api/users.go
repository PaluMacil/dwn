package api

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
)

// GET /api/core/users
func usersHandler(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
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
	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		return err
	}

	return nil
}
