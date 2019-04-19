package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/webserver/errs"
	"github.com/PaluMacil/dwn/module/core"
)

// GET /api/typeahead/users&query={query}
func usersHandler(
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
	qry, err := url.QueryUnescape(vars["query"])
	if len(qry) < 2 || err != nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("invalid user query")}
	}

	user, err := db.Users.CompletionSuggestions(qry) //TODO: check for absence of @ and search by username
	if err != nil {
		return err
	}
	if err := json.NewEncoder(rt.W).Encode(user); err != nil {
		return err
	}
}