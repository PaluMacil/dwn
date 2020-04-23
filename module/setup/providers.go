package setup

type Providers struct {
	Initialization InitializationProvider
}

type InitializationProvider interface {
	Get() (Initialization, error)
	Set(initialization Initialization) error
	EnsureDatabase() error
}
