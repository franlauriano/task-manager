package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client   atomic.Pointer[redis.Client]
	clientMu sync.Mutex

	ErrClientNotFound = errors.New("cache client not found")
)

// Configuration holds Redis connection settings.
type Configuration struct {
	Host               string `toml:"host"`
	Port               uint   `toml:"port"`
	Password           string `toml:"password"`
	DB                 int    `toml:"db"`
	DefaultTTLSeconds  int    `toml:"default_ttl_seconds"`
}

// DefaultTTL returns the configured TTL as a time.Duration.
func (c Configuration) DefaultTTL() time.Duration {
	if c.DefaultTTLSeconds <= 0 {
		return 5 * time.Minute
	}
	return time.Duration(c.DefaultTTLSeconds) * time.Second
}

// Open creates a new Redis client connection and stores it globally.
func Open(config Configuration) error {
	if err := validateConfig(config); err != nil {
		return err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("could not connect to redis: %w", err)
	}

	clientMu.Lock()
	defer clientMu.Unlock()
	client.Store(rdb)

	slog.Info("Cache connection opened")
	return nil
}

// Client returns the current Redis client. It is safe for concurrent use.
func Client() (*redis.Client, error) {
	c := client.Load()
	if c == nil {
		return nil, ErrClientNotFound
	}
	return c, nil
}

// SetClient stores a Redis client.
// This function is thread-safe and protects against race conditions.
func SetClient(c *redis.Client) {
	clientMu.Lock()
	defer clientMu.Unlock()
	client.Store(c)
}

// Close closes the Redis client connection.
// This function is thread-safe and protects against race conditions.
func Close() error {
	clientMu.Lock()
	defer clientMu.Unlock()
	c := client.Load()
	if c == nil {
		return nil
	}
	err := c.Close()
	client.Store(nil)
	slog.Info("Cache connection closed")
	return err
}

// validateConfig validates the configuration; returns error with details on failure.
func validateConfig(config Configuration) error {
	var errs []string
	if strings.TrimSpace(config.Host) == "" {
		errs = append(errs, "config.Host is required")
	}
	if config.Port == 0 {
		errs = append(errs, "config.Port must be greater than zero")
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}
