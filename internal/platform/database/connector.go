package database

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

// Connector encapsulates all database operations: connection lifecycle, context injection, and transaction management.
type Connector interface {
	DB() (*gorm.DB, error)
	InjectDBsIntoContext(ctx context.Context, options ...Option) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Close() error
}

// connector is the concrete implementation of Connector backed by a single *gorm.DB connection.
type connector struct {
	conn *gorm.DB
}

// NewConnector creates a Connector from an existing *gorm.DB connection.
func NewConnector(conn *gorm.DB) Connector {
	return &connector{conn: conn}
}

func (p *connector) DB() (*gorm.DB, error) {
	if p.conn == nil {
		return nil, ErrDBNotFound
	}
	return p.conn, nil
}

func (p *connector) InjectDBsIntoContext(ctx context.Context, options ...Option) (context.Context, error) {
	err := error(nil)
	for _, opt := range options {
		ctx, err = opt(ctx, p.conn)
		if err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

func (p *connector) Commit(ctx context.Context) error {
	conn, err := dbFromContext(ctx, databaseWithTransactionKey)
	if err != nil {
		return err
	}
	return conn.Commit().Error
}

func (p *connector) Rollback(ctx context.Context) error {
	conn, err := dbFromContext(ctx, databaseWithTransactionKey)
	if err != nil {
		return err
	}
	return conn.Rollback().Error
}

func (p *connector) Close() error {
	if p.conn == nil {
		return nil
	}
	sqlDB, err := p.conn.DB()
	if err != nil {
		return err
	}
	sqlDB.Close()
	p.conn = nil
	slog.Info("Database connection closed")
	return nil
}

// DBFromContext extracts the *gorm.DB from the context.
// It checks for a read-only connection first, then falls back to a transactional one.
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
