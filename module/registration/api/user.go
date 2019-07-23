package api

import (
	"encoding/json"
	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/module/registration"
	"github.com/PaluMacil/dwn/webserver/errs"
	"net/http"
)

// POST /api/registration/user
func userHandler(
	db *database.Database,
	config configuration.Configuration,
	cur *core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	var request registration.UserCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}
	// TODO: check if email already exists on a validated user
	// TODO: finish validation
	validationErrors := request.Validate()
	if len(validationErrors) > 0 {
		return errs.StatusError{}
	}
	user, err := db.Users.Get(id)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(user)
}
