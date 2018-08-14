package app

import (
	"fmt"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/configuration/env"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/repo"
	"github.com/PaluMacil/dwn/database/search"
	"github.com/PaluMacil/dwn/database/stores/badgerstore"
)

func New() (*App, error) {
	config := env.Config()
	db, err := defaultDatabase(config.Database)
	if err != nil {
		return nil, fmt.Errorf(`initializing default database: %s`, err)
	}
	return &App{
		Config: config,
		Db:     db,
	}, nil
}

func defaultDatabase(config configuration.DatabaseConfiguration) (*database.Database, error) {
	// initialize store
	store, err := badgerstore.New(config.DataDir, true)
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
	db.SetupInfo = repo.NewSetupInfoRepo(store, db)
	db.Util = &badgerstore.Utility{}

	return db, nil
}

type App struct {
	Config configuration.Configuration `json:"config"`
	Db     *database.Database          `json:"-"`
}
