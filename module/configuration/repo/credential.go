package repo

import (
	"fmt"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/configuration"
)

type CredentialRepo struct {
	store database.Storer
	db    *database.Database
}

func NewCredentialRepo(store database.Storer, db *database.Database) *CredentialRepo {
	return &CredentialRepo{store, db}
}

func (c CredentialRepo) Get(name string, fsType configuration.ForeignSystemType) (configuration.Credential, error) {
	cred := configuration.Credential{
		Name: name,
		Type: fsType,
	}
	item, err := c.store.Get(&cred)
	if err != nil {
		return cred, err
	}
	cred, ok := item.(configuration.Credential)
	if !ok {
		return cred, fmt.Errorf("got data of type %T but wanted Credential", cred)
	}
	return cred, nil
}

func (c CredentialRepo) Exists(name string, fsType configuration.ForeignSystemType) (bool, error) {
	_, err := c.Get(name, fsType)
	if c.db.IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

func (c CredentialRepo) Set(cred configuration.Credential) error {
	return c.store.Set(&cred)
}

func (c CredentialRepo) AllOf(fsType configuration.ForeignSystemType) ([]configuration.Credential, error) {
	var items []database.Item
	prefix := append(configuration.Credential{}.Prefix(), fsType.Bytes()...)
	prefix = append(prefix, []byte(":")[0])
	err := c.store.All(prefix, &items, true)
	users := make([]configuration.Credential, len(items))
	for i, v := range items {
		users[i] = v.(configuration.Credential)
	}

	return users, err
}

func (c CredentialRepo) Delete(name string, fsType configuration.ForeignSystemType) error {
	cred := configuration.Credential{
		Name: name,
		Type: fsType,
	}
	return c.store.Delete(cred)
}
