package badgerstore

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PaluMacil/dwn/database"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

type BadgerStore struct {
	bgr     *badger.DB
	dataDir string
}

func retry(dir string, originalOpts badger.Options) (*BadgerStore, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	bgr, err := badger.Open(retryOpts)
	return &BadgerStore{bgr: bgr, dataDir: dir}, err
}

func New(dir string, useMMAP bool) (*BadgerStore, error) {
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
			if bgr, err := retry(dir, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				return bgr, nil
			}
			log.Println("could not unlock database:", err)

		}
		return nil, err
	}
	bs := &BadgerStore{bgr: bgr, dataDir: dir}
	return bs, nil
}

func (bs BadgerStore) Close() error {
	return bs.bgr.Close()
}

func (bs *BadgerStore) Get(obj database.Item) (database.Item, error) {
	var rawBytes []byte
	err := bs.bgr.View(func(txn *badger.Txn) error {
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

func (bs *BadgerStore) Set(obj database.Item) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&obj)
	if err != nil {
		return err
	}
	return bs.bgr.Update(func(txn *badger.Txn) error {
		err := txn.Set(obj.Key(), buf.Bytes())
		if err != nil {
			return fmt.Errorf(`setting "%s" (%T): %s`, string(obj.Key()), obj, err)
		}
		return nil
	})
}

func (bs *BadgerStore) Delete(obj database.Item) error {
	return bs.bgr.Update(func(txn *badger.Txn) error {
		return txn.Delete(obj.Key())
	})
}

func (bs *BadgerStore) All(pfx []byte, out *[]database.Item, preload bool) error {
	err := bs.bgr.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(pfx); it.ValidForPrefix(pfx); it.Next() {
			item := it.Item()
			v, err := item.Value()
			if err != nil {
				return err
			}
			var buf bytes.Buffer
			var outItem database.Item
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

func (bs *BadgerStore) Count(pfx []byte) (int, error) {
	var items []database.Item
	err := bs.All(pfx, &items, false)
	return len(items), err
}
