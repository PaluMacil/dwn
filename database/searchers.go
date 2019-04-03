package database

import (
	"github.com/PaluMacil/dwn/module/core"
)

type UserSearcher interface {
	Index(u core.User) error
	Deindex(u core.User) error
	Reindex() error
	CompletionSuggestions(query string) ([]core.User, error)
}
