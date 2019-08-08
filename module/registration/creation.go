package registration

import (
	"fmt"
	"github.com/PaluMacil/dwn/database/store"
	"github.com/PaluMacil/dwn/module/core"
	"time"
)

type UserCreationRequest struct {
	Email                      string `json:"email"`
	Password                   string `json:"password"`
	WithVerifiedEmail          bool   `json:"withVerifiedEmail"`
	IgnorePasswordRequirements bool   `json:"ignorePasswordRequirements"`
	MustChangePWNextLogin      bool   `json:"mustChangePWNextLogin"`
	GivenName                  string `json:"givenName"`
	FamilyName                 string `json:"familyName"`
}

// Validate examines a request object for creation of a new user and returns a slice
// of validation errors or an empty slice if the request is valid
func (req UserCreationRequest) Validate() []string {
	validationErrors := make([]string, 0)
	//TODO: check for restricted fields, bad pw choice, cetc
	return validationErrors
}

func (req UserCreationRequest) User(id store.Identity) (core.User, error) {
	displayName := core.DisplayName(fmt.Sprintf("%s %s", req.GivenName, req.FamilyName))
	passwordHash, err := core.CreateHash(req.Password)
	if err != nil {
		return core.User{}, err
	}
	// TODO: send verification email
	user := core.User{
		ID:           id,
		PrimaryEmail: req.Email,
		Emails: []core.Email{
			{
				Email:    req.Email,
				Verified: false,
				// TODO: verified code and date
				// VerificationCode: uuid.Must(uuid.NewV4()).String(),
			},
		},
		//TODO: tag
		DisplayName:  displayName,
		PasswordHash: passwordHash,
		GivenName:    req.GivenName,
		FamilyName:   req.FamilyName,
		ModifiedDate: time.Now(),
		CreatedDate:  time.Now(),
	}

	return user, nil
}
