package application

import (
	"fmt"

	logutilrepo "github.com/PaluMacil/dwn/module/logutil/repo"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/configuration/env"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/stores/badgerstore"
	"github.com/PaluMacil/dwn/module/core/repo"
	"github.com/PaluMacil/dwn/module/core/search"
	"github.com/PaluMacil/dwn/webserver"

	setuprepo "github.com/PaluMacil/dwn/module/setup/repo"
	shoppingrepo "github.com/PaluMacil/dwn/module/shopping/repo"
)

func New() (*App, error) {
	config := env.Config()
	db, err := defaultDatabase(config.Database)
	if err != nil {
		return nil, fmt.Errorf(`initializing default database: %s`, err)
	}

	web := webserver.New(db, config)
	return &App{
		Config: config,
		Db:     db,
		Web:    web,
	}, nil
}

func defaultDatabase(config configuration.DatabaseConfiguration) (*database.Database, error) {
	// initialize store
	store, err := badgerstore.New(config.DataDir)
	if err != nil {
		return nil, fmt.Errorf(`initializing datastore in directory "%s": %s`, config.DataDir, err)
	}
	db := database.New(store)

	// initialize searchers from indexes
	userIndex, err := search.NewUserIndex(db, config.DataDir)
	if err != nil {
		return nil, fmt.Errorf(`initializing user index with dataDir "%s": %s`, config.DataDir, err)
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

	return db, nil
}

type App struct {
	Config configuration.Configuration `json:"config"`
	Db     *database.Database          `json:"-"`
	Web    *webserver.WebServer        `json:"-"`
}
