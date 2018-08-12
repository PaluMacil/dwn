package db

import (
	"errors"
	"fmt"
	"time"
)

type UserGroup struct {
	Email       string    `json:"email"`
	GroupName   string    `json:"groupName"`
	CreatedDate time.Time `json:"createdDate"`
}

func (u UserGroup) Key() []byte {
	partOne := append(u.Prefix(), []byte(u.Email)...)
	return append(partOne, []byte(u.GroupName)...)
}

func (u UserGroup) Prefix() []byte {
	return []byte(userGroupPrefix)
}

type UserGroupProvider struct {
	Db *Db
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
	item, err := p.Db.get(&userGroup)
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
	if IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

func (p *UserGroupProvider) Set(userGroup UserGroup) error {
	if userGroup.CreatedDate.IsZero() {
		userGroup.CreatedDate = time.Now()
	}
	return p.Db.set(&userGroup)
}

func (p *UserGroupProvider) All() ([]UserGroup, error) {
	var items []DbItem
	err := p.Db.all(UserGroup{}.Prefix(), &items, true)
	userGroups := make([]UserGroup, len(items))
	for i, v := range items {
		userGroups[i] = v.(UserGroup)
	}

	return userGroups, err
}

func (p UserGroupProvider) Delete(email, groupName string) error {
	return p.Db.delete(UserGroup{
		Email:     email,
		GroupName: groupName,
	})
}
