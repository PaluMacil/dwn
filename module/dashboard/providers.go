package dashboard

import (
	"github.com/PaluMacil/dwn/database/store"
)

type Providers struct {
	Board    DashboardProvider
	Projects ProjectProvider
}

type DashboardProvider interface {
	Get() (Dashboard, error)
}

type ProjectProvider interface {
	Get(id store.Identity) (Project, error)
	Exists(id store.Identity) (bool, error)
	Set(project Project) error
	All() ([]Project, error)
	Delete(id store.Identity) error
}
