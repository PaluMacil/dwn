package badgerstore

import (
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

func opts(dir string) badger.Options {
	opts := badger.DefaultOptions
	opts.TableLoadingMode = options.FileIO
	opts.ValueLogLoadingMode = options.FileIO
	opts.Dir = dir
	opts.ValueDir = dir
	return opts
}
