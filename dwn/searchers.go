package dwn

type UserSearcher interface {
	Index(u User) error
	ReIndex() error
	CompletionSuggestions(query string) ([]User, error)
}
