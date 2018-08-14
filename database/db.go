package database

func New(store Storer) *Database {
	return &Database{store: store}
}

type Database struct {
	store      Storer
	Sessions   SessionProvider
	Users      UserProvider
	Groups     GroupProvider
	UserGroups UserGroupProvider
	SetupInfo  SetupInfoProvider
	Util       UtilityProvider
}

func (db Database) Close() error {
	return db.store.Close()
}

type Item interface {
	Key() []byte
	Prefix() []byte
}

type Storer interface {
	Get(obj Item) (Item, error)
	Set(obj Item) error
	Delete(obj Item) error
	All(pfx []byte, out *[]Item, preload bool) error
	Count(pfx []byte) (int, error)
	Close() error
}

type UtilityProvider interface {
	IsKeyNotFoundErr(err error) bool
}
