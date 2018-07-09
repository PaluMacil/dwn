package db

import (
	"time"

	"github.com/dgraph-io/badger"
)

type UserGroup struct {
	Email       string
	GroupName   string
	CreatedDate time.Time
}

func (u UserGroup) Key() []byte {
	partOne := append(u.Prefix(), []byte(u.Email)...)
	return append(partOne, []byte(u.GroupName)...)
}

func (u UserGroup) Prefix() []byte {
	return []byte(userGroupPrefix)
}

type UserGroupProvider struct {
	bgr *badger.DB
}
