package setup

import (
	"time"
)

const InitializationPrefix = "SETUPINITIALIZATION:"

type Initialization struct {
	DatabaseInitDate time.Time `json:"databaseInitDate"`
	WizardComplete   bool      `json:"wizardComplete"`
}

func (i Initialization) Key() []byte {
	return i.Prefix()
}

func (i Initialization) Prefix() []byte {
	return []byte(InitializationPrefix)
}
