//go:build test

package transport

import (
	"testing"

	"taskmanager/internal/paths"
	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/testenv"
	"taskmanager/internal/platform/testing/venomtest"
)

func TestCreateTeam(t *testing.T) {
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
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/teams/create/basic.yml"},
		{"with success (edge cases)", func() { resetWithMinimalData(env) }, "success/teams/create/edge_cases.yml"},
		{"with success (corner cases)", func() { resetWithMinimalData(env) }, "success/teams/create/corner_cases.yml"},
		// Failure
		{"with bad request", func() { resetWithMinimalData(env) }, "failure/teams/create/bad_request.yml"},
		{"with validation errors", func() { resetWithMinimalData(env) }, "failure/teams/create/validation_errors.yml"},
		{"with missing content type", func() { resetWithMinimalData(env) }, "failure/teams/create/missing_content_type.yml"},
	}

	for _, tc := range tests {
		t.Run("Create team "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}

func TestListTeams(t *testing.T) {
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
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/teams/list/basic.yml"},
		{"with success (edge cases)", func() { resetWithMinimalData(env) }, "success/teams/list/edge_cases.yml"},
		{"with success (corner cases)", func() { resetWithMinimalData(env) }, "success/teams/list/corner_cases.yml"},
	}

	for _, tc := range tests {
		t.Run("List teams "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}

func TestRetrieveTeam(t *testing.T) {
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
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/teams/retrieve/basic.yml"},
		{"with success (corner cases)", func() { resetWithMinimalData(env) }, "success/teams/retrieve/corner_cases.yml"},
		// Failure
		{"with bad request", func() { resetWithMinimalData(env) }, "failure/teams/retrieve/bad_request.yml"},
		{"with not found", func() { resetWithMinimalData(env) }, "failure/teams/retrieve/not_found.yml"},
	}

	for _, tc := range tests {
		t.Run("Retrieve team "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}

func TestAssociateTaskToTeam(t *testing.T) {
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
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/teams/associate_task/basic.yml"},
		{"with success (corner cases)", func() { resetWithMinimalData(env) }, "success/teams/associate_task/corner_cases.yml"},
		// Failure
		{"with bad request", func() { resetWithMinimalData(env) }, "failure/teams/associate_task/bad_request.yml"},
		{"with validation errors", func() { resetWithMinimalData(env) }, "failure/teams/associate_task/validation_errors.yml"},
		{"with missing content type", func() { resetWithMinimalData(env) }, "failure/teams/associate_task/missing_content_type.yml"},
	}

	for _, tc := range tests {
		t.Run("Associate task to team "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}

func TestDisassociateTaskFromTeam(t *testing.T) {
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
		{"with success (basic)", func() { resetWithMinimalData(env) }, "success/teams/disassociate_task/basic.yml"},
		{"with success (corner cases)", func() { resetWithMinimalData(env) }, "success/teams/disassociate_task/corner_cases.yml"},
		// Failure
		{"with bad request", func() { resetWithMinimalData(env) }, "failure/teams/disassociate_task/bad_request.yml"},
		{"with validation errors", func() { resetWithMinimalData(env) }, "failure/teams/disassociate_task/validation_errors.yml"},
		{"with missing content type", func() { resetWithMinimalData(env) }, "failure/teams/disassociate_task/missing_content_type.yml"},
	}

	for _, tc := range tests {
		t.Run("Disassociate task from team "+tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			env.RunAPISuite(t, tc.suitePath)
		})
	}
}
