//go:build test

package team

import (
	"context"
	"errors"
	"strings"
	"testing"

	taskEntity "taskmanager/internal/entity/task"
	teamEntity "taskmanager/internal/entity/team"
	"taskmanager/internal/platform/database"
	errs "taskmanager/internal/platform/errors"
	"taskmanager/internal/platform/testing/assert"
	taskRepo "taskmanager/internal/repository/task"
	teamRepo "taskmanager/internal/repository/team"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	originalPersist := teamRepo.Persist()

	tests := []struct {
		name    string
		setup   func()
		ctx     context.Context
		team    *teamEntity.Team
		wantErr error
	}{
		{
			"Create team with success",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "Novo time",
				Description: "Descrição do novo time",
			},
			nil,
		},
		{
			"Create team with empty name",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "",
				Description: "Descrição válida",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "name",
						Message: "name is required",
					},
				},
			},
		},
		{
			"Create team with empty description",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "Nome válido",
				Description: "",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Create team with only whitespace name",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "   ",
				Description: "Descrição válida",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "name",
						Message: "name is required",
					},
				},
			},
		},
		{
			"Create team with only whitespace description",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "Nome válido",
				Description: "   ",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Create team with empty name and description",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "",
				Description: "",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "name",
						Message: "name is required",
					},
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Create team with only whitespace name and description",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "   ",
				Description: "   ",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "name",
						Message: "name is required",
					},
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Create team with tab and newline whitespace in name",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "\t\n\r   ",
				Description: "Descrição válida",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "name",
						Message: "name is required",
					},
				},
			},
		},
		{
			"Create team with tab and newline whitespace in description",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "Nome válido",
				Description: "\t\n\r   ",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Create team with persist error",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return database.ErrContextDatabase },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "Novo time",
				Description: "Descrição válida",
			},
			database.ErrContextDatabase,
		},
		{
			"Create team with generic persist error",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return errors.New("database connection failed") },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        "Novo time",
				Description: "Descrição válida",
			},
			errors.New("database connection failed"),
		},
		{
			"Create team with name exceeding 255 characters",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        strings.Repeat("a", 256),
				Description: "Descrição válida",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "name",
						Message: "name must not exceed 255 characters",
					},
				},
			},
		},
		{
			"Create team with name exactly 255 characters",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *teamEntity.Team) error { return nil },
				})
			},
			context.Background(),
			&teamEntity.Team{
				Name:        strings.Repeat("a", 255),
				Description: "Descrição válida",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				teamRepo.SetPersist(originalPersist)
			}()

			if tt.setup != nil {
				tt.setup()
			}

			err := Create(tt.ctx, tt.team)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("Create() error diff: %s", diff)
				return
			}
		})
	}
}

func TestRetrieveByUUIDWithTasks(t *testing.T) {
	originalPersist := teamRepo.Persist()
	originalTaskPersist := taskRepo.Persist()

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		teamUUID uuid.UUID
		wantTeam *teamEntity.Team
		wantErr  error
	}{
		{
			"RetrieveByUUIDWithTasks with success",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
				})
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListByTeamID: func(ctx context.Context, teamID uint) ([]taskEntity.Task, error) {
						return []taskEntity.Task{
							{
								UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
								Title:       "Tarefa 1 do Team",
								Description: "Descrição da tarefa 1",
								Status:      taskEntity.StatusTodo,
								TeamID:      &teamID,
							},
							{
								UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
								Title:       "Tarefa 2 do Team",
								Description: "Descrição da tarefa 2",
								Status:      taskEntity.StatusInProgress,
								TeamID:      &teamID,
							},
						}, nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			&teamEntity.Team{
				UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name:        "Time de Desenvolvimento",
				Description: "Time responsável pelo desenvolvimento",
				Model:       gorm.Model{ID: 1},
				Tasks: []taskEntity.Task{
					{
						UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
						Title:       "Tarefa 1 do Team",
						Description: "Descrição da tarefa 1",
						Status:      taskEntity.StatusTodo,
						TeamID:      func() *uint { id := uint(1); return &id }(),
					},
					{
						UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
						Title:       "Tarefa 2 do Team",
						Description: "Descrição da tarefa 2",
						Status:      taskEntity.StatusInProgress,
						TeamID:      func() *uint { id := uint(1); return &id }(),
					},
				},
			},
			nil,
		},
		{
			"RetrieveByUUIDWithTasks with success - empty tasks list",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
				})
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListByTeamID: func(ctx context.Context, teamID uint) ([]taskEntity.Task, error) {
						return []taskEntity.Task{}, nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			&teamEntity.Team{
				UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name:        "Time de Desenvolvimento",
				Description: "Time responsável pelo desenvolvimento",
				Model:       gorm.Model{ID: 1},
				Tasks:       []taskEntity.Task{},
			},
			nil,
		},
		{
			"RetrieveByUUIDWithTasks with task.ListByTeamID error",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
				})
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListByTeamID: func(ctx context.Context, teamID uint) ([]taskEntity.Task, error) {
						return nil, database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			nil,
			database.ErrContextDatabase,
		},
		{
			"RetrieveByUUIDWithTasks team not found",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return nil, errs.ErrNotFound
					},
				})
			},
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			nil,
			errs.ErrNotFound,
		},
		{
			"RetrieveByUUIDWithTasks with context database error",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return nil, database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			nil,
			database.ErrContextDatabase,
		},
		{
			"RetrieveByUUIDWithTasks with generic persist error",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return nil, errors.New("database connection failed")
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			nil,
			errors.New("database connection failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				teamRepo.SetPersist(originalPersist)
				taskRepo.SetPersist(originalTaskPersist)
			}()

			if tt.setup != nil {
				tt.setup()
			}

			gotTeam, err := RetrieveByUUIDWithTasks(tt.ctx, tt.teamUUID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("RetrieveByUUIDWithTasks() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(gotTeam, tt.wantTeam); diff != "" {
				t.Errorf("RetrieveByUUIDWithTasks() gotTeam diff: %s", diff)
			}
		})
	}
}

func TestListPaginated(t *testing.T) {
	originalPersist := teamRepo.Persist()
	originalConfig := Config

	tests := []struct {
		name    string
		setup   func()
		ctx     context.Context
		page    int
		limit   int
		want    *teamEntity.ListTeams
		wantErr error
	}{
		{
			"ListPaginated all teams with success",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, page, limit int) (*teamEntity.ListTeams, error) {
						return &teamEntity.ListTeams{
							Page:  1,
							Limit: 10,
							Teams: []teamEntity.Team{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
									Name:        "Time 1",
									Description: "Descrição 1",
								},
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
									Name:        "Time 2",
									Description: "Descrição 2",
								},
							},
							TotalItems: 2,
						}, nil
					},
				})
			},
			context.Background(),
			1,
			10,
			&teamEntity.ListTeams{
				Page:  1,
				Limit: 10,
				Teams: []teamEntity.Team{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time 1",
						Description: "Descrição 1",
					},
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
						Name:        "Time 2",
						Description: "Descrição 2",
					},
				},
				TotalItems: 2,
			},
			nil,
		},
		{
			"ListPaginated with pagination - page 2, limit 2",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, page, limit int) (*teamEntity.ListTeams, error) {
						return &teamEntity.ListTeams{
							Page:  2,
							Limit: 2,
							Teams: []teamEntity.Team{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
									Name:        "Time 3",
									Description: "Descrição 3",
								},
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"),
									Name:        "Time 4",
									Description: "Descrição 4",
								},
							},
							TotalItems: 4,
						}, nil
					},
				})
			},
			context.Background(),
			2,
			2,
			&teamEntity.ListTeams{
				Page:  2,
				Limit: 2,
				Teams: []teamEntity.Team{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
						Name:        "Time 3",
						Description: "Descrição 3",
					},
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"),
						Name:        "Time 4",
						Description: "Descrição 4",
					},
				},
				TotalItems: 4,
			},
			nil,
		},
		{
			"ListPaginated with empty result",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, page, limit int) (*teamEntity.ListTeams, error) {
						return &teamEntity.ListTeams{
							Page:       1,
							Limit:      10,
							Teams:      []teamEntity.Team{},
							TotalItems: 0,
						}, nil
					},
				})
			},
			context.Background(),
			1,
			10,
			&teamEntity.ListTeams{
				Page:       1,
				Limit:      10,
				Teams:      []teamEntity.Team{},
				TotalItems: 0,
			},
			nil,
		},
		{
			"ListPaginated with context database error",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, page, limit int) (*teamEntity.ListTeams, error) {
						return nil, database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			1,
			10,
			nil,
			database.ErrContextDatabase,
		},
		{
			"ListPaginated with generic persist error",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, page, limit int) (*teamEntity.ListTeams, error) {
						return nil, errors.New("database connection failed")
					},
				})
			},
			context.Background(),
			1,
			10,
			nil,
			errors.New("database connection failed"),
		},
		{
			"ListPaginated with limit 0 uses default limit",
			func() {
				Config.ListDefaultLimit = 15
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, page, limit int) (*teamEntity.ListTeams, error) {
						return &teamEntity.ListTeams{
							Page:  1,
							Limit: 15,
							Teams: []teamEntity.Team{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
									Name:        "Time 1",
									Description: "Descrição 1",
								},
							},
							TotalItems: 1,
						}, nil
					},
				})
			},
			context.Background(),
			1,
			0,
			&teamEntity.ListTeams{
				Page:  1,
				Limit: 15,
				Teams: []teamEntity.Team{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time 1",
						Description: "Descrição 1",
					},
				},
				TotalItems: 1,
			},
			nil,
		},
		{
			"ListPaginated with negative limit uses default limit",
			func() {
				Config.ListDefaultLimit = 15
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, page, limit int) (*teamEntity.ListTeams, error) {
						return &teamEntity.ListTeams{
							Page:  1,
							Limit: 15,
							Teams: []teamEntity.Team{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
									Name:        "Time 1",
									Description: "Descrição 1",
								},
							},
							TotalItems: 1,
						}, nil
					},
				})
			},
			context.Background(),
			1,
			-5,
			&teamEntity.ListTeams{
				Page:  1,
				Limit: 15,
				Teams: []teamEntity.Team{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time 1",
						Description: "Descrição 1",
					},
				},
				TotalItems: 1,
			},
			nil,
		},
		{
			"ListPaginated with limit exceeding max uses max limit",
			func() {
				Config.ListDefaultLimit = 10
				Config.ListMaxLimit = 20
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, page, limit int) (*teamEntity.ListTeams, error) {
						return &teamEntity.ListTeams{
							Page:  1,
							Limit: 20,
							Teams: []teamEntity.Team{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
									Name:        "Time 1",
									Description: "Descrição 1",
								},
							},
							TotalItems: 1,
						}, nil
					},
				})
			},
			context.Background(),
			1,
			100,
			&teamEntity.ListTeams{
				Page:  1,
				Limit: 20,
				Teams: []teamEntity.Team{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Name:        "Time 1",
						Description: "Descrição 1",
					},
				},
				TotalItems: 1,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				teamRepo.SetPersist(originalPersist)
				Config = originalConfig
			}()

			if tt.setup != nil {
				tt.setup()
			}

			got, err := ListPaginated(tt.ctx, tt.page, tt.limit)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("ListPaginated() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ListPaginated() diff: %s", diff)
			}
		})
	}
}

func TestAssociateTask(t *testing.T) {
	originalPersist := teamRepo.Persist()

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		teamUUID uuid.UUID
		taskUUID uuid.UUID
		wantErr  error
	}{
		{
			"AssociateTask with success",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return nil, nil // Task não está associada a nenhum team
					},
					FnUpdateTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID, teamID *uint) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			nil,
		},
		{
			"AssociateTask team not found",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return nil, errs.ErrNotFound
					},
				})
			},
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "team",
						Message: "team not found",
					},
				},
			},
		},
		{
			"AssociateTask task not found",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return nil, errs.ErrNotFound
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "task",
						Message: "task not found",
					},
				},
			},
		},
		{
			"AssociateTask task already associated with another team",
			func() {
				teamID := uint(1)
				otherTeamID := uint(2)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return &otherTeamID, nil // Task já está associada a outro team
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "task",
						Message: "task is already associated with another team",
					},
				},
			},
		},
		{
			"AssociateTask with retrieve team context database error",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return nil, database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			database.ErrContextDatabase,
		},
		{
			"AssociateTask with retrieve task team ID context database error",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return nil, database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			database.ErrContextDatabase,
		},
		{
			"AssociateTask with update task team ID context database error",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return nil, nil
					},
					FnUpdateTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID, teamID *uint) error {
						return database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			database.ErrContextDatabase,
		},
		{
			"AssociateTask with generic persist error",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return nil, nil
					},
					FnUpdateTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID, teamID *uint) error {
						return errors.New("database connection failed")
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			errors.New("database connection failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				teamRepo.SetPersist(originalPersist)
			}()

			if tt.setup != nil {
				tt.setup()
			}

			err := AssociateTask(tt.ctx, tt.teamUUID, tt.taskUUID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("AssociateTask() error diff: %s", diff)
				return
			}
		})
	}
}

func TestDisassociateTask(t *testing.T) {
	originalPersist := teamRepo.Persist()

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		teamUUID uuid.UUID
		taskUUID uuid.UUID
		wantErr  error
	}{
		{
			"DisassociateTask with success",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return &teamID, nil // Task está associada a este team
					},
					FnUpdateTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID, teamID *uint) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			nil,
		},
		{
			"DisassociateTask team not found",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return nil, errs.ErrNotFound
					},
				})
			},
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "team",
						Message: "team not found",
					},
				},
			},
		},
		{
			"DisassociateTask task not found",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return nil, errs.ErrNotFound
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "task",
						Message: "task not found",
					},
				},
			},
		},
		{
			"DisassociateTask task not associated with this team",
			func() {
				teamID := uint(1)
				otherTeamID := uint(2)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return &otherTeamID, nil // Task está associada a outro team
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "task",
						Message: "task is not associated with this team",
					},
				},
			},
		},
		{
			"DisassociateTask task not associated with any team",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return nil, nil // Task não está associada a nenhum team
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "task",
						Message: "task is not associated with this team",
					},
				},
			},
		},
		{
			"DisassociateTask with retrieve team context database error",
			func() {
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return nil, database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			database.ErrContextDatabase,
		},
		{
			"DisassociateTask with retrieve task team ID context database error",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return nil, database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			database.ErrContextDatabase,
		},
		{
			"DisassociateTask with update task team ID context database error",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return &teamID, nil
					},
					FnUpdateTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID, teamID *uint) error {
						return database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			database.ErrContextDatabase,
		},
		{
			"DisassociateTask with generic persist error",
			func() {
				teamID := uint(1)
				teamRepo.SetPersist(&teamRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
						return &teamEntity.Team{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Name:        "Time de Desenvolvimento",
							Description: "Time responsável pelo desenvolvimento",
							Model:       gorm.Model{ID: teamID},
						}, nil
					},
					FnRetrieveTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
						return &teamID, nil
					},
					FnUpdateTaskTeamID: func(ctx context.Context, taskUUID uuid.UUID, teamID *uint) error {
						return errors.New("database connection failed")
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
			errors.New("database connection failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				teamRepo.SetPersist(originalPersist)
			}()

			if tt.setup != nil {
				tt.setup()
			}

			err := DisassociateTask(tt.ctx, tt.teamUUID, tt.taskUUID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("DisassociateTask() error diff: %s", diff)
				return
			}
		})
	}
}
