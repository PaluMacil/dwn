package userapi

import "net/http"

// api/user/logout
func (rt *UserRoute) handleLogout() {
	// TODO: ensure post or delete (and sync with frontend)
	if rt.Current != nil {
		rt.Db.Sessions.Delete(rt.ID)
	}
	rt.W.WriteHeader(http.StatusNoContent)
}
