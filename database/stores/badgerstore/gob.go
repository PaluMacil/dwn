package badgerstore

import (
	"encoding/gob"

	"github.com/PaluMacil/dwn/module/core"

	"github.com/PaluMacil/dwn/module/shopping"
)

// TODO: rehome this in the inits of each module
func init() {
	gob.Register(core.Session{})
	gob.Register(core.User{})
	gob.Register(core.Group{})
	gob.Register(core.UserGroup{})
	gob.Register(core.SetupInfo{})

	gob.Register(shopping.Item{})
	//TODO: don't forget to register logging-related structs
}
