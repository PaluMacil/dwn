package repo

import "github.com/PaluMacil/dwn/database"

type CredentialRepo struct {
	store database.Storer
	db    *database.Database
}

func NewCredentialRepo(store database.Storer, db *database.Database) *CredentialRepo {
	return &CredentialRepo{store, db}
}
