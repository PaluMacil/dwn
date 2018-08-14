package configuration

import "fmt"

type SetupConfiguration struct {
	InitialAdmin string `json:"initialAdmin"`
}

type Configuration struct {
	WebServer WebServerConfiguration `json:"webServer"`
	Setup     SetupConfiguration     `json:"setup"`
	Database  DatabaseConfiguration  `json:"database"`
}

type WebServerConfiguration struct {
	Protocol    string `json:"protocol"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	BaseURL     string `json:"baseURL"`
	UIProxyPort string `json:"uiProxyPort"`
}

type DatabaseConfiguration struct {
	DataDir string `json:"dataDir"`
}

func (ws WebServerConfiguration) HomePage() string {
	port := ws.UIProxyPort
	if port == "" {
		port = ws.Port
	}
	return fmt.Sprintf("%s://%s:%s%s",
		ws.Protocol, ws.Host, port, ws.BaseURL)
}
