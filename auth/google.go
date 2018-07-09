package auth

import (
	"time"

	"github.com/PaluMacil/dwn/db"
)

type GoogleClaims struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

func (g *GoogleClaims) CreateUser() db.User {
	return db.User{
		GoogleID:         g.ID,
		GoogleImportDate: time.Now(),
		Email:            g.Email,
		PasswordHash:     []byte{},
		VerifiedEmail:    g.VerifiedEmail,
		Locked:           false,
		DisplayName:      g.GivenName + " " + g.FamilyName,
		GivenName:        g.GivenName,
		FamilyName:       g.FamilyName,
		Link:             g.Link,
		Picture:          g.Picture,
		Gender:           g.Gender,
		Locale:           g.Locale,
		LastLogin:        time.Now(),
		ModifiedDate:     time.Now(),
		CreatedDate:      time.Now(),
	}
}