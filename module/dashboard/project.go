package dashboard

import (
	"github.com/PaluMacil/dwn/database/store"
)

const ProjectPrefix = "DASHBOARD:PROJECT:"

type Project struct {
	ID          store.Identity `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}

func (p Project) Key() []byte {
	return append(p.Prefix(), p.ID.Bytes()...)
}

func (p Project) Prefix() []byte {
	return []byte(ProjectPrefix)
}
