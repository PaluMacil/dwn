package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/configuration"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/webserver/errs"
)

// GET /api/typeahead/users&query={query}
func usersHandler(
	db *database.Database,
	config configuration.Configuration,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionViewUsers); err != nil {
		return err
	}
	qry, err := url.QueryUnescape(vars["query"])
	if len(qry) < 2 || err != nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("invalid user query")}
	}

	user, err := db.Users.CompletionSuggestions(qry) //TODO: check for absence of @ and search by username
	if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(core.Users(user).Info()); err != nil {
		return err
	}

	return nil
}
