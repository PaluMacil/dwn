package badgerstore

import (
	"github.com/PaluMacil/dwn/module/configuration"
	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
)

func opts(config configuration.DatabaseConfiguration) badger.Options {
	return badger.DefaultOptions(config.DataDir).
		WithTableLoadingMode(options.FileIO).
		WithValueLogLoadingMode(options.FileIO)
}
