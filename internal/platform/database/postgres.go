package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
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
	// db stores the single *gorm.DB connection atomically and thread-safe.
	db atomic.Pointer[gorm.DB]
	// dbMu protects write operations to db to prevent race conditions.
	dbMu sync.Mutex
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

// Open opens a new database connection and stores it globally.
func Open(config Configuration) error {
	if err := validateConfig(config); err != nil {
		return err
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
		return fmt.Errorf("could not open database: %w", err)
	}

	sqlDB, err := newConn.DB()
	if err == nil {
		applyPoolConfig(sqlDB, config)
	}

	dbMu.Lock()
	defer dbMu.Unlock()
	db.Store(newConn)

	slog.Info("Database connection opened")
	return nil
}

// DB returns the current database connection. It is safe for concurrent use.
func DB() (*gorm.DB, error) {
	conn := db.Load()
	if conn == nil {
		return nil, ErrDBNotFound
	}
	return conn, nil
}

// SetDB stores a database connection.
// This function is thread-safe and protects against race conditions.
func SetDB(conn *gorm.DB) {
	dbMu.Lock()
	defer dbMu.Unlock()
	db.Store(conn)
}

// InjectDBsIntoContext injects the database connection into the context.
func InjectDBsIntoContext(ctx context.Context, options ...Option) (context.Context, error) {
	err := error(nil)
	for _, opt := range options {
		ctx, err = opt(ctx)
		if err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

// Commit commits the transaction in the context.
func Commit(ctx context.Context) error {
	conn, err := dbFromContext(ctx, databaseWithTransactionKey)
	if err != nil {
		return err
	}
	return conn.Commit().Error
}

// Rollback rolls back the transaction in the context.
func Rollback(ctx context.Context) error {
	conn, err := dbFromContext(ctx, databaseWithTransactionKey)
	if err != nil {
		return err
	}
	return conn.Rollback().Error
}

// Close closes the database connection.
// This function is thread-safe and protects against race conditions.
func Close() error {
	dbMu.Lock()
	defer dbMu.Unlock()
	conn := db.Load()
	if conn == nil {
		return nil
	}
	sqlDB, err := conn.DB()
	if err != nil {
		return err
	}
	sqlDB.Close()
	db.Store(nil)
	slog.Info("Database connection closed")
	return nil
}

// DBFromContext returns the database connection from the context.
// It first checks for a non-transactional connection, then a transactional one.
func DBFromContext(ctx context.Context) (*gorm.DB, error) {
	conn, err := dbFromContext(ctx, databaseWithoutTransactionKey)
	if err == nil {
		return conn, nil
	}
	conn, err = dbFromContext(ctx, databaseWithTransactionKey)
	if err == nil {
		return conn, nil
	}
	return nil, err
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
