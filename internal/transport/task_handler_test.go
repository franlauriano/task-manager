//go:build test

package transport

import (
	"testing"

	"taskmanager/internal/paths"
	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/testenv"
	"taskmanager/internal/platform/testing/venomtest"
)

func resetWithMinimalData(env *testenv.Environment) {
	dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
}

func TestCreateTask(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(
			databaseTest,
			dbtest.WithMigrations(paths.MigrationDir()),
		),
		testenv.WithHTTPServer(Routes()),
		testenv.WithAPITest(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	tests := []struct {
		name      string
		setup     func()
		suitePath string
	}{
		// Success
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/tasks/create/basic.yml"},
		{"with success (edge cases)", func() { resetWithMinimalData(env) }, "success/tasks/create/edge_cases.yml"},
		{"with success (corner cases)", func() { resetWithMinimalData(env) }, "success/tasks/create/corner_cases.yml"},
		// Failure
		{"with bad request", func() { resetWithMinimalData(env) }, "failure/tasks/create/bad_request.yml"},
		{"with validation errors", func() { resetWithMinimalData(env) }, "failure/tasks/create/validation_errors.yml"},
		{"with missing content type", func() { resetWithMinimalData(env) }, "failure/tasks/create/missing_content_type.yml"},
	}

	for _, tc := range tests {
		t.Run("Create task "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(
			databaseTest,
			dbtest.WithMigrations(paths.MigrationDir()),
		),
		testenv.WithHTTPServer(Routes()),
		testenv.WithAPITest(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	tests := []struct {
		name      string
		setup     func()
		suitePath string
	}{
		// Success
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/tasks/update/basic.yml"},
		{"with success (edge cases)", func() { resetWithMinimalData(env) }, "success/tasks/update/edge_cases.yml"},
		{"with success (corner cases)", func() { resetWithMinimalData(env) }, "success/tasks/update/corner_cases.yml"},
		// Failure
		{"with bad request", func() { resetWithMinimalData(env) }, "failure/tasks/update/bad_request.yml"},
		{"with validation errors", func() { resetWithMinimalData(env) }, "failure/tasks/update/validation_errors.yml"},
		{"with not found", func() { resetWithMinimalData(env) }, "failure/tasks/update/not_found.yml"},
		{"with missing content type", func() { resetWithMinimalData(env) }, "failure/tasks/update/missing_content_type.yml"},
	}

	for _, tc := range tests {
		t.Run("Update task "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(
			databaseTest,
			dbtest.WithMigrations(paths.MigrationDir()),
		),
		testenv.WithHTTPServer(Routes()),
		testenv.WithAPITest(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	tests := []struct {
		name      string
		setup     func()
		suitePath string
	}{
		// Success
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/tasks/delete/basic.yml"},
		{"with success (corner cases)", func() { resetWithMinimalData(env) }, "success/tasks/delete/corner_cases.yml"},
		// Failure
		{"with bad request", func() { resetWithMinimalData(env) }, "failure/tasks/delete/bad_request.yml"},
		{"with not found", func() { resetWithMinimalData(env) }, "failure/tasks/delete/not_found.yml"},
		{"with missing content type", func() { resetWithMinimalData(env) }, "failure/tasks/delete/missing_content_type.yml"},
	}

	for _, tc := range tests {
		t.Run("Delete task "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}

func TestRetrieveTask(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(
			databaseTest,
			dbtest.WithMigrations(paths.MigrationDir()),
		),
		testenv.WithHTTPServer(Routes()),
		testenv.WithAPITest(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	tests := []struct {
		name      string
		setup     func()
		suitePath string
	}{
		// Success
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/tasks/retrieve/basic.yml"},
		{"with success (corner cases)", func() { resetWithMinimalData(env) }, "success/tasks/retrieve/corner_cases.yml"},
		// Failure
		{"with bad request", func() { resetWithMinimalData(env) }, "failure/tasks/retrieve/bad_request.yml"},
		{"with not found", func() { resetWithMinimalData(env) }, "failure/tasks/retrieve/not_found.yml"},
	}

	for _, tc := range tests {
		t.Run("Retrieve task "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}

func TestListTasks(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(
			databaseTest,
			dbtest.WithMigrations(paths.MigrationDir()),
		),
		testenv.WithHTTPServer(Routes()),
		testenv.WithAPITest(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	tests := []struct {
		name      string
		setup     func()
		suitePath string
	}{
		// Success
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/tasks/list/basic.yml"},
		{"with success (edge cases)", func() { resetWithMinimalData(env) }, "success/tasks/list/edge_cases.yml"},
		{"with success (corner cases)", func() { resetWithMinimalData(env) }, "success/tasks/list/corner_cases.yml"},
		// Failure
		{"with bad request", func() { resetWithMinimalData(env) }, "failure/tasks/list/bad_request.yml"},
	}

	for _, tc := range tests {
		t.Run("List tasks "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(
			databaseTest,
			dbtest.WithMigrations(paths.MigrationDir()),
		),
		testenv.WithHTTPServer(Routes()),
		testenv.WithAPITest(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	tests := []struct {
		name      string
		setup     func()
		suitePath string
	}{
		// Success
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/tasks/status/basic.yml"},
		// Failure
		{"with bad request", func() { resetWithMinimalData(env) }, "failure/tasks/status/bad_request.yml"},
		{"with validation errors", func() { resetWithMinimalData(env) }, "failure/tasks/status/validation_errors.yml"},
		{"with not found", func() { resetWithMinimalData(env) }, "failure/tasks/status/not_found.yml"},
		{"with missing content type", func() { resetWithMinimalData(env) }, "failure/tasks/status/missing_content_type.yml"},
	}

	for _, tc := range tests {
		t.Run("Update task status "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}
