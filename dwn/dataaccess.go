package dwn

type DbItem interface {
	Key() []byte
	Prefix() []byte
}

type DataStorer interface {
	Get(obj DbItem) (DbItem, error)
	Set(obj DbItem) error
	Delete(obj DbItem) error
	All(pfx []byte, out *[]DbItem, preload bool) error
	Count(pfx []byte) (int, error)
	Close() error
}

type DbUtilityProvider interface {
	IsKeyNotFoundErr(err error) bool
}
