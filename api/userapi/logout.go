package userapi

import "net/http"

// api/user/logout
func (rt *UserRoute) handleLogout() {
	if rt.Current != nil {
		rt.Db.Sessions.Delete(rt.ID)
	}
	rt.W.WriteHeader(http.StatusNoContent)
}
