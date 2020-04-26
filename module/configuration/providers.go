package configuration

type Providers struct {
	Config ConfigProvider
}

type ConfigProvider interface {
	Get() (Configuration, error)
	Set(config Configuration)
}
