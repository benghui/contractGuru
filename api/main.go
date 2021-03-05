package main

import (
	"github.com/contractGuru/pkg/application"
	"github.com/contractGuru/pkg/exithandler"
	"github.com/contractGuru/pkg/logger"
	"github.com/contractGuru/pkg/router"
	"github.com/contractGuru/pkg/server"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Info.Println("Failed to load env vars")
	}
}

func main() {
	app, err := application.GetApp()

	if err != nil {
		logger.Error.Fatal(err.Error())
	}

	srv := server.
		GetServer().
		WithAddr(app.Cfg.GetAPIPort()).
		WithRouter(router.GetRouter(app)).
		WithErrLogger(logger.Error)

	go func() {
		logger.Info.Printf("Starting server. Listening at port %s\n", app.Cfg.GetAPIPort())

		if err := srv.StartServer(); err != nil {
			logger.Error.Fatal(err.Error())
		}
	}()

	exithandler.Exit(func() {
		if err := srv.CloseServer(); err != nil {
			logger.Error.Println(err.Error())
		}

		if err := app.DB.CloseDB(); err != nil {
			logger.Error.Println(err.Error())
		}
	})
}
