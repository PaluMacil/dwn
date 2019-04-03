package database

import "github.com/PaluMacil/dwn/core"

type CoreProviders struct {
	Sessions   SessionProvider
	Users      UserProvider
	Groups     GroupProvider
	UserGroups UserGroupProvider
	SetupInfo  SetupInfoProvider
	Util       UtilityProvider
}

type UserProvider interface {
	UserSearcher

	UsersFor(groupName string) ([]core.User, error)
	Get(email string) (core.User, error)
	Exists(email string) (bool, error)
	Set(user core.User) error
	Count() (int, error)
	All() ([]core.User, error)
	Delete(email string) error
	PurgeAll() error
}

type GroupProvider interface {
	GroupsFor(email string) ([]core.Group, error)
	Get(name string) (core.Group, error)
	Exists(name string) (bool, error)
	Set(group core.Group) error
	All() ([]core.Group, error)
	Delete(name string) error
}

type SessionProvider interface {
	Get(token string) (core.Session, error)
	Exists(token string) (bool, error)
	Set(session core.Session) error
	GenerateFor(email, ip string) core.Session
	All() ([]core.Session, error)
	Delete(token string) error
	PurgeAll() error
}

type UserGroupProvider interface {
	Get(email, groupName string) (core.UserGroup, error)
	Exists(email, groupName string) (bool, error)
	Set(userGroup core.UserGroup) error
	All() ([]core.UserGroup, error)
	Delete(email, groupName string) error
}

type SetupInfoProvider interface {
	Get() (core.SetupInfo, error)
	Completed() (bool, error)
	Set(setupInfo core.SetupInfo) error
	Delete() error
}
