package userapi

import (
	"encoding/json"
	"sort"

	"github.com/PaluMacil/dwn/dwn"
)

// api/user/all
func (rt *UserRoute) handleAll() {
	if rt.API().ServeCannot(dwn.PermissionViewUsers) {
		return
	}
	users, err := rt.Db.Users.All()
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
