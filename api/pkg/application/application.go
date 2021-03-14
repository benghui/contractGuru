package application

import (
	"github.com/contractGuru/pkg/config"
	"github.com/contractGuru/pkg/db"
	"github.com/contractGuru/pkg/logger"
	"github.com/contractGuru/pkg/router"
	"github.com/contractGuru/pkg/server"
)

// Application struct holds DB, server & configuration data for dependency injection
type Application struct {
	DB  *db.DB
	Cfg *config.Config
	Srv *server.Server
}

// GetApp captures env variables, initializes server & establishes DB connection then returns reference to all.
func GetApp() (*Application, error) {
	cfg := config.GetConfig()
	cert, err := cfg.GetCert()

	if err != nil {
		return nil, err
	}

	db, err := db.GetDB(cfg.GetDBConnStr())

	if err != nil {
		return nil, err
	}

	srv, err := server.GetServer(
		server.WithAddr(cfg.GetAPIPort()),
		server.WithRouter(router.GetRouter(db)),
		server.WithErrLogger(logger.Error),
		server.WithTLS(cert),
	)

	if err != nil {
		return nil, err
	}

	return &Application{
		DB:  db,
		Cfg: cfg,
		Srv: srv,
	}, nil
}
