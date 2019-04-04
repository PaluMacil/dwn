package core

type Providers struct {
	Sessions   SessionProvider
	Users      UserProvider
	Groups     GroupProvider
	UserGroups UserGroupProvider
}

type UserProvider interface {
	UserSearcher

	UsersFor(groupName string) ([]User, error)
	Get(email string) (User, error)
	Exists(email string) (bool, error)
	Set(user User) error
	Count() (int, error)
	All() ([]User, error)
	Delete(email string) error
	PurgeAll() error
}

type GroupProvider interface {
	GroupsFor(email string) ([]Group, error)
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
	GenerateFor(email, ip string) Session
	All() ([]Session, error)
	Delete(token string) error
	PurgeAll() error
}

type UserGroupProvider interface {
	Get(email, groupName string) (UserGroup, error)
	Exists(email, groupName string) (bool, error)
	Set(userGroup UserGroup) error
	All() ([]UserGroup, error)
	Delete(email, groupName string) error
}
