package core

import (
	"time"
	"github.com/PaluMacil/dwn/database/store"
)

const GroupPrefix = "GROUP:"

type Group struct {
	Name             string    `json:"name"`
	Permissions      []string  `json:"permissions"`
	Requires2FA      bool      `json:"requires2FA"`
	RequiresVaultPIN bool      `json:"requiresVaultPIN"`
	ModifiedBy       store.Identity    `json:"modifiedBy"`
	ModifiedDate     time.Time `json:"modifiedDate"`
}

type GroupCreationRequest struct {
	Name             string `json:"name"`
	Requires2FA      bool   `json:"requires2FA"`
	RequiresVaultPIN bool   `json:"requiresVaultPIN"`
}

func (req GroupCreationRequest) Group(modifiedBy store.Identity) Group {
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
