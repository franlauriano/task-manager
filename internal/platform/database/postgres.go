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
	// databaseDBsWithoutTransactionsKey is the context key for the database without transactions.
	databaseDBsWithoutTransactionsKey contextKey = "database_without_transactions"
	// databaseDBsWithTransactionsKey is the context key for the database with transactions.
	databaseDBsWithTransactionsKey contextKey = "database_with_transactions"

	// DatabaseDefaultAlias is the default alias for the database.
	DatabaseDefaultAlias = "primary"
)

var (
	// dbs stores a map[string]*gorm.DB atomically and thread-safe.
	dbs atomic.Value
	// dbsMu protects write operations to dbs to prevent race conditions
	// when multiple database connections are set in parallel
	dbsMu sync.Mutex
	// ErrContextDatabase is returned when the context is without database.
	ErrContextDatabase = errors.New("context without database")
	// ErrDBNotFound is returned when the database connection is not found.
	ErrDBNotFound = errors.New("database connection not found")
	// ErrAliasConflict is returned when an alias is already used with a different connection type.
	ErrAliasConflict = errors.New("database alias already used with a different connection type")
)

// init ensures that the atomic.Value starts with an empty map to avoid nil errors.
func init() {
	dbs.Store(make(map[string]*gorm.DB))
}

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

// Open opens a new database connection and adds it to the global registry atomically,
// using DatabaseDefaultAlias.
func Open(config Configuration) error {
	return OpenAs(DatabaseDefaultAlias, config)
}

// OpenAs opens a new database connection under the given alias and adds it to the global registry.
func OpenAs(alias string, config Configuration) error {
	if err := validateOpenAsConfig(alias, config); err != nil {
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
		return fmt.Errorf("could not open database %s: %w", alias, err)
	}

	sqlDB, err := newConn.DB()
	if err == nil {
		applyPoolConfig(sqlDB, config)
	}

	dbsMu.Lock()
	defer dbsMu.Unlock()
	currentMap := dbs.Load().(map[string]*gorm.DB)

	newMap := make(map[string]*gorm.DB)
	for k, v := range currentMap {
		newMap[k] = v
	}

	newMap[alias] = newConn

	dbs.Store(newMap)

	slog.Info("Database connection opened", "alias", alias)
	return nil
}

// DByAlias returns a database connection by alias. It is safe for concurrent use.
func DByAlias(alias string) (*gorm.DB, error) {
	m := dbs.Load().(map[string]*gorm.DB)
	db, ok := m[alias]
	if !ok {
		return nil, ErrDBNotFound
	}
	return db, nil
}

// DBDefault returns the default database connection.
func DBDefault() (*gorm.DB, error) {
	return DByAlias(DatabaseDefaultAlias)
}

// InjectDBsIntoContext injects the database connections into the context.
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

// Commit commits all transactions in the context.
func Commit(ctx context.Context, aliases ...string) error {
	if len(aliases) == 0 {
		aliases = []string{DatabaseDefaultAlias}
	}
	for _, alias := range aliases {
		db, err := dbFromContext(ctx, databaseDBsWithTransactionsKey, alias)
		if err != nil {
			return err
		}
		if err := db.Commit().Error; err != nil {
			return err
		}
	}
	return nil
}

// Rollback rolls back all transactions in the context.
func Rollback(ctx context.Context, aliases ...string) error {
	if len(aliases) == 0 {
		aliases = []string{DatabaseDefaultAlias}
	}
	for _, alias := range aliases {
		db, err := dbFromContext(ctx, databaseDBsWithTransactionsKey, alias)
		if err != nil {
			return err
		}
		if err := db.Rollback().Error; err != nil {
			return err
		}
	}
	return nil
}

// CloseAll closes all registered database connections.
// This function is thread-safe and protects against race conditions
func CloseAll() error {
	dbsMu.Lock()
	defer dbsMu.Unlock()
	m := dbs.Load().(map[string]*gorm.DB)
	for alias, db := range m {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		sqlDB.Close()
		slog.Info("Database connection closed", "alias", alias)
	}

	dbs.Store(make(map[string]*gorm.DB))
	return nil
}

// DBDefaultFromContext returns the default database connection from the context.
func DBDefaultFromContext(ctx context.Context) (*gorm.DB, error) {
	return DBFromContext(ctx, DatabaseDefaultAlias)
}

// DBFromContext returns a database connection from the context by alias.
func DBFromContext(ctx context.Context, alias string) (*gorm.DB, error) {
	db, err := dbFromContext(ctx, databaseDBsWithoutTransactionsKey, alias)
	if err == nil {
		return db, nil
	}
	db, err = dbFromContext(ctx, databaseDBsWithTransactionsKey, alias)
	if err == nil {
		return db, nil
	}
	return nil, err
}

// SetDB stores a database connection by alias.
// This function is thread-safe and protects against race conditions
// when multiple database connections are set in parallel
func SetDB(alias string, db *gorm.DB) {
	dbsMu.Lock()
	defer dbsMu.Unlock()
	currentMap := dbs.Load().(map[string]*gorm.DB)
	newMap := make(map[string]*gorm.DB)
	for k, v := range currentMap {
		newMap[k] = v
	}
	newMap[alias] = db
	dbs.Store(newMap)
}

// dbFromContext extracts the map of database connections from the context.
func dbFromContext(ctx context.Context, key contextKey, alias string) (*gorm.DB, error) {
	if ctx == nil {
		return nil, ErrContextDatabase
	}

	dbsGorm, ok := ctx.Value(key).(map[string]*gorm.DB)
	if !ok || dbsGorm == nil {
		return nil, ErrContextDatabase
	}

	dbGorm, ok := dbsGorm[alias]
	if !ok {
		return nil, ErrDBNotFound
	}

	return dbGorm, nil
}

// validateOpenAsConfig validates alias and config; returns ErrInvalidConfig with details on failure.
func validateOpenAsConfig(alias string, config Configuration) error {
	var errs []string
	if strings.TrimSpace(alias) == "" {
		errs = append(errs, "alias is required")
	}
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
