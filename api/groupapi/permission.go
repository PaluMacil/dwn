package groupapi

import (
	"encoding/json"
	"github.com/PaluMacil/dwn/db"
	"net/url"
)

// api/group/permission/{groupName}?permission={permission}
func (rt *GroupRoute) handlePermission() {
	if rt.API().ServeCannot(db.PermissionEditGroups) {
		return
	}
	group, err := rt.Db.Groups.Get(rt.ID)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	permission, err := url.QueryUnescape(rt.R.URL.Query().Get("permission"))
	if permission == "" || err != nil {
		rt.API().ServeBadRequest()
	}

	switch rt.R.Method {
	case "PUT": //TODO: Check for elevated permissions
		if !group.HasPermission(permission) {
			group.Permissions = append(group.Permissions, permission)
			err := rt.Db.Groups.Set(group)
			if err != nil {
				rt.API().ServeInternalServerError(err)
				return
			}
		}
	case "DELETE":
		if group.HasPermission(permission) {
			group.Permissions = remove(group.Permissions, permission)
			err := rt.Db.Groups.Set(group)
			if err != nil {
				rt.API().ServeInternalServerError(err)
				return
			}
		}
	default:
		rt.API().ServeMethodNotAllowed()
	}
	if err := json.NewEncoder(rt.W).Encode(group); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}