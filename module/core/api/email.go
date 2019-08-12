package api

import (
	"encoding/json"
	"errors"
	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/store"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/webserver/errs"
	"net/http"
	"strings"
	"time"
)

// DELETE /api/core/email?userID=123&email=blah@example.com
func deleteEmailHandler(
	db *database.Database,
	config configuration.Configuration,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	userID, err := store.StringToIdentity(vars["userID"])
	if err != nil {
		return err
	}
	email := strings.ToLower(vars["email"])
	// Must have permission to edit users unless editing oneself
	if err := cur.Can(core.PermissionEditUserInfo); err != nil && cur.User.ID != userID {
		return err
	}
	user, err := db.Users.Get(userID)
	if db.IsKeyNotFoundErr(err) {
		return errs.StatusNotFound
	} else if err != nil {
		return err
	}
	if email == user.PrimaryEmail {
		return errs.StatusError{http.StatusBadRequest, errors.New("cannot delete primary email")}
	}
	updatedEmailList, err := deleteEmail(user.Emails, email)
	if err != nil {
		return err
	}
	user.Emails = updatedEmailList

	return db.Users.Set(user)
}

func deleteEmail(emails []core.Email, email string) ([]core.Email, error) {
	var remainingVerifiedCount int
	var emailRecordForRemoval core.Email
	for i, record := range emails {
		if record.Email == email {
			emails = emails[:i+copy(emails[i:], emails[i+1:])]
			emailRecordForRemoval = record
		} else if record.Verified {
			// count number of verified emails
			remainingVerifiedCount++
		}
	}

	// you can always remove emails that are not verified
	if !emailRecordForRemoval.Verified {
		return emails, nil
	}
	// if this will eliminate the last verified email, you can't delete it
	if remainingVerifiedCount == 0 {
		return nil, errs.StatusError{http.StatusBadRequest, errors.New("cannot delete last verified email")}
	}

	return emails, nil
}

// POST /api/core/email?userID=123&email=blah@example.com&action=something
func emailActionHandler(
	db *database.Database,
	config configuration.Configuration,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	userID, err := store.StringToIdentity(vars["userID"])
	if err != nil {
		return err
	}
	email := strings.ToLower(vars["email"])
	action := vars["action"]
	// Must have permission to edit users unless editing oneself
	if err := cur.Can(core.PermissionEditUserInfo); err != nil && cur.User.ID != userID {
		return err
	}
	user, err := db.Users.Get(userID)
	if db.IsKeyNotFoundErr(err) {
		return errs.StatusNotFound
	} else if err != nil {
		return err
	}

	// process actions
	switch action {
	case "setPrimary":
		if !verifiedEmailExists(user.Emails, email) {
			return errs.StatusError{http.StatusBadRequest, errors.New("primary email must be verified")}
		}
		user.PrimaryEmail = email
		user.ModifiedDate = time.Now()
	case "addEmail":
		exists, err := db.Users.EmailExists(email)
		if err != nil {
			return err
		}
		if exists {
			return errs.StatusError{http.StatusBadRequest, errors.New("email already in use")}
		}
		// TODO: set verification code
		// TODO: send email, once implemented
		record := core.Email{
			Email:                email,
			Verified:             false,
			VerifiedDate:         time.Time{},
			VerificationCode:     "",
			VerificationCodeDate: time.Time{},
		}
		user.Emails = append(user.Emails, record)
	case "resendVerificationMessage":
		// TODO: update code and send email, once implemented
	}
	if err = db.Users.Set(user); err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(user.Info())
}

func verifiedEmailExists(emails []core.Email, email string) bool {
	for _, record := range emails {
		if record.Email == email && record.Verified {
			return true
		}
	}
	return false
}
