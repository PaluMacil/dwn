package usergroupapi

import (
	"encoding/json"
	"time"

	"github.com/PaluMacil/dwn/db"
)

// api/usergroup/add
func (rt *UserGroupRoute) handleAdd() {
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
	ug.CreatedDate = time.Now()
	err = rt.Db.UserGroups.Set(ug)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
}
