package configuration

import (
	"fmt"

	"golang.org/x/oauth2"
)

type SetupConfiguration struct {
	InitialAdmin string `json:"initialAdmin"`
}

type Configuration struct {
	WebServer WebServerConfiguration `json:"webServer"`
	Setup     SetupConfiguration     `json:"setup"`
	Database  DatabaseConfiguration  `json:"database"`
	Auth      AuthConfiguration      `json:"auth"`
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
	Google *oauth2.Config `json:"google"`
}
