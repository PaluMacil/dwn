package api

import (
	"log"
	"net/http"
)

func (rt Route) ServeCannot(permissions ...string) bool {
	for _, permission := range permissions {
		if can, err := rt.Current.Can(permission); err != nil {
			http.Error(rt.W, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return true
		} else if !can {
			http.Error(rt.W, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return true
		}
	}
	return false
}

func (rt Route) ServeInternalServerError(err error) {
	log.Println("serving route", rt.Name, err) //TODO: once I fix up stringer for route, update this
	http.Error(rt.W, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	return
}

func (rt Route) ServeBadRequest() {
	//TODO: log once I have a verbose log level
	http.Error(rt.W, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	return
}
