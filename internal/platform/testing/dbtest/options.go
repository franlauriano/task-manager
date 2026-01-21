package dbtest

// Option is a function that configures the Container.
type Option func(*Config)

// WithImage sets the Docker image for the PostgreSQL container.
func WithImage(image string) Option {
	return func(c *Config) {
		c.image = image
	}
}

// WithMigrations sets the migration directory to run on startup.
func WithMigrations(dir string) Option {
	return func(c *Config) {
		c.migrationDir = dir
	}
}
