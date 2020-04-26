package api

import (
	"encoding/json"
	"errors"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/store"
	"github.com/PaluMacil/dwn/module/configuration"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/webserver/errs"
	"net/http"
	"time"
)

// GET /api/registration/verify/{verificationMessage}
func verifyHandler(
	db *database.Database,
	config configuration.Configuration,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	// TODO: Unwrap message, verify authenticity, and set specified email as verified
	return nil
}

type VerificationRequest struct {
	UserID store.Identity `json:"userID"`
	Email  string         `json:"email"`
}

// POST /api/registration/verify
func adminVerifyHandler(
	db *database.Database,
	config configuration.Configuration,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionEditUserInfo); err != nil {
		return err
	}
	if r.Body == nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("no request body")}
	}
	var request VerificationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return err
	}

	user, err := db.Users.Get(request.UserID)
	if err != nil {
		return err
	}
	userHasEmail := false
	for i, email := range user.Emails {
		if email.Email == request.Email {
			userHasEmail = true
			// set the email of the struct in the original user, not the loop instance copy
			user.Emails[i].Verified = true
			user.Emails[i].VerifiedDate = time.Now()
			user.Emails[i].VerificationCode = ""
			user.Emails[i].VerificationCodeDate = time.Time{}
			if err := db.Users.Set(user); err != nil {
				return err
			}
		}
	}
	if !userHasEmail {
		return errs.StatusNotFound
	}

	return json.NewEncoder(w).Encode(user.Info())
}
