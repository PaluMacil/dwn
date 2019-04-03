package repo

import (
	"github.com/PaluMacil/dwn/database"
)

type EntryRepo struct {
	store database.Storer
	db    *database.Database
}

//func NewEntryRepo(store database.Storer, db *database.Database) *EntryRepo {
//	entryRepo := &EntryRepo{store, db}
//	log.SetOutput(entryRepo)
//	return entryRepo
//}

func (r EntryRepo) Write(p []byte) (n int, err error) {
	return 0, nil
}
