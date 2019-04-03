package database

import (
	"github.com/PaluMacil/dwn/module/gamelibrary"
	"github.com/PaluMacil/dwn/module/logutil"
	"github.com/PaluMacil/dwn/module/shopping"
)

func New(store Storer) *Database {
	return &Database{store: store}
}

type Database struct {
	store Storer
	CoreProviders
	GameLibrary gamelibrary.Providers
	Shopping    shopping.Providers
	Log         logutil.Providers
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
