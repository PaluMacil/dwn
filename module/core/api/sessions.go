package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/webserver/errs"
)

// POST /api/core/sessions/login
func loginHandler(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	var request core.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}
	ip := core.IP(r)
	userInfo, session, result, err := request.Do(db.Providers, ip)
	// TODO: add minimum wait delay
	switch result {
	case core.LoginResultSuccess:
		groups, err := db.Groups.GroupsFor(userInfo.ID)
		if err != nil {
			return err
		}
		me := core.Me{
			User:    userInfo,
			Session: session,
			Groups:  groups,
		}
		if err := json.NewEncoder(w).Encode(me); err != nil {
			return err
		}
		return nil
	case core.LoginResultBadCredentials:
		return errs.StatusUnauthorized
	case core.LoginResult2FA:
		// TODO: create 2FA
		return nil
	case core.LoginResultChangePassword:
		// TODO: create and return a change password specialized token
		return nil
	case core.LoginResultLockedOrDisabled:
		return errs.StatusLocked
	case core.LoginResultError:
		return err
	}

	return nil
}

// GET /api/core/sessions
func sessionsHandler(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionViewUsers); err != nil {
		return err
	}
	sessions, err := db.Sessions.All()
	if err != nil {
		return err
	}
	details := make([]SessionDetails, len(sessions))
	for i, s := range sessions {
		u, err := db.Users.Get(s.UserID)
		if err != nil {
			return err
		}
		details[i].Session = s
		details[i].User = u
	}
	return json.NewEncoder(w).Encode(details)
}

type SessionDetails struct {
	User    core.User    `json:"user"`
	Session core.Session `json:"session"`
}

// DELETE /api/core/sessions/logout/{token}
func logoutHandler(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if cur.Authenticated() {
		if err := db.Sessions.Delete(vars["token"]); err != nil {
			log.Println("deleting session:", err.Error())
		}
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// GET /api/core/sessions/me
func meHandler(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	groups, err := db.Groups.GroupsFor(cur.User.ID)
	if err != nil {
		return err
	}
	me := core.Me{
		User:    cur.User,
		Session: cur.Session,
		Groups:  groups,
	}
	if err := json.NewEncoder(w).Encode(me); err != nil {
		return err
	}

	return nil
}
