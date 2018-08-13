package usergroupapi

import (
	"encoding/json"

	"github.com/PaluMacil/dwn/dwn"
)

// api/usergroup/remove
func (rt *UserGroupRoute) handleRemove() {
	if rt.API().ServeCannot(dwn.PermissionEditGroups) {
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
	var ug dwn.UserGroup
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
