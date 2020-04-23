package badgerstore

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
)

func opts(dir string) badger.Options {
	return badger.DefaultOptions(dir).
		WithTableLoadingMode(options.FileIO).
		WithValueLogLoadingMode(options.FileIO)
}
