package groupapi

import (
	"encoding/json"
	"sort"

	"github.com/PaluMacil/dwn/db"
)

// api/group/users-for/{groupname}
func (rt *GroupRoute) handleUsersForGroup() {
	if rt.API().ServeCannot(db.PermissionViewGroups) {
		return
	}
	users, err := rt.Db.Users.UsersFor(rt.ID)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].DisplayName < users[j].DisplayName
	})
	if err := json.NewEncoder(rt.W).Encode(users); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
}
