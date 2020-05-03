package setup

import (
	"github.com/PaluMacil/dwn/database/store"
	"time"
)

const InitializationPrefix = "SETUPINITIALIZATION:"

type Initialization struct {
	DatabaseInitDate time.Time      `json:"databaseInitDate"`
	WizardComplete   bool           `json:"wizardComplete"`
	SetupUserID      store.Identity `json:"setupUserID"`
	EnvUserID        store.Identity `json:"envUserID"`
}

func (i Initialization) Key() []byte {
	return i.Prefix()
}

func (i Initialization) Prefix() []byte {
	return []byte(InitializationPrefix)
}
