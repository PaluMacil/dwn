package badgerstore

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/PaluMacil/dwn/database"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
)

type BadgerStore struct {
	bgr      *badger.DB
	gcTicker *time.Ticker
}

func retry(originalOpts badger.Options) (*BadgerStore, error) {
	lockPath := filepath.Join(originalOpts.Dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	bgr, err := badger.Open(retryOpts)
	return &BadgerStore{bgr: bgr}, err
}

func New(dir string) (*BadgerStore, error) {
	bgr, err := badger.Open(opts(dir))
	if err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			log.Println("database locked, probably due to improper shutdown")
			if bgr, err := retry(opts(dir)); err == nil {
				log.Println("database unlocked, value log truncated")
				return bgr, nil
			}
			log.Println("could not unlock database:", err)

		}
		return nil, err
	}

	gcTicker := time.NewTicker(5 * time.Minute)
	bs := &BadgerStore{
		bgr:      bgr,
		gcTicker: gcTicker,
	}
	go func() {
		for range gcTicker.C {
			bs.runGC()
		}
	}()
	return bs, nil
}

func (bs *BadgerStore) runGC() {
	log.Println("Running GC...")
	var logFiles int
again:
	err := bs.bgr.RunValueLogGC(0.7)
	if err == nil {
		logFiles++
		goto again
	}
	log.Println(logFiles, "log files removed during GC.")
}

func (bs BadgerStore) Close() error {
	bs.gcTicker.Stop()
	bs.runGC()
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
