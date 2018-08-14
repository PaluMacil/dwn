package dwn

import (
	"strings"
	"time"
)

const (
	SessionPrefix    = "SESSION:"
	UserPrefix       = "USER:"
	GroupPrefix      = "GROUP:"
	PermissionPrefix = "PERMISSION:"
	UserGroupPrefix  = "USERGROUP:"
	SetupInfoPrefix  = "SETUPINFO:"
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
	return []byte(UserPrefix)
}

type Group struct {
	Name         string    `json:"name"`
	Permissions  []string  `json:"permissions"`
	ModifiedBy   string    `json:"modifiedBy"`
	ModifiedDate time.Time `json:"modifiedDate"`
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
	BuiltInGroupAdmin    = "ADMIN"
	BuiltInGroupSpouse   = "SPOUSE"
	BuiltInGroupResident = "RESIDENT"
	BuiltInGroupFriend   = "FRIEND"
	BuiltInGroupLandlord = "LANDLORD"
	BuiltInGroupTenant   = "TENANT"
	BuiltInGroupUser     = "USER"
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
