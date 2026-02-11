package testenv

import (
	"net/http"

	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/venomtest"
)

// WithDatabase passes an existing database container to the test environment.
func WithDatabase(container *dbtest.Container, opts ...dbtest.Option) Option {
	return func(c *Config) {
		c.containerDB = container
		c.needsDB = true
		c.optionsDB = opts
	}
}

// WithNewDatabase creates a new PostgreSQL container for the test environment.
func WithNewDatabase(opts ...dbtest.Option) Option {
	return func(c *Config) {
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
