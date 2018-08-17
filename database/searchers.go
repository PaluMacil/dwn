package database

import "github.com/PaluMacil/dwn/dwn"

type UserSearcher interface {
	Index(u dwn.User) error
	Deindex(u dwn.User) error
	Reindex() error
	CompletionSuggestions(query string) ([]dwn.User, error)
}
