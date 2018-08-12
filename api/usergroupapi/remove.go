package usergroupapi

import (
	"encoding/json"

	"github.com/PaluMacil/dwn/db"
)

// api/usergroup/remove
func (rt *UserGroupRoute) handleRemove() {
	if rt.API().ServeCannot(db.PermissionEditGroups) {
		return
	}
	if rt.R.Method != "POST" {
		rt.API().ServeMethodNotAllowed()
		return
	}
	if rt.R.Body == nil {
		rt.API().ServeBadRequest()
		return
	}
	var ug db.UserGroup
	err := json.NewDecoder(rt.R.Body).Decode(&ug)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	err = rt.Db.UserGroups.Delete(ug.Email, ug.GroupName)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	if err := json.NewEncoder(rt.W).Encode(ug); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
}
