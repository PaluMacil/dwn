package dashboard

import (
	"github.com/PaluMacil/dwn/database/store"
)

const ProjectPrefix = "DASHBOARD:PROJECT:"

type Project struct {
	ID          store.Identity `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Links       []Link         `json:"links"`
	// Owners, Participants, both as slices of user id (email?)
}

type Link struct {
	Text string `json:"name"`
	Order int `json:"order"`
	Audience UserType `json:"audience"`
	External bool `json:"external"`
	Ref string `json:"ref"`
}

type UserType string
const (
	UserTypeOwner UserType = "OWNER"
	UserTypeParticipant UserType = "USER"
	UserTypePublic UserType = "PUBLIC"
)

func (p Project) Key() []byte {
	return append(p.Prefix(), p.ID.Bytes()...)
}

func (p Project) Prefix() []byte {
	return []byte(ProjectPrefix)
}
