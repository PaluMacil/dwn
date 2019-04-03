package userapi

import (
	"encoding/json"
	"github.com/PaluMacil/dwn/module/core"
	"net/url"
	"sort"
)

// api/user/groups-for/{email}
func (rt *UserRoute) handleGroupsForUser() {
	if rt.API().ServeCannot(core.PermissionViewGroups) {
		return
	}
	email, err := url.QueryUnescape(rt.ID)
	if err != nil {
		rt.API().ServeBadRequest()
		return
	}
	groups, err := rt.Db.Groups.GroupsFor(email)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	if err := json.NewEncoder(rt.W).Encode(groups); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
}
