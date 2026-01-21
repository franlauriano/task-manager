package configtest

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"taskmanager/internal/config"
)

// loadEnvFile loads environment variables from a .env file.
// It reads the file line by line and sets environment variables using os.Setenv.
// Lines starting with # are treated as comments and ignored.
// Empty lines are ignored.
func loadEnvFile(envPath string) error {
	file, err := os.Open(envPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		os.Setenv(key, value)
	}

	return scanner.Err()
}

// Load loads configuration from a TOML file after loading environment variables from a .env file.
// It first loads environment variables from envPath, then calls config.Load() with the configPath.
// The configPath and envPath can be relative to the project root or absolute paths.
// The v parameter is the interface that will be used to load the configuration.
func Load(configPath, envPath string, v interface{}) error {
	// Load environment variables from envPath
	absEnvPath, err := filepath.Abs(envPath)
	if err != nil {
		return err
	}

	if err := loadEnvFile(absEnvPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	return config.Load(absConfigPath, v)
}
