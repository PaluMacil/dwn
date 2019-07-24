package api

import (
	"encoding/json"
	"fmt"
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
	exists, err := db.Users.VerifiedEmailExists(request.Email)
	if err != nil {
		return err
	}
	// TODO: finish validation
	validationErrors := request.Validate()
	if exists {
		alreadyExistsMessage := fmt.Sprintf("user with email %s already exists", request.Email)
		validationErrors = append(validationErrors, alreadyExistsMessage)
	}
	if len(validationErrors) > 0 {
		// TODO: fix 400 errors to hold more context, such as a string list
		return errs.StatusError{Code: 400, Err: fmt.Errorf("validation errors")}
	}
	userID, err := db.NextID()
	if err != nil {
		return err
	}
	user, err := request.User(userID)
	if err != nil {
		return err
	}
	err = db.Users.Set(user)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(user)
}
