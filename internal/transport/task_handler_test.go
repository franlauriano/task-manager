//go:build test

package transport

import (
	"testing"

	"taskmanager/internal/paths"
	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/testenv"
	"taskmanager/internal/platform/testing/venomtest"
)

// Task Tests

func TestCreateTask(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(
			databaseTest,
			dbtest.WithMigrations(paths.MigrationDir()),
		),
		testenv.WithHTTPServer(Routes()),
		testenv.WithVenom(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name  string
		setup func()
		path  string
	}{
		// Success
		{"Create task with success (basic)", resetWithMinimalData, "success/tasks/create/basic.yml"},
		{"Create task with success (edge cases)", resetWithMinimalData, "success/tasks/create/edge_cases.yml"},
		{"Create task with success (corner cases)", resetWithMinimalData, "success/tasks/create/corner_cases.yml"},
		// Failure
		{"Create task with bad request", resetWithMinimalData, "failure/tasks/create/bad_request.yml"},
		{"Create task with validation errors", resetWithMinimalData, "failure/tasks/create/validation_errors.yml"},
		{"Create task with missing content type", resetWithMinimalData, "failure/tasks/create/missing_content_type.yml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			env.RunVenomSuite(t, tt.path)
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
		testenv.WithVenom(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name  string
		setup func()
		path  string
	}{
		// Success
		{"Update task with success (basic)", resetWithMinimalData, "success/tasks/update/basic.yml"},
		{"Update task with success (edge cases)", resetWithMinimalData, "success/tasks/update/edge_cases.yml"},
		{"Update task with success (corner cases)", resetWithMinimalData, "success/tasks/update/corner_cases.yml"},
		// Failure
		{"Update task with bad request", resetWithMinimalData, "failure/tasks/update/bad_request.yml"},
		{"Update task with validation errors", resetWithMinimalData, "failure/tasks/update/validation_errors.yml"},
		{"Update task with not found", resetWithMinimalData, "failure/tasks/update/not_found.yml"},
		{"Update task with missing content type", resetWithMinimalData, "failure/tasks/update/missing_content_type.yml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			env.RunVenomSuite(t, tt.path)
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
		testenv.WithVenom(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name  string
		setup func()
		path  string
	}{
		// Success
		{"Delete task with success (basic)", resetWithMinimalData, "success/tasks/delete/basic.yml"},
		{"Delete task with success (corner cases)", resetWithMinimalData, "success/tasks/delete/corner_cases.yml"},
		// Failure
		{"Delete task with bad request", resetWithMinimalData, "failure/tasks/delete/bad_request.yml"},
		{"Delete task with not found", resetWithMinimalData, "failure/tasks/delete/not_found.yml"},
		{"Delete task with missing content type", resetWithMinimalData, "failure/tasks/delete/missing_content_type.yml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			env.RunVenomSuite(t, tt.path)
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
		testenv.WithVenom(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name  string
		setup func()
		path  string
	}{
		// Success
		{"Retrieve task with success (basic)", resetWithMinimalData, "success/tasks/retrieve/basic.yml"},
		{"Retrieve task with success (corner cases)", resetWithMinimalData, "success/tasks/retrieve/corner_cases.yml"},
		// Failure
		{"Retrieve task with bad request", resetWithMinimalData, "failure/tasks/retrieve/bad_request.yml"},
		{"Retrieve task with not found", resetWithMinimalData, "failure/tasks/retrieve/not_found.yml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			env.RunVenomSuite(t, tt.path)
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
		testenv.WithVenom(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name  string
		setup func()
		path  string
	}{
		// Success
		{"List tasks with success (basic)", resetWithMinimalData, "success/tasks/list/basic.yml"},
		{"List tasks with success (edge cases)", resetWithMinimalData, "success/tasks/list/edge_cases.yml"},
		{"List tasks with success (corner cases)", resetWithMinimalData, "success/tasks/list/corner_cases.yml"},
		// Failure
		{"List tasks with bad request", resetWithMinimalData, "failure/tasks/list/bad_request.yml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			env.RunVenomSuite(t, tt.path)
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
		testenv.WithVenom(
			venomtest.WithSuiteRoot(paths.APITestDir()),
			venomtest.WithVerbose(1),
		),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name  string
		setup func()
		path  string
	}{
		// Success
		{"Update task status with success (basic)", resetWithMinimalData, "success/tasks/status/basic.yml"},
		// Failure
		{"Update task status with bad request", resetWithMinimalData, "failure/tasks/status/bad_request.yml"},
		{"Update task status with validation errors", resetWithMinimalData, "failure/tasks/status/validation_errors.yml"},
		{"Update task status with not found", resetWithMinimalData, "failure/tasks/status/not_found.yml"},
		{"Update task status with missing content type", resetWithMinimalData, "failure/tasks/status/missing_content_type.yml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			env.RunVenomSuite(t, tt.path)
		})
	}
}
