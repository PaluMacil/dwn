package configuration

type Providers struct {
	ConfigProvider
	Credential CredentialProvider
}

type ConfigProvider interface {
	Get() Configuration
	Set(config Configuration)
}

type CredentialProvider interface {
	Get(name string, fsType ForeignSystemType) (Credential, error)
	Set(cred Credential) error
	AllOf(fsType ForeignSystemType) ([]Credential, error)
	Delete(name string, fsType ForeignSystemType) error
}
