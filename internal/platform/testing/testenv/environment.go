package testenv

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"taskmanager/internal/platform/database"
	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/venomtest"

	"gorm.io/gorm"
)

// Environment holds the full test environment state.
// Can be used for both unit and integration tests.
type Environment struct {
	// Database
	DB        *gorm.DB
	container *dbtest.Container

	// HTTP Server
	server  *httptest.Server
	baseURL string

	// Venom Runner
	venomRunner *venomtest.Runner

	// Config
	config *Config
}

// Config holds the test environment configuration.
type Config struct {
	// Database options
	needsDB     bool
	containerDB *dbtest.Container
	optionsDB   []dbtest.Option

	// HTTP server options
	needsHTTP bool
	handler   http.Handler

	// Venom options
	needsVenom   bool
	venomOptions []venomtest.Option
}

// Option is a function that configures the Environment.
type Option func(*Config)

// Setup creates and configures a unified test environment.
func Setup(t *testing.T, opts ...Option) *Environment {
	t.Helper()

	cfg := &Config{}

	for _, opt := range opts {
		opt(cfg)
	}

	env := &Environment{
		config: cfg,
	}

	if cfg.needsDB {
		env.setupDatabase(t)
	}

	if cfg.needsHTTP {
		env.setupHTTPServer(t)
	}

	if cfg.needsVenom {
		env.setupVenomRunner(t)
	}

	return env
}

func (e *Environment) RunVenomSuite(t *testing.T, suitePath string) {
	t.Helper()
	e.venomRunner.Run(t, suitePath)
}

// setupDatabase configures the PostgreSQL container and connection.
func (e *Environment) setupDatabase(t *testing.T) {
	t.Helper()

	var container *dbtest.Container
	var err error

	if e.config.containerDB != nil {
		container = e.config.containerDB
	} else {
		container, err = dbtest.SetupDatabase(t, e.config.optionsDB...)
		if err != nil {
			t.Fatalf("failed to setup database: %v", err)
		}
	}

	e.container = container
	e.DB = container.DB()

	database.SetDB(e.DB)
}

// setupHTTPServer creates an httptest.Server for the test.
func (e *Environment) setupHTTPServer(t *testing.T) {
	t.Helper()

	if e.config.handler == nil {
		t.Fatal("HTTP handler is required when needsHTTP is true")
	}

	e.server = httptest.NewServer(e.config.handler)
	e.baseURL = e.server.URL

	t.Cleanup(func() {
		e.server.Close()
	})
}

// setupVenomRunner creates a Venom runner configured with this Environment's resources.
func (e *Environment) setupVenomRunner(t *testing.T) {
	t.Helper()

	if e.baseURL == "" {
		t.Fatal("VenomRunner requires HTTP server, use testenv.WithHTTPServer()")
	}

	e.venomRunner = venomtest.NewRunner(e.baseURL, e.config.venomOptions...)
}
