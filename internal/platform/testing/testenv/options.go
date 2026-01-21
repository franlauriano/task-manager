package testenv

import (
	"net/http"

	"taskmanager/internal/platform/database"
	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/venomtest"
)

// WithContainerDatabaseAlias sets the database alias and passes an existing database container to the test environment.
func WithContainerDatabaseAlias(alias string, container *dbtest.Container, opts ...dbtest.Option) Option {
	return func(c *Config) {
		c.aliasDB = alias
		c.containerDB = container
		c.needsDB = true
		c.optionsDB = opts
	}
}

// WithContainerDatabase sets the database alias to the default alias and passes an existing database container to the test environment.
func WithContainerDatabase(container *dbtest.Container, opts ...dbtest.Option) Option {
	return func(c *Config) {
		c.aliasDB = database.DatabaseDefaultAlias
		c.containerDB = container
		c.needsDB = true
		c.optionsDB = opts
	}
}

// WithDatabaseAlias sets the database alias and enables a PostgreSQL database in the test environment.
func WithDatabaseAlias(alias string, opts ...dbtest.Option) Option {
	return func(c *Config) {
		c.aliasDB = alias
		c.needsDB = true
		c.optionsDB = opts
	}
}

// WithDatabase sets the database alias to the default alias and enables a PostgreSQL database in the test environment.
func WithDatabase(opts ...dbtest.Option) Option {
	return func(c *Config) {
		c.aliasDB = database.DatabaseDefaultAlias
		c.needsDB = true
		c.optionsDB = opts
	}
}

// WithHTTPServer enables an HTTP test server; the given handler is used to create it.
func WithHTTPServer(handler http.Handler) Option {
	return func(c *Config) {
		c.needsHTTP = true
		c.handler = handler
	}
}

// WithVenom enables Venom test suites.
func WithVenom(opts ...venomtest.Option) Option {
	return func(c *Config) {
		c.needsVenom = true
		c.venomOptions = opts
	}
}
