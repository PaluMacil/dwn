package api

import (
	"encoding/json"
	"fmt"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/module/registration"
	"github.com/PaluMacil/dwn/webserver/errs"
	"log"
	"net/http"
	"time"
)

// POST /api/registration/user
func userHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	var request registration.UserCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return fmt.Errorf("decoding user creation request: %w", err)
	}
	exists, err := db.Users.VerifiedEmailExists(request.Email)
	if err != nil {
		return fmt.Errorf("checking if email %s exists and is verified: %w", request.Email, err)
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
		return fmt.Errorf("getting next ID for user creation request: %w", err)
	}
	user, err := request.User(userID)
	if err != nil {
		return fmt.Errorf("converting user creation request into a user: %w", err)
	}
	config := db.Config.Get()
	err = db.Users.Set(user)
	if err != nil {
		return fmt.Errorf("saving newly created user: %w", err)
	}
	// if this is the initial admin, add it to the admin group; email and password must both match the one configured
	if config.Setup.InitialAdmin != "" && config.Setup.InitialAdmin == request.Email && config.Setup.InitialPassword == request.Password {
		err = db.UserGroups.Set(core.UserGroup{
			UserID:    userID,
			GroupName: core.BuiltInGroupAdmin,
		})
		if err != nil {
			return fmt.Errorf("saving initial admin user to admin group: %w", err)
		}
		user.Emails[0].Verified = true
		user.Emails[0].VerifiedDate = time.Now()
		err = db.Users.Set(user)
		if err != nil {
			return fmt.Errorf("updating admin user to verify email: %w", err)
		}
		log.Printf("user '%s' (%s) created as initial admin", user.DisplayName, user.PrimaryEmail)
	}

	return json.NewEncoder(w).Encode(user.Info())
}
