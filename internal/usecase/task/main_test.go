//go:build test

package task

import (
	"log"
	"os"
	"testing"

	"taskmanager/internal/paths"
	"taskmanager/internal/platform/database"
	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/testing/configtest"
)

var databaseTest *dbtest.Container

func TestMain(m *testing.M) {
	os.Exit(func(m *testing.M) int {
		appConfig := struct {
			Database database.Configuration `toml:"database"`
		}{}

		// Loading configs
		if err := configtest.Load(paths.TestConfigPath(), paths.TestEnvPath(), &appConfig); err != nil {
			log.Fatalf("Error on load config on struct. Err: %s", err)
		}

		return m.Run()
	}(m))
}
