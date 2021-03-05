package application

import (
	"github.com/contractGuru/pkg/config"
	"github.com/contractGuru/pkg/db"
)

// Application struct holds DB & configuration data for dependency injection
type Application struct {
	DB  *db.DB
	Cfg *config.Config
}

// GetApp captures env variables, establishes DB connection & returns reference to both.
func GetApp() (*Application, error) {
	cfg := config.GetConfig()

	db, err := db.GetDB(cfg.GetDBConnStr())

	if err != nil {
		return nil, err
	}

	return &Application{
		DB: db,
		Cfg: cfg,
	}, nil
}