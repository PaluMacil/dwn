package oauth

import (
	"time"

	"github.com/PaluMacil/dwn/database/store"
	"github.com/PaluMacil/dwn/module/core"
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

func (g *GoogleClaims) CreateUser(id store.Identity, displayName core.DisplayName) core.User {
	// TODO: must check if displayname as tag exists once bleve search is up
	return core.User{
		ID:               id,
		GoogleID:         g.ID,
		GoogleImportDate: time.Now(),
		PrimaryEmail:     g.Email,
		Emails: []core.Email{
			core.Email{
				Email:        g.Email,
				Verified:     g.VerifiedEmail,
				VerifiedDate: time.Now(),
			},
		},
		Tag:          displayName.Tag(),
		PreviousTags: []string{},
		PasswordHash: []byte{},
		Locked:       false,
		DisplayName:  displayName,
		GivenName:    g.GivenName,
		FamilyName:   g.FamilyName,
		Link:         g.Link,
		Picture:      g.Picture,
		Gender:       g.Gender,
		Locale:       g.Locale,
		LastLogin:    time.Now(),
		ModifiedDate: time.Now(),
		CreatedDate:  time.Now(),
	}
}
