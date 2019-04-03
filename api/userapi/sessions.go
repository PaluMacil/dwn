package userapi

import (
	"encoding/json"
	"github.com/PaluMacil/dwn/module/core"
)

// api/user/sessions
func (rt *UserRoute) handleSessionDetails() {
	// TODO: give sessions their own api
	if rt.API().ServeCannot(core.PermissionViewUsers) {
		return
	}
	sessions, err := rt.Db.Sessions.All()
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	details := make([]SessionDetails, len(sessions))
	for i, s := range sessions {
		u, err := rt.Db.Users.Get(s.Email)
		if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
		details[i].Session = s
		details[i].User = u
	}

	if err := json.NewEncoder(rt.W).Encode(details); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
}

type SessionDetails struct {
	User    core.User    `json:"user"`
	Session core.Session `json:"session"`
}
