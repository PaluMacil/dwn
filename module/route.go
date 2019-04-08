package module

import (
	"net/http"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/gorilla/mux"
)

// HandlerOption is a type for option constands
type HandlerOption int

// HandlerOptions is a type representing multiple options and providing a Contains method.
type HandlerOptions []HandlerOption

// Contains returns whether an instance of HandlerOptions containts a specific option.
func (options HandlerOptions) Contains(option HandlerOption) bool {
	for _, o := range options {
		if o == option {
			return true
		}
	}
	return false
}

const (
	OptionAllowCORS HandlerOption = iota
	OptionAllowAnonymous
	OptionVaultPINRequired
)

// HandlerFunc is the dwn-specific signature for handler functions
type HandlerFunc func(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
)

// Handler is a wrapper struct for giving HandlerFunc functionality
type Handler struct {
	Handler               HandlerFunc
	db                    *database.Database
	config                configuration.Configuration
	options               HandlerOptions
	assumeJSONContentType bool
}

// ServeHTTP allows the module Handler to match the interface requirements for the
// standard library Handler while adding application specific context.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the current session information
	cur, err := core.GetCurrent(r, h.db.Providers)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// unless anonymous browsing is allowed, check for authentication
	if !h.options.Contains(OptionAllowAnonymous) && !cur.Authenticated() {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// for preflight requests, check if CORS is allowed for the route and send appropriate headers
	if h.options.Contains(OptionAllowCORS) && r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		return
	}
	if h.assumeJSONContentType {
		w.Header().Set("Content-Type", "application/json")
	}
	h.Handler(h.db, h.config, cur, mux.Vars(r), w, r)
}

// HandlerFactory stores repeated requirements for constructing handlers
type HandlerFactory struct {
	Db                    *database.Database
	Config                configuration.Configuration
	AssumeJSONContentType bool
}

// Handler returns a Handler from a HandlerFunc with database and configuration
func (factory HandlerFactory) Handler(handler HandlerFunc, options ...HandlerOption) Handler {
	return Handler{
		Handler:               handler,
		db:                    factory.Db,
		config:                factory.Config,
		options:               HandlerOptions(options),
		assumeJSONContentType: factory.AssumeJSONContentType,
	}
}
