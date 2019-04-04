package core

import (
	"encoding/gob"
)

func init() {
	gob.Register(Session{})
	gob.Register(User{})
	gob.Register(Group{})
	gob.Register(UserGroup{})
}
