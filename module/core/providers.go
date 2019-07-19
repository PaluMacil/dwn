package core

import (
	"github.com/PaluMacil/dwn/database/store"
)

type Providers struct {
	Sessions   SessionProvider
	Users      UserProvider
	Groups     GroupProvider
	UserGroups UserGroupProvider
}

type UserProvider interface {
	UserSearcher

	UsersFor(groupName string) ([]User, error)
	Get(userID store.Identity) (User, error)
	Exists(userID store.Identity) (bool, error)
	Set(user User) error
	Count() (int, error)
	All() ([]User, error)
	Delete(userID store.Identity) error
	PurgeAll() error
}

type GroupProvider interface {
	GroupsFor(userID store.Identity) ([]Group, error)
	Get(name string) (Group, error)
	Exists(name string) (bool, error)
	Set(group Group) error
	All() ([]Group, error)
	Delete(name string) error
}

type SessionProvider interface {
	Get(token string) (Session, error)
	Exists(token string) (bool, error)
	Set(session Session) error
	GenerateFor(userID store.Identity, ip string) Session
	All() ([]Session, error)
	Delete(token string) error
	PurgeAll() error
	UpdateHeartbeat(session *Session, ip string) error
}

type UserGroupProvider interface {
	Get(userID store.Identity, groupName string) (UserGroup, error)
	Exists(userID store.Identity, groupName string) (bool, error)
	Set(userGroup UserGroup) error
	All() ([]UserGroup, error)
	Delete(userID store.Identity, groupName string) error
}
