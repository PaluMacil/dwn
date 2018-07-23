package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

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

type Db struct {
	bgr *badger.DB

	Sessions   SessionProvider
	Users      UserProvider
	Groups     GroupProvider
	UserGroups UserGroupProvider
	SetupInfo  SetupInfoProvider
}

func (db *Db) WithProviders() {
	db.Sessions = SessionProvider{db}
	db.Users = UserProvider{db}
	db.Groups = GroupProvider{db}
	db.UserGroups = UserGroupProvider{db}
	db.SetupInfo = SetupInfoProvider{db}
}

type DbItem interface {
	Key() []byte
	Prefix() []byte
}

func retry(dir string, originalOpts badger.Options) (*Db, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	bgr, err := badger.Open(retryOpts)
	return &Db{bgr: bgr}, err
}

func New(dir string, useMMAP bool) (*Db, error) {
	registerGobs()
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir
	if useMMAP {
		opts.ValueLogLoadingMode = options.MemoryMap
	} else {
		opts.ValueLogLoadingMode = options.FileIO
	}
	bgr, err := badger.Open(opts)
	if err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			log.Println("database locked, probably due to improper shutdown")
			if db, err := retry(dir, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				db.WithProviders()
				return db, nil
			}
			log.Println("could not unlock database:", err)

		}
		return nil, err
	}
	database := &Db{bgr: bgr}
	database.WithProviders()
	return database, nil
}

func (db Db) Close() error {
	return db.bgr.Close()
}

func (db *Db) get(obj DbItem) (DbItem, error) {
	var rawBytes []byte
	err := db.bgr.View(func(txn *badger.Txn) error {
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
		return obj, fmt.Errorf(`getting "%s" (%T): %s`, string(obj.Key()), obj, err)
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

func (db *Db) set(obj DbItem) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&obj)
	if err != nil {
		return err
	}
	return db.bgr.Update(func(txn *badger.Txn) error {
		err := txn.Set(obj.Key(), buf.Bytes())
		if err != nil {
			return fmt.Errorf(`setting "%s" (%T): %s`, string(obj.Key()), obj, err)
		}
		return nil
	})
}

func (db *Db) delete(obj DbItem) error {
	return db.bgr.Update(func(txn *badger.Txn) error {
		return txn.Delete(obj.Key())
	})
}

func (db *Db) all(pfx []byte, out *[]DbItem, preload bool) error {
	err := db.bgr.View(func(txn *badger.Txn) error {
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

func (db *Db) count(pfx []byte) (int, error) {
	var items []DbItem
	err := db.all(pfx, &items, false)
	return len(items), err
}

func IsKeyNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), badger.ErrKeyNotFound.Error())
}
