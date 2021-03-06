package repo

import (
	"errors"
	"fmt"
	"github.com/PaluMacil/dwn/module/core"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/store"
)

type GroupRepo struct {
	store database.Storer
	db    *database.Database
}

func NewGroupRepo(store database.Storer, db *database.Database) *GroupRepo {
	return &GroupRepo{store, db}
}

func (p GroupRepo) GroupsFor(userID store.Identity) ([]core.Group, error) {
	var items []database.Item
	extendedPrefix := append(core.UserGroup{}.Prefix(), userID.Bytes()...)
	err := p.store.All(extendedPrefix, &items, true)
	if err != nil {
		return nil, err
	}
	userGroups := make([]core.UserGroup, len(items))
	for i, v := range items {
		userGroups[i] = v.(core.UserGroup)
	}

	groups := make([]core.Group, len(items))
	for i, ug := range userGroups {
		groups[i], err = p.Get(ug.GroupName)
		if err != nil {
			return groups, err
		}
	}

	return groups, err
}

func (p GroupRepo) Get(name string) (core.Group, error) {
	var group = core.Group{Name: name}
	if name == "" {
		return group, errors.New("GroupRepo.Get requires a name but got an empty string")
	}
	item, err := p.store.Get(&group)
	if err != nil {
		return group, err
	}
	group, ok := item.(core.Group)
	if !ok {
		return group, fmt.Errorf("got data of type %T but wanted Group", group)
	}
	return group, err
}

func (p GroupRepo) Exists(name string) (bool, error) {
	_, err := p.Get(name)
	if p.db.IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

func (p GroupRepo) Set(group core.Group) error {
	return p.store.Set(&group)
}

func (p GroupRepo) All() ([]core.Group, error) {
	var items []database.Item
	err := p.store.All(core.Group{}.Prefix(), &items, true)
	groups := make([]core.Group, len(items))
	for i, v := range items {
		groups[i] = v.(core.Group)
	}

	return groups, err
}

func (p GroupRepo) Delete(name string) error {
	return p.store.Delete(core.Group{Name: name})
}
