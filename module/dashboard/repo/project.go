package repo

import (
	"fmt"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/store"
	"github.com/PaluMacil/dwn/module/dashboard"
)

type ProjectRepo struct {
	store database.Storer
	db    *database.Database
}

func NewProjectRepo(store database.Storer, db *database.Database) *ProjectRepo {
	return &ProjectRepo{store, db}
}

func (p ProjectRepo) Get(id store.Identity) (dashboard.Project, error) {
	var project = dashboard.Project{ID: id}
	item, err := p.store.Get(&project)
	if err != nil {
		return project, err
	}
	project, ok := item.(dashboard.Project)
	if !ok {
		return project, fmt.Errorf("got data of type %T but wanted dashboard.Project", project)
	}
	return project, err
}

// Exists checks whether a specific dashboard.Project exists
func (p ProjectRepo) Exists(id store.Identity) (bool, error) {
	_, err := p.Get(id)
	if p.db.IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

// Set saves a dashboard.Project to the database
func (p ProjectRepo) Set(project dashboard.Project) error {
	return p.store.Set(&project)
}

func (p ProjectRepo) All() ([]dashboard.Project, error) {
	var items []database.Item
	pfx := dashboard.Project{}.Prefix()
	err := p.store.All(pfx, &items, true)
	projects := make([]dashboard.Project, len(items))
	for i, v := range items {
		projects[i] = v.(dashboard.Project)
	}

	return projects, err
}

func (p ProjectRepo) Delete(id store.Identity) error {
	return p.store.Delete(dashboard.Project{ID: id})
}
