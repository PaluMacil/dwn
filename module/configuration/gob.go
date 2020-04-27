package configuration

import "encoding/gob"

func init() {
	gob.Register(Credential{})
}
