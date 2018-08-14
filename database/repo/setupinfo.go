package repo

import (
	"fmt"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/dwn"
)

type SetupInfoRepo struct {
	store dwn.DataStorer
	db    *database.Database
}

func NewSetupInfoRepo(store dwn.DataStorer, db *database.Database) *SetupInfoRepo {
	return &SetupInfoRepo{store, db}
}

func (p SetupInfoRepo) Get() (dwn.SetupInfo, error) {
	var setupInfo = dwn.SetupInfo{}
	item, err := p.store.Get(&setupInfo)
	if err != nil {
		return setupInfo, err
	}
	setupInfo, ok := item.(dwn.SetupInfo)
	if !ok {
		return setupInfo, fmt.Errorf("got data of type %T but wanted dwn.SetupInfo", setupInfo)
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

func (p SetupInfoRepo) Set(setupInfo dwn.SetupInfo) error {
	return p.store.Set(&setupInfo)
}

func (p SetupInfoRepo) Delete() error {
	return p.store.Delete(dwn.SetupInfo{})
}
