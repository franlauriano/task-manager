package dbtest

import (
	"context"
	"testing"

	"taskmanager/internal/platform/database"
)

// SetupDBWithTransaction sets up a database transaction for the test, attaching it to ctx.
func SetupDBWithTransaction(t *testing.T, ctx context.Context) context.Context {
	t.Helper()

	if ctx == nil {
		ctx = context.Background()
	}

	var err error
	ctx, err = database.InjectDBsIntoContext(ctx, database.WithDBTransaction())
	if err != nil {
		t.Fatalf("SetDBWithTransactionIntoContext: %v", err)
	}

	t.Cleanup(func() {
		database.Rollback(ctx)
	})

	return ctx
}

// SetupDBWithoutTransaction sets up a database without transaction for the test, attaching it to ctx.
func SetupDBWithoutTransaction(t *testing.T, ctx context.Context) context.Context {
	t.Helper()

	if ctx == nil {
		ctx = context.Background()
	}
	var err error
	ctx, err = database.InjectDBsIntoContext(ctx, database.WithDBWithoutTransaction())
	if err != nil {
		t.Fatalf("SetDBWithoutTransactionIntoContext: %v", err)
	}
	return ctx
}
