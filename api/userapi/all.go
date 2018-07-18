package userapi

import (
	"encoding/json"

	"github.com/PaluMacil/dwn/db"
)

func (rt *UserRoute) handleAll() {
	if rt.API().ServeCannot(db.PermissionViewUsers) {
		return
	}
	users, err := rt.Db.Users.All()
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	if err := json.NewEncoder(rt.W).Encode(users); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
}
