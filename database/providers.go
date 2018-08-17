package database

import "github.com/PaluMacil/dwn/dwn"

type UserProvider interface {
	UserSearcher

	UsersFor(groupName string) ([]dwn.User, error)
	Get(email string) (dwn.User, error)
	Exists(email string) (bool, error)
	Set(user dwn.User) error
	Count() (int, error)
	All() ([]dwn.User, error)
	Delete(email string) error
	PurgeAll() error
}

type GroupProvider interface {
	GroupsFor(email string) ([]dwn.Group, error)
	Get(name string) (dwn.Group, error)
	Exists(name string) (bool, error)
	Set(group dwn.Group) error
	All() ([]dwn.Group, error)
	Delete(name string) error
}

type SessionProvider interface {
	Get(token string) (dwn.Session, error)
	Exists(token string) (bool, error)
	Set(session dwn.Session) error
	GenerateFor(email, ip string) dwn.Session
	All() ([]dwn.Session, error)
	Delete(token string) error
	PurgeAll() error
}

type UserGroupProvider interface {
	Get(email, groupName string) (dwn.UserGroup, error)
	Exists(email, groupName string) (bool, error)
	Set(userGroup dwn.UserGroup) error
	All() ([]dwn.UserGroup, error)
	Delete(email, groupName string) error
}

type SetupInfoProvider interface {
	Get() (dwn.SetupInfo, error)
	Completed() (bool, error)
	Set(setupInfo dwn.SetupInfo) error
	Delete() error
}
