package auth

import (
	"github.com/PaluMacil/dwn/module/core"
)

func generateDisplayName(givenName string) (core.DisplayName, error) {
	//TODO: generate alternatives and check for repeats
	return core.DisplayName(givenName), nil
}
