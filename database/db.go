package database

import (
	"github.com/PaluMacil/dwn/database/store"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/module/dashboard"
	"github.com/PaluMacil/dwn/module/gamelibrary"
	"github.com/PaluMacil/dwn/module/logutil"
	"github.com/PaluMacil/dwn/module/setup"
	"github.com/PaluMacil/dwn/module/shopping"
)

func New(store Storer) *Database {
	return &Database{store: store}
}

type Database struct {
	store Storer
	core.Providers
	GameLibrary gamelibrary.Providers
	Shopping    shopping.Providers
	Log         logutil.Providers
	Setup       setup.Providers
	Dashboard   dashboard.Providers
}

func (db Database) Close() error {
	return db.store.Close()
}

func (db Database) NextID() (store.Identity, error) {
	return db.store.NextID()
}

func (db Database) IsKeyNotFoundErr(err error) bool {
	return db.store.IsKeyNotFoundErr(err)
}

func (db Database) KeyNotFoundErr() error {
	return db.store.KeyNotFoundErr()
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
	IsKeyNotFoundErr(err error) bool
	KeyNotFoundErr() error
	NextID() (store.Identity, error)
	Close() error
}
