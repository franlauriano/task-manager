package dbtest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Container represents a PostgreSQL test container.
type Container struct {
	container testcontainers.Container
	db        *gorm.DB
	connStr   string
}

// Config holds configuration for the PostgreSQL test container.
type Config struct {
	image          string
	startupTimeout time.Duration
	migrationDir   string
}

// defaultConfig returns the default container configuration.
func defaultConfig() *Config {
	return &Config{
		image:          "postgres:18-alpine",
		startupTimeout: 30 * time.Second,
	}
}

// SetupDatabase creates and configures a PostgreSQL container for tests.
func SetupDatabase(t *testing.T, opts ...Option) (*Container, error) {
	if t != nil {
		t.Helper()
	}

	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	container, err := createContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres container: %w", err)
	}

	if t != nil {
		t.Cleanup(func() {
			_ = CleanDatabase(container.db)
		})
	}

	if cfg.migrationDir != "" {
		if t != nil {
			if err := container.runMigrations(t, cfg.migrationDir); err != nil {
				return nil, fmt.Errorf("failed to run migrations: %w", err)
			}
		} else {
			if err := container.runMigrationsWithoutT(cfg.migrationDir); err != nil {
				return nil, fmt.Errorf("failed to run migrations: %w", err)
			}
		}
	}

	return container, nil
}

// DB returns the GORM database connection.
func (c *Container) DB() *gorm.DB {
	return c.db
}

// TeardownDatabase terminates the container.
func (c *Container) TeardownDatabase() error {
	ctx := context.Background()
	return c.container.Terminate(ctx)
}

// runMigrations runs all .up.sql migration files from the given directory in order.
func (c *Container) runMigrations(t *testing.T, migrationDir string) error {
	t.Helper()
	return c.runMigrationsWithoutT(migrationDir)
}

// runMigrationsWithoutT runs migrations without requiring *testing.T.
func (c *Container) runMigrationsWithoutT(migrationDir string) error {
	absMigrationDir, err := filepath.Abs(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute migration path: %w", err)
	}

	entries, err := os.ReadDir(absMigrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}

	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		filePath := filepath.Join(absMigrationDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		if err := c.db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
	}

	return nil
}

// createContainer creates a PostgreSQL container using the provided configuration.
func createContainer(cfg *Config) (*Container, error) {
	ctx := context.Background()

	pgContainer, err := postgrescontainer.Run(ctx,
		cfg.image,
		testcontainers.WithCmd([]string{
			"postgres",
			"-c", "fsync=off",
			"-c", "synchronous_commit=off",
			"-c", "full_page_writes=off",
		}...),
		testcontainers.WithTmpfs(map[string]string{
			"/var/lib/postgresql": "rw",
		}),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
				wait.ForListeningPort("5432/tcp"),
			).WithDeadline(cfg.startupTimeout),
		),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, err
	}

	db, err := gorm.Open(postgresdriver.Open(connStr), &gorm.Config{})
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		return nil, err
	}

	return &Container{
		container: pgContainer,
		db:        db,
		connStr:   connStr,
	}, nil
}
