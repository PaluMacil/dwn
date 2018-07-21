package userapi

import (
	"encoding/json"

	"github.com/PaluMacil/dwn/db"
)

// api/user/{email or username}
func (rt *UserRoute) handleUser() {
	switch rt.R.Method {
	case "GET":
		if rt.API().ServeCannot(db.PermissionViewUsers) {
			return
		}
		user, err := rt.Db.Users.Get(rt.ID) //TODO: check for absence of @ and search by username
		if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
		if err := json.NewEncoder(rt.W).Encode(user); err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
	case "PUT": //TODO: Except password; set modified by and modified date
		if rt.API().ServeCannot(db.PermissionEditUserInfo) {
			return
		}
		var user db.User
		err := json.NewDecoder(rt.R.Body).Decode(&user)
		if err != nil {
			rt.API().ServeBadRequest()
			return
		}
		rt.Db.Users.Set(user)
	case "POST": //TODO: check for conflict (exists) before setting
		if rt.API().ServeCannot(db.PermissionEditUserInfo) {
			return
		}
		var user db.User
		err := json.NewDecoder(rt.R.Body).Decode(&user)
		if err != nil {
			rt.API().ServeBadRequest()
			return
		}
		rt.Db.Users.Set(user)
	case "DELETE": //TODO: determine if this actually deletes or if it sets an inactive bit that put can't modify
		if rt.API().ServeCannot(db.PermissionEditUserInfo) {
			return
		}
		if err := rt.Db.Users.Delete(rt.ID); err != nil {
			rt.API().ServeInternalServerError(err)
		}
	default:
		rt.API().ServeMethodNotAllowed()
	}
}
