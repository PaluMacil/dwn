package repo

import (
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/dashboard"
)

type DashboardRepo struct {
	store database.Storer
	db    *database.Database
}

func NewDashboardRepo(store database.Storer, db *database.Database) *DashboardRepo {
	return &DashboardRepo{store, db}
}

func (r DashboardRepo) Get() (dashboard.Dashboard, error) {
	projects, err := r.db.Dashboard.Projects.All()
	return dashboard.Dashboard{Projects: projects}, err
}
