// +build !windows

package badgerstore

import (
	"github.com/PaluMacil/dwn/configuration"
	"github.com/dgraph-io/badger/v2"
)

func opts(config configuration.DatabaseConfiguration) badger.Options {
	return badger.DefaultOptions(config.DataDir).
		WithEncryptionKey(config.EncryptionKey)
}
