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
	InitialAdmin string
}

type App struct {
	Protocol    string
	Host        string
	Port        string
	BaseURL     string
	UIProxyPort string
	Db          *db.Db
	Setup       Setup
}

func (app App) HomePage() string {
	port := app.UIProxyPort
	if port == "" {
		port = app.Port
	}
	return fmt.Sprintf("%s://%s:%s%s",
		app.Protocol, app.Host, port, app.BaseURL)
}
