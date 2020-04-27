package application

import (
	"fmt"
	configrepo "github.com/PaluMacil/dwn/module/configuration/repo"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/store/badgerstore"
	"github.com/PaluMacil/dwn/module/configuration"
	"github.com/PaluMacil/dwn/module/core/search"
	"github.com/PaluMacil/dwn/webserver"

	"github.com/PaluMacil/dwn/module/core/repo"
	dashboardrepo "github.com/PaluMacil/dwn/module/dashboard/repo"
	logutilrepo "github.com/PaluMacil/dwn/module/logutil/repo"
	setuprepo "github.com/PaluMacil/dwn/module/setup/repo"
	shoppingrepo "github.com/PaluMacil/dwn/module/shopping/repo"
)

func New(prod bool) (*App, error) {
	configProvider, err := configrepo.NewConfigurationRepo(prod)
	if err != nil {
		return nil, fmt.Errorf("creating new config repo: %w", err)
	}
	config := configProvider.Get()
	db, store, err := defaultDatabase(config.Database)
	if err != nil {
		return nil, fmt.Errorf(`initializing default database: %w`, err)
	}

	configProvider.ConnectDatabase(store, db)
	db.Config.ConfigProvider = configProvider
	web := webserver.New(db)
	return &App{
		Db:  db,
		Web: web,
	}, nil
}

func defaultDatabase(config configuration.DatabaseConfiguration) (*database.Database, database.Storer, error) {

	// initialize store
	store, err := badgerstore.New(config)
	if err != nil {
		return nil, nil, fmt.Errorf(`initializing datastore in directory "%s": %w`, config.DataDir, err)
	}
	db := database.New(store)

	// initialize searchers from indexes
	userIndex, err := search.NewUserIndex(db, config.DataDir)
	if err != nil {
		return nil, nil, fmt.Errorf(`initializing user index with dataDir "%s": %w`, config.DataDir, err)
	}

	// initialize providers from repo package
	db.Sessions = repo.NewSessionRepo(store, db)
	db.Users = repo.NewUserRepo(store, db, userIndex)
	db.Groups = repo.NewGroupRepo(store, db)
	db.UserGroups = repo.NewUserGroupRepo(store, db)

	db.Shopping.Items = shoppingrepo.NewItemRepo(store, db)
	db.Log.Config = logutilrepo.NewLogConfigRepo(store, db)
	//db.Log.Writer = repo.NewEntryRepo(store, db)
	db.Setup.Initialization = setuprepo.NewInitializationRepo(store, db)

	// Initialize dashboard providers from repos
	db.Dashboard.Board = dashboardrepo.NewDashboardRepo(store, db)
	db.Dashboard.Projects = dashboardrepo.NewProjectRepo(store, db)

	return db, store, nil
}

type App struct {
	Config configuration.Configuration `json:"config"`
	Db     *database.Database          `json:"-"`
	Web    *webserver.WebServer        `json:"-"`
}
