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
	log.Printf("serving route %s: %s", rt.Name, err) // TODO: once I fix up stringer for route, update this
	http.Error(rt.W, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	return
}

func (rt Route) ServeBadRequest() {
	// TODO: log once I have a verbose log level
	http.Error(rt.W, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	return
}

func (rt Route) ServeMethodNotAllowed() {
	// TODO: log once I have a verbose log level (or maybe for this it is an error since it represents
	// code calling wrong method, not improper query params).
	http.Error(rt.W, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	return
}
