package auth

import (
	"github.com/PaluMacil/dwn/dwn"
)

func generateDisplayName(givenName string) (dwn.DisplayName, error) {
	//TODO: generate alternatives and check for repeats
	return dwn.DisplayName(givenName), nil
}
