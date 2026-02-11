package redistest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Container represents a Redis test container.
type Container struct {
	container testcontainers.Container
	client    *redis.Client
}

// Config holds configuration for the Redis test container.
type Config struct {
	image          string
	startupTimeout time.Duration
}

// defaultConfig returns the default container configuration.
func defaultConfig() *Config {
	return &Config{
		image:          "redis:8-alpine",
		startupTimeout: 30 * time.Second,
	}
}

// SetupRedis creates and configures a Redis container for tests.
func SetupRedis(t *testing.T, opts ...Option) (*Container, error) {
	if t != nil {
		t.Helper()
	}

	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	container, err := createContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis container: %w", err)
	}

	if t != nil {
		t.Cleanup(func() {
			_ = container.client.FlushAll(context.Background()).Err()
		})
	}

	return container, nil
}

// Client returns the Redis client connected to the container.
func (c *Container) Client() *redis.Client {
	return c.client
}

// TeardownRedis terminates the container.
func (c *Container) TeardownRedis() error {
	if c.client != nil {
		_ = c.client.Close()
	}
	return c.container.Terminate(context.Background())
}

// FlushAll removes all keys from the Redis instance.
func (c *Container) FlushAll() error {
	return c.client.FlushAll(context.Background()).Err()
}

// createContainer creates a Redis container using the provided configuration.
func createContainer(cfg *Config) (*Container, error) {
	ctx := context.Background()

	redisContainer, err := rediscontainer.Run(ctx,
		cfg.image,
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(cfg.startupTimeout),
		),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		_ = redisContainer.Terminate(ctx)
		return nil, err
	}

	opts, err := redis.ParseURL(connStr)
	if err != nil {
		_ = redisContainer.Terminate(ctx)
		return nil, err
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		_ = redisContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &Container{
		container: redisContainer,
		client:    client,
	}, nil
}
