package db

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger"
)

type Prefix string

const (
	sessionPrefix    = "SESSION:"
	userPrefix       = "USER:"
	groupPrefix      = "GROUP:"
	permissionPrefix = "PERMISSION:"
	userGroupPrefix  = "USERGROUP:"
	setupInfoPrefix  = "SETUPINFO:"
)

func registerGobs() {
	gob.Register(Session{})
	gob.Register(User{})
	gob.Register(Group{})
	gob.Register(UserGroup{})
	gob.Register(SetupInfo{})
}

type DbItem interface {
	Key() []byte
	Prefix() []byte
}

type Db struct {
	Close      func()
	Sessions   SessionProvider
	Users      UserProvider
	Groups     GroupProvider
	UserGroups UserGroupProvider
	SetupInfo  SetupInfoProvider
}

func New(dir string) *Db {
	registerGobs()
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir
	bgr, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	database := &Db{
		Close: func() {
			bgr.Close()
		},
	}
	database.Sessions = SessionProvider{bgr, database}
	database.Users = UserProvider{bgr, database}
	database.Groups = GroupProvider{bgr, database}
	database.UserGroups = UserGroupProvider{bgr, database}
	database.SetupInfo = SetupInfoProvider{bgr, database}
	return database
}

func get(bgr *badger.DB, obj DbItem) (DbItem, error) {
	var rawBytes []byte
	err := bgr.View(func(txn *badger.Txn) error {
		item, err := txn.Get(obj.Key())
		if err != nil {
			return err
		}
		value, err := item.Value()
		if err != nil {
			return err
		}
		rawBytes = make([]byte, len(value))
		copy(rawBytes, value)
		return nil
	})
	if err != nil {
		return obj, err
	}
	var buf bytes.Buffer
	_, err = buf.Write(rawBytes)
	if err != nil {
		return obj, err
	}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&obj)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

func set(bgr *badger.DB, obj DbItem) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&obj)
	if err != nil {
		return err
	}
	return bgr.Update(func(txn *badger.Txn) error {
		err := txn.Set(obj.Key(), buf.Bytes())
		if err != nil {
			return err
		}
		return nil
	})
}

func delete(bgr *badger.DB, obj DbItem) error {
	return bgr.Update(func(txn *badger.Txn) error {
		return txn.Delete(obj.Key())
	})
}

func all(bgr *badger.DB, pfx []byte, out *[]DbItem, preload bool) error {
	err := bgr.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(pfx); it.ValidForPrefix(pfx); it.Next() {
			item := it.Item()
			v, err := item.Value()
			if err != nil {
				return err
			}
			var buf bytes.Buffer
			var outItem DbItem
			_, err = buf.Write(v)
			if err != nil {
				return err
			}
			dec := gob.NewDecoder(&buf)
			err = dec.Decode(&outItem)
			if err != nil {
				return err
			}
			*out = append(*out, outItem)
		}
		return nil
	})
	return err
}

func count(bgr *badger.DB, pfx []byte) (int, error) {
	var items []DbItem
	err := all(bgr, pfx, &items, false)
	return len(items), err
}
