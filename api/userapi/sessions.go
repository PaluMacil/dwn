package userapi

import (
	"encoding/json"

	"github.com/PaluMacil/dwn/dwn"
)

// api/user/sessions
func (rt *UserRoute) handleSessionDetails() {
	if rt.API().ServeCannot(dwn.PermissionViewUsers) {
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
	User    dwn.User    `json:"user"`
	Session dwn.Session `json:"session"`
}
