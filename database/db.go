package database

import (
	"github.com/PaluMacil/dwn/dwn"
)

func New(store dwn.DataStorer) *Database {
	return &Database{store: store}
}

type Database struct {
	store      dwn.DataStorer
	Sessions   dwn.SessionProvider
	Users      dwn.UserProvider
	Groups     dwn.GroupProvider
	UserGroups dwn.UserGroupProvider
	SetupInfo  dwn.SetupInfoProvider
	Util       dwn.DbUtilityProvider
}

func (db Database) Close() error {
	return db.store.Close()
}
