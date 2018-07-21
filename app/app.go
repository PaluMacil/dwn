package app

import (
	"fmt"
	"os"

	"github.com/PaluMacil/dwn/db"
)

func New() (App, error) {
	dataDir := os.Getenv("DWN_DATA_DIR")
	database, err := db.New(dataDir)
	if err != nil {
		return App{}, fmt.Errorf(`initializing database in directory "%s": %s`, dataDir, err)
	}
	return App{
		Protocol:    os.Getenv("DWN_PROTOCOL"),
		Host:        os.Getenv("DWN_HOST"),
		Port:        os.Getenv("DWN_PORT"),
		BaseURL:     os.Getenv("DWN_BASE_URL"),
		UIProxyPort: os.Getenv("DWN_UI_PROXY_PORT"),
		Db:          database,
		Setup: Setup{
			InitialAdmin: os.Getenv("DWN_INITIAL_ADMIN"),
		},
	}, nil
}

type Setup struct {
	InitialAdmin string `json:"initialAdmin"`
}

type App struct {
	Protocol    string `json:"protocol"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	BaseURL     string `json:"baseURL"`
	UIProxyPort string `json:"uiProxyPort"`
	Db          *db.Db `json:"-"`
	Setup       Setup  `json:"setup"`
}

func (app App) HomePage() string {
	port := app.UIProxyPort
	if port == "" {
		port = app.Port
	}
	return fmt.Sprintf("%s://%s:%s%s",
		app.Protocol, app.Host, port, app.BaseURL)
}
