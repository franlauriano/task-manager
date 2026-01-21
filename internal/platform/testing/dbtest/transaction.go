package dbtest

import (
	"context"
	"testing"

	"taskmanager/internal/platform/database"
)

// SetupDBWithTransaction sets up a database transaction for the test, attaching it to ctx.
// If alias is empty, uses the default database alias.
func SetupDBWithTransaction(t *testing.T, ctx context.Context, alias string) context.Context {
	t.Helper()

	if ctx == nil {
		ctx = context.Background()
	}

	if alias == "" {
		alias = database.DatabaseDefaultAlias
	}

	var err error
	ctx, err = database.InjectDBsIntoContext(ctx, database.WithDBTransaction(alias))
	if err != nil {
		t.Fatalf("SetDBWithTransactionIntoContext: %v", err)
	}

	t.Cleanup(func() {
		database.Rollback(ctx, alias)
	})

	return ctx
}

// SetupDBWithoutTransaction sets up a database without transaction for the test, attaching it to ctx.
func SetupDBWithoutTransaction(t *testing.T, ctx context.Context, alias string) context.Context {
	t.Helper()

	if ctx == nil {
		ctx = context.Background()
	}
	var err error
	ctx, err = database.InjectDBsIntoContext(ctx, database.WithDBWithoutTransaction(alias))
	if err != nil {
		t.Fatalf("SetDBWithoutTransactionIntoContext: %v", err)
	}
	return ctx
}
