package core

type UserSearcher interface {
	Index(u User) error
	Deindex(u User) error
	Reindex() error
	CompletionSuggestions(query string) ([]User, error)
}
