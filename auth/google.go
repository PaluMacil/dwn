package auth

import (
	"time"

	"github.com/PaluMacil/dwn/core"
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

func (g *GoogleClaims) CreateUser(displayName core.DisplayName) core.User {
	// TODO: must check if displayname as tag exists once bleve search is up
	return core.User{
		GoogleID:         g.ID,
		GoogleImportDate: time.Now(),
		Email:            g.Email,
		Tag:              displayName.Tag(),
		PreviousTags:     []string{},
		PasswordHash:     []byte{},
		VerifiedEmail:    g.VerifiedEmail,
		Locked:           false,
		DisplayName:      displayName,
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
