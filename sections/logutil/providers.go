package logutil

type Providers struct {
	Config ConfigProvider
	Writer EntryProvider
}

type ConfigProvider interface {
	Get() (Config, error)
	Set(config Config) error
}

type EntryProvider interface {
	Write(p []byte) (n int, err error)
	Cached() ([]Entry, error)
	Log(level int, message string)
} //TODO: implement this with multiwriter and mutexes
