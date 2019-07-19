package core

import (
	"time"

	"github.com/PaluMacil/dwn/database/store"
)

const UserGroupPrefix = "USERGROUP:"

type UserGroup struct {
	UserID      store.Identity `json:"userID"`
	GroupName   string         `json:"groupName"`
	CreatedDate time.Time      `json:"createdDate"`
}

func (u UserGroup) Key() []byte {
	partOne := append(u.Prefix(), u.UserID.Bytes()...)
	return append(partOne, []byte(u.GroupName)...)
}

func (u UserGroup) Prefix() []byte {
	return []byte(UserGroupPrefix)
}
