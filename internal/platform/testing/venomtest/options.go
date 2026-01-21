package venomtest

// Option is a function that configures the Venom runner.
type Option func(*Runner)

// WithSuiteRoot sets the root directory for resolving relative suite paths.
func WithSuiteRoot(suiteRoot string) Option {
	return func(r *Runner) {
		r.suiteRoot = suiteRoot
	}
}

// WithVerbose sets the Venom verbosity level (0–3).
func WithVerbose(level int) Option {
	return func(r *Runner) {
		r.verbose = level
	}
}

// WithVariables sets the variables for the Venom runner.
func WithVariables(variables map[string]interface{}) Option {
	return func(r *Runner) {
		r.variables = variables
	}
}
