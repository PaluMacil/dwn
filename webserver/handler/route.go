package handler

import (
	"log"
	"net/http"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/configuration"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/webserver/errs"
	"github.com/gorilla/mux"
)

// Option is a type for option constants
type Option int

// Options is a type representing multiple options and providing a Contains method.
type Options []Option

// Contains returns whether an instance of Options contains a specific Option.
func (options Options) Contains(option Option) bool {
	for _, o := range options {
		if o == option {
			return true
		}
	}
	return false
}

const (
	OptionAllowCORS Option = iota
	OptionAllowAnonymous
	OptionVaultPINRequired
)

// Func is the dwn-specific signature for handler functions
type Func func(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error

// Handler is a wrapper struct for giving HandlerFunc functionality
type Handler struct {
	Handler               Func
	db                    *database.Database
	options               Options
	assumeJSONContentType bool
}

// ServeHTTP allows the module Handler to match the interface requirements for the
// standard library Handler while adding application specific context.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the current session information
	token := r.Header.Get("dwn-token")
	cur, err := core.GetCurrent(token, h.db.Providers)
	if err != nil { //TODO: if session not found, 401 Unauthorized is more appropriate
		log.Printf("couldn't get current user with token %s requesting %s: %s", token, r.URL.Path, err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// update heartbeat, logging any error (but don't fail the request if update fails)
	ip := core.IP(r)
	if cur.Authenticated() {
		if err := h.db.Sessions.UpdateHeartbeat(&cur.Session, ip); err != nil {
			log.Printf("updating heartbeat for %s: %s", ip, err)
		}
	}
	// unless anonymous browsing is allowed, check if anonymous and return Unauthorized
	if !h.options.Contains(OptionAllowAnonymous) && cur.Anonymous() {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// for preflight requests, check if CORS is allowed for the route and send appropriate headers
	if h.options.Contains(OptionAllowCORS) && r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		return
	}
	if h.assumeJSONContentType {
		w.Header().Set("Content-Type", "application/json")
	}
	err = h.Handler(h.db, cur, mux.Vars(r), w, r)
	if err != nil {
		switch e := err.(type) {
		case errs.Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// look for database errors and mark these as a 404
			if h.db.IsKeyNotFoundErr(err) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			log.Printf("HTTP 500 - %s", e)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

// Factory stores repeated requirements for constructing handlers
type Factory struct {
	Db                    *database.Database
	Config                configuration.Configuration
	AssumeJSONContentType bool
}

// Handler returns a Handler from a HandlerFunc with database and configuration
func (factory Factory) Handler(handler Func, options ...Option) Handler {
	return Handler{
		Handler:               handler,
		db:                    factory.Db,
		options:               options,
		assumeJSONContentType: factory.AssumeJSONContentType,
	}
}
