package app

import (
	"fmt"
	"os"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/repo"
	"github.com/PaluMacil/dwn/database/search"
	"github.com/PaluMacil/dwn/database/stores/badgerstore"

	"strconv"
)

func New() (*App, error) {
	useMMAPEnv := os.Getenv("DWN_VALUE_LOG_USE_MMAP")
	useMMAP, err := strconv.ParseBool(useMMAPEnv)
	if err != nil {
		return nil, fmt.Errorf("parsing MMAP setting")
	}
	dataDir := os.Getenv("DWN_DATA_DIR")
	db, err := defaultDatabase(dataDir)
	if err != nil {
		return nil, fmt.Errorf(`initializing database in directory "%s": %s`, dataDir, err)
	}
	return &App{
		Protocol:        os.Getenv("DWN_PROTOCOL"),
		Host:            os.Getenv("DWN_HOST"),
		Port:            os.Getenv("DWN_PORT"),
		BaseURL:         os.Getenv("DWN_BASE_URL"),
		UIProxyPort:     os.Getenv("DWN_UI_PROXY_PORT"),
		ValueLogUseMMAP: useMMAP,
		Db:              db,
		Setup: Setup{
			InitialAdmin: os.Getenv("DWN_INITIAL_ADMIN"),
		},
	}, nil
}

func defaultDatabase(dataDir string) (database.Database, error) {
	// initialize store
	store, err := badgerstore.New(dataDir, true)
	if err != nil {
		return database.Database{}, fmt.Errorf(`initializing datastore in directory "%s": %s`, dataDir, err)
	}
	db := database.New(store)

	// initialize searchers from indexes
	userIndex, err := search.NewUserIndex(db, dataDir)
	if err != nil {
		return database.Database{}, fmt.Errorf(`initializing user index with dataDir "%s": %s`, dataDir, err)
	}

	// initialize providers from repo package
	db.Sessions = repo.NewSessionRepo(store, db)
	db.Users = repo.NewUserRepo(store, db, userIndex)
	db.Groups = repo.NewGroupRepo(store, db)
	db.UserGroups = repo.NewUserGroupRepo(store, db)
	db.SetupInfo = repo.NewSetupInfoRepo(store, db)
	db.Util = &badgerstore.Utility{}

	return db, nil
}

type Setup struct {
	InitialAdmin string `json:"initialAdmin"`
}

type App struct {
	Protocol        string            `json:"protocol"`
	Host            string            `json:"host"`
	Port            string            `json:"port"`
	BaseURL         string            `json:"baseURL"`
	UIProxyPort     string            `json:"uiProxyPort"`
	ValueLogUseMMAP bool              `json:"valueLogUseMMAP"`
	Db              database.Database `json:"-"`
	Setup           Setup             `json:"setup"`
}

func (app App) HomePage() string {
	port := app.UIProxyPort
	if port == "" {
		port = app.Port
	}
	return fmt.Sprintf("%s://%s:%s%s",
		app.Protocol, app.Host, port, app.BaseURL)
}
