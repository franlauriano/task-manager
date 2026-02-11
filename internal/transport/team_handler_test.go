//go:build test

package transport

import (
	"testing"

	"taskmanager/internal/paths"
	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/testenv"
	"taskmanager/internal/platform/testing/venomtest"
)

// Team Tests

func TestCreateTeam(t *testing.T) {
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
		{"Create team with success (basic)", resetWithMinimalData, "success/teams/create/basic.yml"},
		{"Create team with success (edge cases)", resetWithMinimalData, "success/teams/create/edge_cases.yml"},
		{"Create team with success (corner cases)", resetWithMinimalData, "success/teams/create/corner_cases.yml"},
		// Failure
		{"Create team with bad request", resetWithMinimalData, "failure/teams/create/bad_request.yml"},
		{"Create team with validation errors", resetWithMinimalData, "failure/teams/create/validation_errors.yml"},
		{"Create team with missing content type", resetWithMinimalData, "failure/teams/create/missing_content_type.yml"},
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

func TestListTeams(t *testing.T) {
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
		{"List teams with success (basic)", resetWithMinimalData, "success/teams/list/basic.yml"},
		{"List teams with success (edge cases)", resetWithMinimalData, "success/teams/list/edge_cases.yml"},
		{"List teams with success (corner cases)", resetWithMinimalData, "success/teams/list/corner_cases.yml"},
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

func TestRetrieveTeam(t *testing.T) {
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
		{"Retrieve team with success (basic)", resetWithMinimalData, "success/teams/retrieve/basic.yml"},
		{"Retrieve team with success (corner cases)", resetWithMinimalData, "success/teams/retrieve/corner_cases.yml"},
		// Failure
		{"Retrieve team with bad request", resetWithMinimalData, "failure/teams/retrieve/bad_request.yml"},
		{"Retrieve team with not found", resetWithMinimalData, "failure/teams/retrieve/not_found.yml"},
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

func TestAssociateTaskToTeam(t *testing.T) {
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
		{"Associate task to team with success (basic)", resetWithMinimalData, "success/teams/associate_task/basic.yml"},
		{"Associate task to team with success (corner cases)", resetWithMinimalData, "success/teams/associate_task/corner_cases.yml"},
		// Failure
		{"Associate task to team with bad request", resetWithMinimalData, "failure/teams/associate_task/bad_request.yml"},
		{"Associate task to team with validation errors", resetWithMinimalData, "failure/teams/associate_task/validation_errors.yml"},
		{"Associate task to team with missing content type", resetWithMinimalData, "failure/teams/associate_task/missing_content_type.yml"},
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

func TestDisassociateTaskFromTeam(t *testing.T) {
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
		{"Disassociate task from team with success (basic)", resetWithMinimalData, "success/teams/disassociate_task/basic.yml"},
		{"Disassociate task from team with success (corner cases)", resetWithMinimalData, "success/teams/disassociate_task/corner_cases.yml"},
		// Failure
		{"Disassociate task from team with bad request", resetWithMinimalData, "failure/teams/disassociate_task/bad_request.yml"},
		{"Disassociate task from team with validation errors", resetWithMinimalData, "failure/teams/disassociate_task/validation_errors.yml"},
		{"Disassociate task from team with missing content type", resetWithMinimalData, "failure/teams/disassociate_task/missing_content_type.yml"},
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
