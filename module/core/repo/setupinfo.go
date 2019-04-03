package repo

import (
	"fmt"
	"github.com/PaluMacil/dwn/module/core"

	"github.com/PaluMacil/dwn/database"
)

type SetupInfoRepo struct {
	store database.Storer
	db    *database.Database
}

func NewSetupInfoRepo(store database.Storer, db *database.Database) *SetupInfoRepo {
	return &SetupInfoRepo{store, db}
}

func (p SetupInfoRepo) Get() (core.SetupInfo, error) {
	var setupInfo = core.SetupInfo{}
	item, err := p.store.Get(&setupInfo)
	if err != nil {
		return setupInfo, err
	}
	setupInfo, ok := item.(core.SetupInfo)
	if !ok {
		return setupInfo, fmt.Errorf("got data of type %T but wanted core.SetupInfo", setupInfo)
	}
	return setupInfo, err
}

func (p SetupInfoRepo) Completed() (bool, error) {
	_, err := p.Get()
	if p.db.Util.IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

func (p SetupInfoRepo) Set(setupInfo core.SetupInfo) error {
	return p.store.Set(&setupInfo)
}

func (p SetupInfoRepo) Delete() error {
	return p.store.Delete(core.SetupInfo{})
}
