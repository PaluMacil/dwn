package registration

import (
	"github.com/PaluMacil/dwn/module/core"
	"github.com/gofrs/uuid"
)

type UserCreationRequest struct {
	Email                 string `json:"email"`
	Password              string `json:"password"`
	MustChangePWNextLogin bool   `json:"mustChangePWNextLogin"`
	GivenName             string `json:"givenName"`
	FamilyName            string `json:"familyName"`
}

// Validate examines a request object for creation of a new user and returns a slice
// of validation errors or an empty slice if the request is valid
func (req UserCreationRequest) Validate() []string {
	validationErrors := make([]string, 0)
	return validationErrors
}

func (req UserCreationRequest) User() (core.User, error) {
	user := core.User{
		// TODO: make a new user with an unverified email
		// VerificationCode: uuid.Must(uuid.NewV4()).String(),
	}

	return user, nil
}
