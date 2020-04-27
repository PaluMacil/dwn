package configuration

type Providers struct {
	ConfigProvider
}

type ConfigProvider interface {
	Get() Configuration
	Set(config Configuration)
}
