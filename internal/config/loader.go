package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Load loads configuration from a TOML file with environment variable expansion.
// It supports the format ${VAR_NAME} and ${VAR:-default} for environment variables.
func Load(path string, v interface{}) error {
	rawBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	expandedContent := os.Expand(string(rawBytes), func(placeholder string) string {
		key := placeholder
		defaultValue := ""

		if strings.Contains(placeholder, ":-") {
			parts := strings.SplitN(placeholder, ":-", 2)
			key = parts[0]
			defaultValue = parts[1]
		}

		if val, ok := os.LookupEnv(key); ok {
			return val
		}

		return defaultValue
	})

	if err := toml.Unmarshal([]byte(expandedContent), v); err != nil {
		var decodeErr *toml.DecodeError
		if errors.As(err, &decodeErr) {
			row, col := decodeErr.Position()
			return fmt.Errorf("error in file '%s' at line %d, column %d: %w", path, row, col, err)
		}
		return fmt.Errorf("error decoding TOML in file '%s': %w", path, err)
	}

	return nil
}
