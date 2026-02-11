package redistest

// Option is a function that configures the Container.
type Option func(*Config)

// WithImage sets the Docker image for the Redis container.
func WithImage(image string) Option {
	return func(c *Config) {
		c.image = image
	}
}
