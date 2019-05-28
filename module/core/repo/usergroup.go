package repo

import (
	"errors"
	"fmt"
	"time"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/store"
	"github.com/PaluMacil/dwn/module/core"
)

type UserGroupRepo struct {
	store database.Storer
	db    *database.Database
}

func NewUserGroupRepo(store database.Storer, db *database.Database) *UserGroupRepo {
	return &UserGroupRepo{store, db}
}

func (p UserGroupRepo) Get(userID store.Identity, groupName string) (core.UserGroup, error) {
	var userGroup = core.UserGroup{
		UserID:    userID,
		GroupName: groupName,
	}
	if groupName == "" {
		return userGroup, errors.New("UserGroupRepo.Get requires a groupName but got an empty string")
	}
	item, err := p.store.Get(&userGroup)
	if err != nil {
		return userGroup, err
	}
	userGroup, ok := item.(core.UserGroup)
	if !ok {
		return userGroup, fmt.Errorf("got data of type %T but wanted UserGroup", userGroup)
	}
	return userGroup, err
}

func (p UserGroupRepo) Exists(userID, groupName string) (bool, error) {
	_, err := p.Get(userID, groupName)
	if p.db.IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

func (p UserGroupRepo) Set(userGroup core.UserGroup) error {
	if userGroup.CreatedDate.IsZero() {
		userGroup.CreatedDate = time.Now()
	}
	return p.store.Set(&userGroup)
}

func (p UserGroupRepo) All() ([]core.UserGroup, error) {
	var items []database.Item
	err := p.store.All(core.UserGroup{}.Prefix(), &items, true)
	userGroups := make([]core.UserGroup, len(items))
	for i, v := range items {
		userGroups[i] = v.(core.UserGroup)
	}

	return userGroups, err
}

func (p UserGroupRepo) Delete(userID store.Identity, groupName string) error {
	return p.store.Delete(core.UserGroup{
		UserID:    userID,
		GroupName: groupName,
	})
}
