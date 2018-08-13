package dwn

type UserSearcher interface {
	ReIndex() error
	CompletionSuggestions(query string) ([]User, error)
}
