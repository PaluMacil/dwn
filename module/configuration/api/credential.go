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
	"net/url"
	"time"
)

// GET /api/configuration/credential?type=credType&name=credName
func getCredentialHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionViewAppSettings); err != nil {
		return err
	}
	var credential configuration.Credential
	var credentialList []configuration.Credential
	credName, err := url.QueryUnescape(vars["credName"])
	if err != nil {
		return err
	}
	credType := configuration.ForeignSystemType(vars["credType"])

	if credType == "" {
		return errs.StatusError{http.StatusBadRequest, errors.New("credential query must specify type")}
	}

	// handle requests for a list of credentials
	if credName == "" {
		if configuration.ValidForeignSystemTypes.Includes(credType) {
			var err error
			credentialList, err = db.Config.Credential.AllOf(credType)
			if err != nil {
				return err
			}
			return json.NewEncoder(w).Encode(credentialList)
		}
		return errs.StatusNotFound
	}

	// handle requests for a specific credential
	credential, err = db.Config.Credential.Get(credName, credType)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(credential)
}

type CredentialCreationRequest struct {
	Name   string                          `json:"name"`
	Type   configuration.ForeignSystemType `json:"type"`
	ID     string                          `json:"id"`
	Secret string                          `json:"secret"`
}

func (c CredentialCreationRequest) Credential(createdBy store.Identity) configuration.Credential {
	return configuration.Credential{
		Name:        c.Name,
		Type:        c.Type,
		ID:          c.ID,
		Secret:      c.Secret,
		CreatedBy:   createdBy,
		CreatedDate: time.Now(),
	}
}

// POST /api/configuration/credential
func postCredentialHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionChangeAppSettings); err != nil {
		return errs.StatusForbidden
	}
	var request CredentialCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}
	credential := request.Credential(cur.User.ID)

	return db.Config.Credential.Set(credential)
}
