//go:build test

package team

import (
	"context"
	"testing"
	"time"

	"taskmanager/internal/entity/team"
	"taskmanager/internal/paths"
	"taskmanager/internal/platform/database"
	errs "taskmanager/internal/platform/errors"
	"taskmanager/internal/platform/testing/assert"
	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/testenv"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Test_datasource_Create(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name    string
		setup   func()
		ctx     context.Context
		team    *team.Team
		wantErr error
	}{
		{
			"Create team with success",
			resetWithMinimalData,
			context.Background(),
			&team.Team{
				Name:        "Novo time",
				Description: "Descrição do novo time",
			},
			nil,
		},
		{
			"Create team with context nil",
			resetWithMinimalData,
			nil,
			&team.Team{
				Name:        "Novo time",
				Description: "Descrição do novo time",
			},
			database.ErrContextDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{}
			err := p.Create(ctx, tt.team)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("datasource.Create() error diff: %s", diff)
				return
			}
		})
	}
}

func Test_datasource_RetrieveByUUID(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		teamUUID uuid.UUID
		want     *team.Team
		wantErr  error
	}{
		{
			"Retrieve team by UUID with success",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("111e4567-e89b-12d3-a456-426614174000"),
			&team.Team{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
				},
				UUID:        uuid.MustParse("111e4567-e89b-12d3-a456-426614174000"),
				Name:        "Time de Desenvolvimento",
				Description: "Equipe responsável pelo desenvolvimento de features e manutenção do código",
			},
			nil,
		},
		{
			"Retrieve team by UUID not found",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			nil,
			errs.ErrNotFound,
		},
		{
			"Retrieve team by UUID with context nil",
			nil,
			nil,
			uuid.MustParse("111e4567-e89b-12d3-a456-426614174000"),
			nil,
			database.ErrContextDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{}
			got, err := p.RetrieveByUUID(ctx, tt.teamUUID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("datasource.RetrieveByUUID() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("datasource.RetrieveByUUID() diff: %s", diff)
			}
		})
	}
}

func Test_datasource_ListPaginated(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name    string
		setup   func()
		ctx     context.Context
		page    int
		limit   int
		want    *team.ListTeams
		wantErr error
	}{
		{
			"ListPaginated all teams - page 1, limit 3",
			resetWithMinimalData,
			context.Background(),
			1,
			3,
			&team.ListTeams{
				Page:  1,
				Limit: 3,
				Teams: []team.Team{
					{
						Model: gorm.Model{
							ID:        4,
							CreatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
						},
						UUID:        uuid.MustParse("444e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time de UX/UI",
						Description: "Equipe responsável por design e experiência do usuário",
					},
					{
						Model: gorm.Model{
							ID:        3,
							CreatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
						},
						UUID:        uuid.MustParse("333e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time de QA",
						Description: "Equipe responsável por testes e garantia de qualidade",
					},
					{
						Model: gorm.Model{
							ID:        2,
							CreatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
						},
						UUID:        uuid.MustParse("222e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time de DevOps",
						Description: "Equipe responsável por infraestrutura, CI/CD e deploy",
					},
				},
				TotalItems: 4,
			},
			nil,
		},
		{
			"ListPaginated all teams - page 2, limit 3",
			resetWithMinimalData,
			context.Background(),
			2,
			3,
			&team.ListTeams{
				Page:  2,
				Limit: 3,
				Teams: []team.Team{
					{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
						},
						UUID:        uuid.MustParse("111e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time de Desenvolvimento",
						Description: "Equipe responsável pelo desenvolvimento de features e manutenção do código",
					},
				},
				TotalItems: 4,
			},
			nil,
		},
		{
			"ListPaginated all teams - page 1, limit 10 (all items)",
			resetWithMinimalData,
			context.Background(),
			1,
			10,
			&team.ListTeams{
				Page:  1,
				Limit: 10,
				Teams: []team.Team{
					{
						Model: gorm.Model{
							ID:        4,
							CreatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
						},
						UUID:        uuid.MustParse("444e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time de UX/UI",
						Description: "Equipe responsável por design e experiência do usuário",
					},
					{
						Model: gorm.Model{
							ID:        3,
							CreatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
						},
						UUID:        uuid.MustParse("333e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time de QA",
						Description: "Equipe responsável por testes e garantia de qualidade",
					},
					{
						Model: gorm.Model{
							ID:        2,
							CreatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
						},
						UUID:        uuid.MustParse("222e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time de DevOps",
						Description: "Equipe responsável por infraestrutura, CI/CD e deploy",
					},
					{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 20, 0, 0, time.UTC),
						},
						UUID:        uuid.MustParse("111e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time de Desenvolvimento",
						Description: "Equipe responsável pelo desenvolvimento de features e manutenção do código",
					},
				},
				TotalItems: 4,
			},
			nil,
		},
		{
			"ListPaginated all teams - page 2, limit 10 (empty page)",
			resetWithMinimalData,
			context.Background(),
			2,
			10,
			&team.ListTeams{
				Page:       2,
				Limit:      10,
				Teams:      []team.Team{},
				TotalItems: 4,
			},
			nil,
		},
		{
			"ListPaginated with context nil",
			nil,
			nil,
			1,
			10,
			nil,
			database.ErrContextDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{}
			got, err := p.ListPaginated(ctx, tt.page, tt.limit)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("datasource.ListPaginated() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("datasource.ListPaginated() diff: %s", diff)
			}

		})
	}
}

func Test_datasource_RetrieveTaskTeamID(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	teamID1 := uint(1)
	teamID2 := uint(2)

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		want     *uint
		wantErr  error
	}{
		{
			"RetrieveTaskTeamID with task associated to team",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			&teamID1,
			nil,
		},
		{
			"RetrieveTaskTeamID with task not associated to team",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			nil,
			nil,
		},
		{
			"RetrieveTaskTeamID with task associated to different team",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
			&teamID2,
			nil,
		},
		{
			"RetrieveTaskTeamID task not found",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			nil,
			errs.ErrNotFound,
		},
		{
			"RetrieveTaskTeamID with context nil",
			nil,
			nil,
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			nil,
			database.ErrContextDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{}
			got, err := p.RetrieveTaskTeamID(ctx, tt.taskUUID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("datasource.RetrieveTaskTeamID() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("datasource.RetrieveTaskTeamID() diff: %s", diff)
			}
		})
	}
}

func Test_datasource_UpdateTaskTeamID(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	teamID1 := uint(1)
	teamID2 := uint(2)

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		teamID   *uint
		wantErr  error
	}{
		{
			"UpdateTaskTeamID with success - associate task to team",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			&teamID1,
			nil,
		},
		{
			"UpdateTaskTeamID with success - disassociate task from team",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			nil,
			nil,
		},
		{
			"UpdateTaskTeamID with success - change team association",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			&teamID2,
			nil,
		},
		{
			"UpdateTaskTeamID task not found",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			&teamID1,
			errs.ErrNotFound,
		},
		{
			"UpdateTaskTeamID with context nil",
			nil,
			nil,
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			&teamID1,
			database.ErrContextDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{}
			err := p.UpdateTaskTeamID(ctx, tt.taskUUID, tt.teamID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("datasource.UpdateTaskTeamID() error diff: %s", diff)
				return
			}
		})
	}
}
