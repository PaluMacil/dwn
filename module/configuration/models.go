package configuration

import (
	"fmt"
	"golang.org/x/oauth2"
	"os"
)

const (
	CredentialPrefix = "CREDENTIAL:"
)

type SetupConfiguration struct {
	InitialAdmin string `json:"initialAdmin"`
}

type Configuration struct {
	WebServer WebServerConfiguration     `json:"webServer"`
	Setup     SetupConfiguration         `json:"setup"`
	Database  DatabaseConfiguration      `json:"database"`
	FS        ForeignSystemConfiguration `json:"fs"`
	Prod      bool                       `json:"prod"`
}

type WebServerConfiguration struct {
	Protocol    string `json:"protocol"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	UIProxyPort string `json:"uiProxyPort"`
	ContentRoot string `json:"contentRoot"`
}

type DatabaseConfiguration struct {
	DataDir       string `json:"dataDir"`
	EncryptionKey []byte `json:"-"`
}

func (ws WebServerConfiguration) HomePage() string {
	port := ws.UIProxyPort
	if port == "" {
		port = ws.Port
	}
	return fmt.Sprintf("%s://%s:%s",
		ws.Protocol, ws.Host, port)
}

type AuthConfiguration struct {
	Google *oauth2.Config `json:"google,omitempty"`
}

type Credential struct {
	Key      string `json:"key"`
	Secret   string `json:"-"`
	Readonly bool   `json:"readonly"`
}

type SMTPConfiguration struct {
	SendGrid *Credential `json:"sendGrid,omitempty"`
}

type ForeignSystemConfiguration struct {
	Auth AuthConfiguration `json:"auth"`
	SMTP SMTPConfiguration `json:"smtp"`
}

type Mode struct {
	Prod bool
}

func (mode Mode) Coalesce(envKey, prodDefault, devDefault string) string {
	value := os.Getenv(envKey)
	if value == "" {
		if mode.Prod {
			return prodDefault
		} else {
			return devDefault
		}
	}
	return value
}
