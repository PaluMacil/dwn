package badgerstore

import (
	"encoding/gob"
	"github.com/PaluMacil/dwn/module/core"

	"github.com/PaluMacil/dwn/module/shopping"
)

func init() {
	gob.Register(core.Session{})
	gob.Register(core.User{})
	gob.Register(core.Group{})
	gob.Register(core.UserGroup{})
	gob.Register(core.SetupInfo{})

	gob.Register(shopping.Item{})
}
