package badgerstore

import (
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

func opts(dir string) badger.Options {
	return badger.DefaultOptions(dir).
		WithTableLoadingMode(options.FileIO).
		WithValueLogLoadingMode(options.FileIO)
}
