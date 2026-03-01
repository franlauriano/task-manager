package database

import (
	"context"

	"gorm.io/gorm"
)

// Option is a function that injects a database connection into the context.
// It receives the base connection from the provider.
type Option func(ctx context.Context, conn *gorm.DB) (context.Context, error)

// WithDBWithoutTransaction injects the database connection into the context without a transaction.
func WithDBWithoutTransaction() Option {
	return func(ctx context.Context, conn *gorm.DB) (context.Context, error) {
		if hasDBInContext(ctx, databaseWithTransactionKey) {
			return nil, ErrConflict
		}
		return context.WithValue(ctx, databaseWithoutTransactionKey, conn), nil
	}
}

// WithDBTransaction injects the database connection into the context with a transaction.
func WithDBTransaction() Option {
	return func(ctx context.Context, conn *gorm.DB) (context.Context, error) {
		if hasDBInContext(ctx, databaseWithoutTransactionKey) {
			return nil, ErrConflict
		}
		tx := conn.Begin()
		if tx.Error != nil {
			return nil, tx.Error
		}
		return context.WithValue(ctx, databaseWithTransactionKey, tx), nil
	}
}

// hasDBInContext reports whether a database connection exists in the context under the given key.
func hasDBInContext(ctx context.Context, key contextKey) bool {
	conn, _ := ctx.Value(key).(*gorm.DB)
	return conn != nil
}
