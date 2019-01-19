// +build !windows

package badgerstore

import (
	"github.com/dgraph-io/badger"
)

func opts(dir string) badger.Options {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir
	return opts
}
