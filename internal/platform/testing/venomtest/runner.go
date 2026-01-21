package venomtest

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ovh/venom"
	"github.com/ovh/venom/executors"
)

// Runner runs Venom test suites.
type Runner struct {
	baseURL   string
	verbose   int
	variables map[string]interface{}
	suiteRoot string
}

// NewRunner creates a new Venom runner with the given base URL and options.
func NewRunner(baseURL string, opts ...Option) *Runner {
	r := &Runner{
		baseURL:   baseURL,
		verbose:   1,
		variables: make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(r)
	}

	r.variables["base_url"] = baseURL

	return r
}

// Run executes a Venom test suite from the given path (relative to suiteRoot or absolute).
func (r *Runner) Run(t *testing.T, suitePath string) {
	t.Helper()

	absPath := r.resolvePath(t, suitePath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		t.Fatalf("Venom suite not found: %s", absPath)
	}

	v := venom.New()
	v.Verbose = r.verbose
	v.StopOnFailure = false

	for name, executorFunc := range executors.Registry {
		v.RegisterExecutorBuiltin(name, executorFunc())
	}

	v.AddVariables(venom.H(r.variables))

	venom.InitTestLogger(t)

	ctx := context.Background()

	if err := v.Parse(ctx, []string{absPath}); err != nil {
		t.Fatalf("failed to parse venom suite: %v", err)
	}

	if len(v.Tests.TestSuites) == 0 {
		content, _ := os.ReadFile(absPath)
		t.Logf("YAML file content:\n%s", string(content))
		t.Fatal("no test suites parsed")
	}

	if err := v.Process(ctx, []string{absPath}); err != nil {
		t.Fatalf("failed to process venom suite: %v", err)
	}

	if len(v.Tests.TestSuites) == 0 {
		t.Fatal("no test suites executed")
	}

	r.reportResults(t, v)
}

// resolvePath resolves the suite path to an absolute path (using suiteRoot when relative).
func (r *Runner) resolvePath(t *testing.T, suitePath string) string {
	t.Helper()
	if filepath.IsAbs(suitePath) {
		return suitePath
	}
	if r.suiteRoot == "" {
		t.Fatal("WithSuiteRoot required for relative suite paths")
	}

	return filepath.Join(r.suiteRoot, suitePath)
}

// reportResults checks Venom results and fails the test with error details on failure.
func (r *Runner) reportResults(t *testing.T, v *venom.Venom) {
	t.Helper()

	ts := v.Tests.TestSuites[0]
	if ts.Status == venom.StatusFail {
		for _, tc := range ts.TestCases {
			if tc.Status == venom.StatusFail {
				for _, tsr := range tc.TestStepResults {
					if tsr.Status == venom.StatusFail {
						for _, failure := range tsr.Errors {
							t.Errorf("Test case '%s' failed: %s", tc.Name, failure.Value)
						}
					}
				}
			}
		}
		t.FailNow()
	}
}
