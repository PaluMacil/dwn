package shopping

type Providers struct {
	Items ItemProvider
}

type ItemProvider interface {
	All() ([]Item, error)
	Delete(name string) error
}
