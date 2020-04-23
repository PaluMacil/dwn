package setup

import (
	"encoding/gob"
)

func init() {
	gob.Register(Initialization{})
}
