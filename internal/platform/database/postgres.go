package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// contextKey is the type used for context keys in this package to avoid collisions.
type contextKey string

const (
	// databaseWithoutTransactionKey is the context key for the database without transaction.
	databaseWithoutTransactionKey contextKey = "database_without_transaction"
	// databaseWithTransactionKey is the context key for the database with transaction.
	databaseWithTransactionKey contextKey = "database_with_transaction"
)

var (
	// ErrContextDatabase is returned when the context is without database.
	ErrContextDatabase = errors.New("context without database")
	// ErrDBNotFound is returned when the database connection is not set.
	ErrDBNotFound = errors.New("database connection not found")
	// ErrConflict is returned when the database is already used with a different connection type in context.
	ErrConflict = errors.New("database already used with a different connection type")
)

// Configuration holds PostgreSQL connection settings.
type Configuration struct {
	Host                   string `toml:"host"`
	User                   string `toml:"user"`
	Password               string `toml:"password"`
	Name                   string `toml:"name"`
	Port                   uint   `toml:"port"`
	Debug                  bool   `toml:"debug"`
	SSLMode                string `toml:"ssl_mode"`
	MaxOpenConns           int    `toml:"max_open_conns"`
	MaxIdleConns           int    `toml:"max_idle_conns"`
	ConnMaxLifetimeSeconds int    `toml:"conn_max_lifetime_seconds"`
}

// Open opens a new database connection and returns a Connector.
func Open(config Configuration) (Connector, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	sslMode := config.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s timezone=UTC",
		config.User, config.Password, config.Name, config.Host, config.Port, sslMode)

	debugLevel := logger.Silent
	if config.Debug {
		debugLevel = logger.Info
	}

	newConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(debugLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	sqlDB, err := newConn.DB()
	if err == nil {
		applyPoolConfig(sqlDB, config)
	}

	slog.Info("Database connection opened")
	return NewConnector(newConn), nil
}

// dbFromContext extracts the database connection from the context under the given key.
func dbFromContext(ctx context.Context, key contextKey) (*gorm.DB, error) {
	if ctx == nil {
		return nil, ErrContextDatabase
	}

	conn, ok := ctx.Value(key).(*gorm.DB)
	if !ok || conn == nil {
		return nil, ErrContextDatabase
	}

	return conn, nil
}

// validateConfig validates the configuration; returns error with details on failure.
func validateConfig(config Configuration) error {
	var errs []string
	if strings.TrimSpace(config.Host) == "" {
		errs = append(errs, "config.Host is required")
	}
	if strings.TrimSpace(config.User) == "" {
		errs = append(errs, "config.User is required")
	}
	if strings.TrimSpace(config.Name) == "" {
		errs = append(errs, "config.Name is required")
	}
	if config.Port == 0 {
		errs = append(errs, "config.Port must be greater than zero")
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}

// applyPoolConfig sets connection pool limits from config on sqlDB when values are greater than zero.
func applyPoolConfig(sqlDB *sql.DB, config Configuration) {
	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetimeSeconds > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetimeSeconds) * time.Second)
	}
}
