package shopping

import (
	"encoding/gob"
)

func init() {
	gob.Register(Item{})
}
