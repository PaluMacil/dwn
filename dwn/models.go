package dwn

import (
	"strings"
	"time"
)

type DisplayName string

func (d DisplayName) Tag() string {
	return "@" + strings.ToLower(strings.Replace(string(d), " ", "", -1))
}

type User struct {
	GoogleID         string      `json:"googleId"`
	GoogleImportDate time.Time   `json:"googleImportDate"`
	Email            string      `json:"email"`
	Tag              string      `json:"tag"`
	PreviousTags     []string    `json:"previousTags"`
	PasswordHash     []byte      `json:"-"`
	VerifiedEmail    bool        `json:"verifiedEmail"`
	Locked           bool        `json:"locked"`
	DisplayName      DisplayName `json:"displayName"`
	GivenName        string      `json:"givenName"`
	FamilyName       string      `json:"familyName"`
	Link             string      `json:"link"`
	Picture          string      `json:"picture"`
	Gender           string      `json:"gender"`
	Locale           string      `json:"locale"`
	LastLogin        time.Time   `json:"lastLogin"`
	ModifiedDate     time.Time   `json:"modifiedDate"`
	CreatedDate      time.Time   `json:"createdDate"`
}

func (u User) Key() []byte {
	return append(u.Prefix(), []byte(u.Email)...)
}

func (u User) Prefix() []byte {
	return []byte(userPrefix)
}
