// +build !windows

package badgerstore

import (
	"github.com/dgraph-io/badger/v2"
)

func opts(dir string) badger.Options {
	return badger.DefaultOptions(dir)
}
