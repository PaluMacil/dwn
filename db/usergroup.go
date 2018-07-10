package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
)

type UserGroup struct {
	Email       string
	GroupName   string
	CreatedDate time.Time
}

func (u UserGroup) Key() []byte {
	partOne := append(u.Prefix(), []byte(u.Email)...)
	return append(partOne, []byte(u.GroupName)...)
}

func (u UserGroup) Prefix() []byte {
	return []byte(userGroupPrefix)
}

type UserGroupProvider struct {
	bgr *badger.DB
	Db  *Db
}

func (p *UserGroupProvider) Get(email, groupName string) (UserGroup, error) {
	var userGroup = UserGroup{
		Email:     email,
		GroupName: groupName,
	}
	if email == "" {
		return userGroup, errors.New("UserGroupProvider.Get requires an email but got an empty string")
	}
	if groupName == "" {
		return userGroup, errors.New("UserGroupProvider.Get requires a groupName but got an empty string")
	}
	item, err := get(p.bgr, &userGroup)
	if err != nil {
		return userGroup, err
	}
	userGroup, ok := item.(UserGroup)
	if !ok {
		return userGroup, fmt.Errorf("got data of type %T but wanted UserGroup", userGroup)
	}
	return userGroup, err
}

func (p *UserGroupProvider) Exists(email, groupName string) (bool, error) {
	_, err := p.Get(email, groupName)
	return err != badger.ErrKeyNotFound, err
}

func (p *UserGroupProvider) Set(userGroup UserGroup) error {
	if userGroup.CreatedDate.IsZero() {
		userGroup.CreatedDate = time.Now()
	}
	return set(p.bgr, &userGroup)
}

func (p *UserGroupProvider) All() ([]UserGroup, error) {
	var items []DbItem
	err := all(p.bgr, UserGroup{}.Prefix(), &items, true)
	userGroups := make([]UserGroup, len(items))
	for i, v := range items {
		userGroups[i] = v.(UserGroup)
	}

	return userGroups, err
}

func (p UserGroupProvider) Delete(email, groupName string) error {
	return delete(p.bgr, UserGroup{
		Email:     email,
		GroupName: groupName,
	})
}
