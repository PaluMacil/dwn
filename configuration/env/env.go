package env

import (
	"os"

	"github.com/PaluMacil/dwn/configuration"
)

func Config() configuration.Configuration {
	return configuration.Configuration{
		WebServer: configuration.WebServerConfiguration{
			Protocol:    os.Getenv("DWN_PROTOCOL"),
			Host:        os.Getenv("DWN_HOST"),
			Port:        os.Getenv("DWN_PORT"),
			BaseURL:     os.Getenv("DWN_BASE_URL"),
			UIProxyPort: os.Getenv("DWN_UI_PROXY_PORT"),
		},
		Setup: configuration.SetupConfiguration{
			InitialAdmin: os.Getenv("DWN_INITIAL_ADMIN"),
		},
		Database: configuration.DatabaseConfiguration{
			DataDir: os.Getenv("DWN_DATA_DIR"),
		},
	}
}
