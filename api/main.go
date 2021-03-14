package main

import (
	"github.com/contractGuru/pkg/application"
	"github.com/contractGuru/pkg/exithandler"
	"github.com/contractGuru/pkg/logger"
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

	go func() {
		logger.Info.Printf("Starting server. Listening at port %s\n", app.Cfg.GetAPIPort())

		if err := app.Srv.StartServerTLS(); err != nil {
			logger.Error.Fatal(err.Error())
		}
	}()

	exithandler.Exit(func() {
		if err := app.Srv.CloseServer(); err != nil {
			logger.Error.Println(err.Error())
		}

		if err := app.DB.CloseDB(); err != nil {
			logger.Error.Println(err.Error())
		}
	})
}
