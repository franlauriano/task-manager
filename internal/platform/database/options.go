package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Option is a function that returns a map of alias -> *gorm.DB.
type Option func(ctx context.Context) (context.Context, error)

// WithDBWithoutTransaction returns a map of alias -> *gorm.DB without transaction.
// None of the aliases may have been used in WithDBTransaction.
func WithDBWithoutTransaction(aliases ...string) Option {
	return func(ctx context.Context) (context.Context, error) {
		a := normalizeAliases(aliases)
		if err := validateAliasesNotInContext(ctx, databaseDBsWithTransactionsKey, a); err != nil {
			return nil, err
		}
		dbMap, err := buildDBMap(a, func(db *gorm.DB) (*gorm.DB, error) { return db, nil })
		if err != nil {
			return nil, err
		}
		return context.WithValue(ctx, databaseDBsWithoutTransactionsKey, dbMap), nil
	}
}

// WithDBTransaction returns a map of alias -> *gorm.DB with transaction.
// None of the aliases may have been used in WithDB.
func WithDBTransaction(aliases ...string) Option {
	return func(ctx context.Context) (context.Context, error) {
		a := normalizeAliases(aliases)
		if err := validateAliasesNotInContext(ctx, databaseDBsWithoutTransactionsKey, a); err != nil {
			return nil, err
		}
		dbMap, err := buildDBMap(a, func(db *gorm.DB) (*gorm.DB, error) {
			tx := db.Begin()
			return tx, tx.Error
		})
		if err != nil {
			return nil, err
		}
		return context.WithValue(ctx, databaseDBsWithTransactionsKey, dbMap), nil
	}
}

// aliasInContext reports whether the alias is already in the context under the given key.
func aliasInContext(ctx context.Context, key contextKey, alias string) bool {
	m, _ := ctx.Value(key).(map[string]*gorm.DB)
	_, ok := m[alias]
	return ok
}

// normalizeAliases returns aliases if non-empty, otherwise []string{DatabaseDefaultAlias}.
func normalizeAliases(aliases []string) []string {
	if len(aliases) == 0 {
		return []string{DatabaseDefaultAlias}
	}
	return aliases
}

// validateAliasesNotInContext returns an error if any alias is already in the context under conflictKey.
func validateAliasesNotInContext(ctx context.Context, conflictKey contextKey, aliases []string) error {
	for _, alias := range aliases {
		if aliasInContext(ctx, conflictKey, alias) {
			return fmt.Errorf("alias %q: %w", alias, ErrAliasConflict)
		}
	}
	return nil
}

// buildDBMap fetches *gorm.DB via DByAlias for each alias, applies wrap, and builds the map.
func buildDBMap(aliases []string, wrap func(*gorm.DB) (*gorm.DB, error)) (map[string]*gorm.DB, error) {
	dbMap := make(map[string]*gorm.DB, len(aliases))
	for _, alias := range aliases {
		db, err := DByAlias(alias)
		if err != nil {
			return nil, err
		}
		out, err := wrap(db)
		if err != nil {
			return nil, err
		}
		dbMap[alias] = out
	}
	return dbMap, nil
}
