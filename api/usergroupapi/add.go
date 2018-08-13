package usergroupapi

import (
	"encoding/json"
	"time"

	"github.com/PaluMacil/dwn/dwn"
)

// api/usergroup/add
func (rt *UserGroupRoute) handleAdd() {
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
	// check both user and group exist
	userExists, err := rt.Db.Users.Exists(ug.Email)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	groupExists, err := rt.Db.Groups.Exists(ug.GroupName)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	if !userExists || !groupExists {
		rt.API().ServeBadRequest()
		return
	}
	ug.CreatedDate = time.Now()
	err = rt.Db.UserGroups.Set(ug)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	if err := json.NewEncoder(rt.W).Encode(ug); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
}
