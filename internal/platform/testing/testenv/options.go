package testenv

import (
	"net/http"

	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/redistest"
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

// WithRedis passes an existing Redis container to the test environment.
func WithRedis(container *redistest.Container) Option {
	return func(c *Config) {
		c.containerRedis = container
		c.needsRedis = true
	}
}

// WithNewRedis creates a new Redis container for the test environment.
func WithNewRedis(opts ...redistest.Option) Option {
	return func(c *Config) {
		c.needsRedis = true
		c.optionsRedis = opts
	}
}

// WithHTTPServer enables an HTTP test server; the given handler is used to create it.
func WithHTTPServer(handler http.Handler) Option {
	return func(c *Config) {
		c.needsHTTP = true
		c.handler = handler
	}
}

// WithAPITest enables API test suites using Venom.
func WithAPITest(opts ...venomtest.Option) Option {
	return func(c *Config) {
		c.needsVenom = true
		c.venomOptions = opts
	}
}
