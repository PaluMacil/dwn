package core

import (
	"strings"
	"time"

	"github.com/PaluMacil/dwn/database/store"
	"golang.org/x/crypto/bcrypt"
)

const UserPrefix = "USER:"

type DisplayName string

func (d DisplayName) Tag() string {
	return "@" + strings.ToLower(strings.Replace(string(d), " ", "", -1))
}

type Hash []byte

func CreateHash(password string) (Hash, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return Hash(bytes), err
}

func (hash Hash) Check(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Email struct {
	Email                string    `json:"email"`
	Verified             bool      `json:"verified"`
	VerifiedDate         time.Time `json:"verifiedDate"`
	VerificationCode     string    `json:"-"`
	VerificationCodeDate time.Time `json:"-"`
}

type User struct {
	ID                    store.Identity `json:"id"`
	GoogleID              string         `json:"googleId"`
	GoogleImportDate      time.Time      `json:"googleImportDate"`
	PrimaryEmail          string         `json:"primaryEmail"`
	Emails                []Email        `json:"emails"`
	Tag                   string         `json:"tag"`
	PreviousTags          []string       `json:"previousTags"`
	PasswordHash          Hash           `json:"-"`
	LastPasswordHash      Hash           `json:"-"`
	MustChangePWNextLogin bool           `json:"mustChangePWNextLogin"`
	Require2FA            bool           `json:"require2FA"`
	VaultPIN              string         `json:"-"`
	Locked                bool           `json:"locked"`
	Disabled              bool           `json:"disabled"`
	LoginAttempts         int            `json:"loginAttempts"`
	LastFailedLogin       time.Time      `json:"lastFailedLogin"`
	DisplayName           DisplayName    `json:"displayName"`
	GivenName             string         `json:"givenName"`
	FamilyName            string         `json:"familyName"`
	Link                  string         `json:"link"`
	Picture               string         `json:"picture"`
	Gender                string         `json:"gender"`
	Locale                string         `json:"locale"`
	LastLogin             time.Time      `json:"lastLogin"`
	ModifiedDate          time.Time      `json:"modifiedDate"`
	CreatedDate           time.Time      `json:"createdDate"`
}

func (u User) Info() UserInfo {
	return UserInfo{
		User:        u,
		HasPassword: len(u.PasswordHash) > 0,
		HasVaultPIN: u.VaultPIN != "",
	}
}

type Users []User

func (users Users) Info() []UserInfo {
	userInfo := make([]UserInfo, len(users))
	for i, user := range users {
		userInfo[i] = user.Info()
	}
	return userInfo
}

type UserInfo struct {
	User
	HasPassword bool `json:"hasPassword"`
	HasVaultPIN bool `json:"hasVaultPIN"`
}

func (u User) Key() []byte {
	return append(u.Prefix(), u.ID.Bytes()...)
}

func (u User) Prefix() []byte {
	return []byte(UserPrefix)
}

func (u User) CanLogin() bool {
	return !u.Locked && !u.Disabled
}
