package core

import (
	"strings"
	"time"
)

const (
	SessionPrefix   = "SESSION:"
	UserPrefix      = "USER:"
	GroupPrefix     = "GROUP:"
	UserGroupPrefix = "USERGROUP:"
	SetupInfoPrefix = "SETUPINFO:"
)

type DisplayName string

func (d DisplayName) Tag() string {
	return "@" + strings.ToLower(strings.Replace(string(d), " ", "", -1))
}

type User struct {
	GoogleID              string      `json:"googleId"`
	GoogleImportDate      time.Time   `json:"googleImportDate"`
	Email                 string      `json:"email"`
	Tag                   string      `json:"tag"`
	PreviousTags          []string    `json:"previousTags"`
	PasswordHash          []byte      `json:"-"`
	LastPasswordHash      []byte      `json:"-"`
	MustChangePWNextLogin bool        `json:"mustChangePWNextLogin"`
	Require2FA            bool        `json:"require2FA"`
	VerifiedEmail         bool        `json:"verifiedEmail"`
	VerifiedEmailDate     time.Time   `json:"verifiedEmailDate"`
	VerificationCode      string      `json:"-"`
	VaultPIN              string      `json:"-"`
	Locked                bool        `json:"locked"`
	LoginAttempts         int         `json:"loginAttempts"`
	LastFailedLogin       time.Time   `json:"lastFailedLogin"`
	DisplayName           DisplayName `json:"displayName"`
	GivenName             string      `json:"givenName"`
	FamilyName            string      `json:"familyName"`
	Link                  string      `json:"link"`
	Picture               string      `json:"picture"`
	Gender                string      `json:"gender"`
	Locale                string      `json:"locale"`
	LastLogin             time.Time   `json:"lastLogin"`
	ModifiedDate          time.Time   `json:"modifiedDate"`
	CreatedDate           time.Time   `json:"createdDate"`
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
	userInfo := make([]UserInfo, len(users), len(users))
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
	return append(u.Prefix(), []byte(u.Email)...)
}

func (u User) Prefix() []byte {
	return []byte(UserPrefix)
}

type Group struct {
	Name             string    `json:"name"`
	Permissions      []string  `json:"permissions"`
	Requires2FA      bool      `json:"requires2FA"`
	RequiresVaultPIN bool      `json:"requiresVaultPIN"`
	ModifiedBy       string    `json:"modifiedBy"`
	ModifiedDate     time.Time `json:"modifiedDate"`
}

type GroupCreationRequest struct {
	Name             string `json:"name"`
	Requires2FA      bool   `json:"requires2FA"`
	RequiresVaultPIN bool   `json:"requiresVaultPIN"`
}

func (req GroupCreationRequest) Group(modifiedBy string) Group {
	return Group{
		Name:             req.Name,
		Permissions:      []string{},
		Requires2FA:      req.Requires2FA,
		RequiresVaultPIN: req.RequiresVaultPIN,
		ModifiedBy:       modifiedBy,
		ModifiedDate:     time.Now(),
	}
}

func (g Group) HasPermission(permission string) bool {
	for _, p := range g.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

const (
	BuiltInGroupAdmin     = "ADMIN"
	BuiltInGroupSpouse    = "SPOUSE"
	BuiltInGroupResident  = "RESIDENT"
	BuiltInGroupFriend    = "FRIEND"
	BuiltInGroupLandlord  = "LANDLORD"
	BuiltInGroupTenant    = "TENANT"
	BuiltInGroupPowerUser = "POWER_USER"
	BuiltInGroupUser      = "USER"
)

func (g Group) Key() []byte {
	return append(g.Prefix(), []byte(g.Name)...)
}

func (g Group) Prefix() []byte {
	return []byte(GroupPrefix)
}

type Session struct {
	Token       string    `json:"token"`
	Email       string    `json:"email"`
	IP          string    `json:"ip"`
	CreatedDate time.Time `json:"createdDate"`
	Heartbeat   time.Time `json:"heartbeat"`
}

func (s Session) Key() []byte {
	return append(s.Prefix(), []byte(s.Token)...)
}

func (s Session) Prefix() []byte {
	return []byte(SessionPrefix)
}

type SetupInfo struct {
	InitializedDate time.Time `json:"initializedDate"`
}

func (s SetupInfo) Key() []byte {
	return s.Prefix()
}

func (s SetupInfo) Prefix() []byte {
	return []byte(SetupInfoPrefix)
}

type UserGroup struct {
	Email       string    `json:"email"`
	GroupName   string    `json:"groupName"`
	CreatedDate time.Time `json:"createdDate"`
}

func (u UserGroup) Key() []byte {
	partOne := append(u.Prefix(), []byte(u.Email)...)
	return append(partOne, []byte(u.GroupName)...)
}

func (u UserGroup) Prefix() []byte {
	return []byte(UserGroupPrefix)
}
