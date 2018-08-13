package badgerstore

import (
	"encoding/gob"

	"github.com/PaluMacil/dwn/dwn"
)

func init() {
	gob.Register(dwn.Session{})
	gob.Register(dwn.User{})
	gob.Register(dwn.Group{})
	gob.Register(dwn.UserGroup{})
	gob.Register(dwn.SetupInfo{})
}
