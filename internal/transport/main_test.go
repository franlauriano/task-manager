//go:build test

package transport

import (
	"log"
	"os"
	"testing"
	"time"

	"taskmanager/internal/paths"
	"taskmanager/internal/platform/cache"
	"taskmanager/internal/platform/database"
	"taskmanager/internal/platform/logger"
	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/redistest"
	taskRepo "taskmanager/internal/repository/task"
	"taskmanager/internal/testing/configtest"
	"taskmanager/internal/usecase/task"
	"taskmanager/internal/usecase/team"
)

var databaseTest *dbtest.Container
var redisTest *redistest.Container

func TestMain(m *testing.M) {
	os.Exit(func(m *testing.M) int {
		appConfig := struct {
			Database database.Configuration `toml:"database"`
			Logger   logger.Configuration   `toml:"logger"`
			Task     task.Configuration     `toml:"task"`
			Team     team.Configuration     `toml:"team"`
		}{}

		// Loading configs
		if err := configtest.Load(paths.TestConfigPath(), paths.TestEnvPath(), &appConfig); err != nil {
			log.Fatalf("Error on load config on struct. Err: %s", err)
		}

		// Set logger config
		appConfig.Logger.Initialize()

		// Load task config
		if err := task.LoadConfig(&appConfig.Task); err != nil {
			log.Fatalf("Error on load task config. Err: %s", err)
		}

		// Load team config
		if err := team.LoadConfig(&appConfig.Team); err != nil {
			log.Fatalf("Error on load team config. Err: %s", err)
		}

		// Setup database container for all tests in this package
		var err error
		if databaseTest, err = dbtest.SetupDatabase(nil,
			dbtest.WithMigrations(paths.MigrationDir()),
		); err != nil {
			log.Fatalf("Failed to setup database: %v", err)
		}
		defer func() {
			if err := databaseTest.TeardownDatabase(); err != nil {
				log.Printf("Failed to teardown database: %v", err)
			}
		}()

		// Setup Redis container for cache tests
		if redisTest, err = redistest.SetupRedis(nil); err != nil {
			log.Fatalf("Failed to setup redis: %v", err)
		}
		defer func() {
			if err := redisTest.TeardownRedis(); err != nil {
				log.Printf("Failed to teardown redis: %v", err)
			}
		}()

		// Wire cache-aside decorator (mirrors production wiring in cmd/main.go)
		cache.SetClient(redisTest.Client())
		taskRepo.SetPersist(taskRepo.NewCachedPersist(
			taskRepo.Persist(),
			redisTest.Client(),
			5*time.Minute,
		))

		return m.Run()
	}(m))
}
