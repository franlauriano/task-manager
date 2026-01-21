package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	rootDir      string
	rootDirOnce  sync.Once
	rootDirError error
)

// RootDir returns the project root directory.
// It walks up the directory hierarchy until it finds go.mod.
// The result is cached after the first call.
func RootDir() (string, error) {
	rootDirOnce.Do(func() {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			rootDirError = os.ErrNotExist
			return
		}

		dir := filepath.Dir(filename)

		// Walk up until go.mod is found.
		for {
			goModPath := filepath.Join(dir, "go.mod")
			if _, err := os.Stat(goModPath); err == nil {
				rootDir = dir
				return
			}

			parent := filepath.Dir(dir)
			if parent == dir {
				rootDirError = os.ErrNotExist
				return
			}
			dir = parent
		}
	})

	return rootDir, rootDirError
}

// MustRootDir returns the root directory or panics if not found.
func MustRootDir() string {
	dir, err := RootDir()
	if err != nil {
		panic("go.mod not found: " + err.Error())
	}
	return dir
}

// MigrationDir returns the migrations directory.
func MigrationDir() string {
	return filepath.Join(MustRootDir(), "db", "migrate")
}

// SeedDir returns the seed directory.
func SeedDir() string {
	return filepath.Join(MustRootDir(), "db", "seed")
}

// FixtureDir returns the fixtures directory for tests.
func FixtureDir() string {
	return filepath.Join(MustRootDir(), "db", "fixtures")
}

// APITestDir returns the API tests directory.
func APITestDir() string {
	return filepath.Join(MustRootDir(), "api_test")
}

// ConfigDir returns the configuration directory.
func ConfigDir() string {
	return filepath.Join(MustRootDir(), "etc")
}

// ConfigPath returns the path to the application configuration file.
func ConfigPath() string {
	return filepath.Join(MustRootDir(), "etc", "config.toml")
}

// TestConfigPath returns the path to the test configuration file.
func TestConfigPath() string {
	return filepath.Join(MustRootDir(), "etc", "config_test.toml")
}

// TestEnvPath returns the path to the test environment file.
func TestEnvPath() string {
	return filepath.Join(MustRootDir(), "etc", ".env.test")
}
