package main

import (
	"fmt"
	"log"

	"taskmanager/internal/config"
	"taskmanager/internal/paths"
	"taskmanager/internal/platform/database"
	"taskmanager/internal/platform/logger"
	"taskmanager/internal/platform/server"
	"taskmanager/internal/transport"
	"taskmanager/internal/usecase/task"
	"taskmanager/internal/usecase/team"
)

func main() {
	appConfig := struct {
		Server   server.Configuration   `toml:"server"`
		Database database.Configuration `toml:"database"`
		Logger   logger.Configuration   `toml:"logger"`
		Task     task.Configuration     `toml:"task"`
		Team     team.Configuration     `toml:"team"`
	}{}

	// Load configuration from file with environment variable expansion
	if err := config.Load(paths.ConfigPath(), &appConfig); err != nil {
		log.Fatalf("Error on load config file: %s", err)
	}

	// Set logger config
	appConfig.Logger.Initialize()

	// Load task config
	if err := task.LoadConfig(&appConfig.Task); err != nil {
		log.Fatal("Error on load task config", "error", err)
	}

	// Load team config
	if err := team.LoadConfig(&appConfig.Team); err != nil {
		log.Fatal("Error on load team config", "error", err)
	}

	// Connect to database
	database.Open(appConfig.Database)
	defer func() {
		if err := database.CloseAll(); err != nil {
			log.Print("Error closing database:", err)
		}
	}()

	// Start http server
	address := fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port)
	server.ListenAndServe(address, transport.Routes())
}
