package repo

import (
	"errors"
	"fmt"
	"time"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/dwn"
)

type UserGroupRepo struct {
	store database.Storer
	db    *database.Database
}

func NewUserGroupRepo(store database.Storer, db *database.Database) *UserGroupRepo {
	return &UserGroupRepo{store, db}
}

func (p UserGroupRepo) Get(email, groupName string) (dwn.UserGroup, error) {
	var userGroup = dwn.UserGroup{
		Email:     email,
		GroupName: groupName,
	}
	if email == "" {
		return userGroup, errors.New("UserGroupRepo.Get requires an email but got an empty string")
	}
	if groupName == "" {
		return userGroup, errors.New("UserGroupRepo.Get requires a groupName but got an empty string")
	}
	item, err := p.store.Get(&userGroup)
	if err != nil {
		return userGroup, err
	}
	userGroup, ok := item.(dwn.UserGroup)
	if !ok {
		return userGroup, fmt.Errorf("got data of type %T but wanted UserGroup", userGroup)
	}
	return userGroup, err
}

func (p UserGroupRepo) Exists(email, groupName string) (bool, error) {
	_, err := p.Get(email, groupName)
	if p.db.Util.IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

func (p UserGroupRepo) Set(userGroup dwn.UserGroup) error {
	if userGroup.CreatedDate.IsZero() {
		userGroup.CreatedDate = time.Now()
	}
	return p.store.Set(&userGroup)
}

func (p UserGroupRepo) All() ([]dwn.UserGroup, error) {
	var items []database.Item
	err := p.store.All(dwn.UserGroup{}.Prefix(), &items, true)
	userGroups := make([]dwn.UserGroup, len(items))
	for i, v := range items {
		userGroups[i] = v.(dwn.UserGroup)
	}

	return userGroups, err
}

func (p UserGroupRepo) Delete(email, groupName string) error {
	return p.store.Delete(dwn.UserGroup{
		Email:     email,
		GroupName: groupName,
	})
}
