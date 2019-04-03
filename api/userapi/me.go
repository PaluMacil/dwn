package userapi

import (
	"encoding/json"
	"net/http"

	"github.com/PaluMacil/dwn/auth"
	"github.com/PaluMacil/dwn/core"
)

// api/user/me
func (rt *UserRoute) handleMe() {
	if rt.Current == nil {
		http.Error(rt.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	groups, err := rt.Db.Groups.GroupsFor(rt.Current.User.Email)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	me := struct {
		auth.Current
		Groups []core.Group `json:"groups"`
	}{
		*rt.Current,
		groups,
	}
	if err := json.NewEncoder(rt.W).Encode(me); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
}
