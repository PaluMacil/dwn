package badgerstore

import (
	"encoding/gob"

	"github.com/PaluMacil/dwn/dwn"
	"github.com/PaluMacil/dwn/sections/shopping"
)

func init() {
	gob.Register(dwn.Session{})
	gob.Register(dwn.User{})
	gob.Register(dwn.Group{})
	gob.Register(dwn.UserGroup{})
	gob.Register(dwn.SetupInfo{})

	gob.Register(shopping.Item{})
}
