package database

import (
	"context"

	"gorm.io/gorm"
)

// Option is a function that injects a database connection into the context.
type Option func(ctx context.Context) (context.Context, error)

// WithDBWithoutTransaction injects the database connection into the context without a transaction.
func WithDBWithoutTransaction() Option {
	return func(ctx context.Context) (context.Context, error) {
		if hasDBInContext(ctx, databaseWithTransactionKey) {
			return nil, ErrConflict
		}
		conn, err := DB()
		if err != nil {
			return nil, err
		}
		return context.WithValue(ctx, databaseWithoutTransactionKey, conn), nil
	}
}

// WithDBTransaction injects the database connection into the context with a transaction.
func WithDBTransaction() Option {
	return func(ctx context.Context) (context.Context, error) {
		if hasDBInContext(ctx, databaseWithoutTransactionKey) {
			return nil, ErrConflict
		}
		conn, err := DB()
		if err != nil {
			return nil, err
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
