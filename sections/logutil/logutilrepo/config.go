package logutilrepo

import (
	"fmt"
	"time"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/sections/logutil"
)

type LogConfigRepo struct {
	store database.Storer
	db    *database.Database
}

func NewLogConfigRepo(store database.Storer, db *database.Database) *LogConfigRepo {
	return &LogConfigRepo{store, db}
}

func (r *LogConfigRepo) Get() (logutil.Config, error) {
	// if the LogConfigRepo is not yet set up or no config is
	// saved to the database, return the default configuration.
	if r == nil {
		return defaultConfig, nil
	}
	var config = logutil.Config{}
	item, err := r.store.Get(&config)
	if r.db.Util.IsKeyNotFoundErr(err) {
		return defaultConfig, nil
	} else if err != nil {
		return config, err
	}
	config, ok := item.(logutil.Config)
	if !ok {
		return config, fmt.Errorf("got data of type %T but wanted logutil.Config", config)
	}
	return config, err
}

func (r LogConfigRepo) Set(config logutil.Config) error {
	return r.store.Set(&config)
}

var defaultConfig logutil.Config = logutil.Config{
	ConsoleLevel:  logutil.LevelDebug,
	UseColorCodes: false,
	QueueLevel:    logutil.LevelDebug,
	Modified:      time.Now(),
}
