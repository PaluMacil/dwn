package shopping

type Providers struct {
	Items ItemProvider
}

type ItemProvider interface {
	All() ([]Item, error)
	Delete(name string) error
	Get(name string) (Item, error)
	Set(item Item) error
}
