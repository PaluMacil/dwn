package core

import (
	"time"
)

const UserGroupPrefix = "USERGROUP:"

type UserGroup struct {
	Email       string    `json:"email"`
	GroupName   string    `json:"groupName"`
	CreatedDate time.Time `json:"createdDate"`
}

func (u UserGroup) Key() []byte {
	partOne := append(u.Prefix(), []byte(u.Email)...)
	return append(partOne, []byte(u.GroupName)...)
}

func (u UserGroup) Prefix() []byte {
	return []byte(UserGroupPrefix)
}
