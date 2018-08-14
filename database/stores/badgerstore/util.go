package badgerstore

import (
	"strings"

	"github.com/dgraph-io/badger"
)

type Utility struct{}

func (u Utility) IsKeyNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), badger.ErrKeyNotFound.Error())
}
