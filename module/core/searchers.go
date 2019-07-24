package core

type UserSearcher interface {
	Index(u User) error
	Deindex(u User) error
	Reindex() error
	CompletionSuggestions(query string) ([]User, error)
	FromEmail(email string) (User, error)
	WithEmail(email string) ([]User, error)
	EmailExists(email string) (bool, error)
	VerifiedEmailExists(email string) (bool, error)
}
