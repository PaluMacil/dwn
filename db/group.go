package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
)

type Group struct {
	Name         string
	Permissions  []string
	ModifiedBy   string
	ModifiedDate time.Time
}

const (
	BuiltInGroupAdmin    = "ADMIN"
	BuiltInGroupSpouse   = "SPOUSE"
	BuiltInGroupResident = "RESIDENT"
	BuiltInGroupFriend   = "FRIEND"
	BuiltInGroupTenant   = "TENANT"
	BuiltInGroupUser     = "USER"
)

func (g Group) Key() []byte {
	return append(g.Prefix(), []byte(g.Name)...)
}

func (g Group) Prefix() []byte {
	return []byte(groupPrefix)
}

/*
func (p *GroupProvider) GroupsFor(email string) ([]Group, error) {
	var items []DbItem
	err := all(p.bgr, Group{}.Prefix(), &items, true)
	groups := make([]Group, len(items))
	for i, v := range items {
		groups[i] = v.(Group)
	}

	return groups, err
}*/

type GroupProvider struct {
	bgr *badger.DB
}

func (p *GroupProvider) Get(name string) (Group, error) {
	var group = Group{Name: name}
	if name == "" {
		return group, errors.New("GroupProvider.Get requires a name but got an empty string")
	}
	item, err := get(p.bgr, &group)
	if err != nil {
		return group, err
	}
	group, ok := item.(Group)
	if !ok {
		return group, fmt.Errorf("got data of type %T but wanted Group", group)
	}
	return group, err
}

func (p *GroupProvider) Exists(name string) (bool, error) {
	_, err := p.Get(name)
	return err != badger.ErrKeyNotFound, err
}

func (p *GroupProvider) Set(group Group) error {
	return set(p.bgr, &group)
}

func (p *GroupProvider) All() ([]Group, error) {
	var items []DbItem
	err := all(p.bgr, Group{}.Prefix(), &items, true)
	groups := make([]Group, len(items))
	for i, v := range items {
		groups[i] = v.(Group)
	}

	return groups, err
}

func (p GroupProvider) Delete(name string) error {
	return delete(p.bgr, Group{Name: name})
}
