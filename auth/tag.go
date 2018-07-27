package auth

import (
	"github.com/PaluMacil/dwn/db"
)

func generateDisplayName(givenName string) (db.DisplayName, error) {
	//TODO: generate alternatives and check for repeats
	return db.DisplayName(givenName), nil
}
