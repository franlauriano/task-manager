//go:build test

package task

import (
	"context"
	"errors"
	"strings"
	"testing"

	taskEntity "taskmanager/internal/entity/task"
	"taskmanager/internal/platform/database"
	errs "taskmanager/internal/platform/errors"
	"taskmanager/internal/platform/testing/assert"
	taskRepo "taskmanager/internal/repository/task"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	originalPersist := taskRepo.Persist()

	tests := []struct {
		name    string
		setup   func()
		ctx     context.Context
		task    *taskEntity.Task
		wantErr error
	}{
		{
			"Create task with success",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "Nova tarefa",
				Description: "Descrição da nova tarefa",
			},
			nil,
		},
		{
			"Create task with empty title",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "",
				Description: "Descrição válida",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
				},
			},
		},
		{
			"Create task with empty description",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "Título válido",
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
			"Create task with only whitespace title",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "   ",
				Description: "Descrição válida",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
				},
			},
		},
		{
			"Create task with only whitespace description",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "Título válido",
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
			"Create task with empty title and description",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "",
				Description: "",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Create task with only whitespace title and description",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "   ",
				Description: "   ",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Create task with tab and newline whitespace in title",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "\t\n\r   ",
				Description: "Descrição válida",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
				},
			},
		},
		{
			"Create task with tab and newline whitespace in description",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "Título válido",
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
			"Create task with persist error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return database.ErrContextDatabase },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "Nova tarefa",
				Description: "Descrição válida",
			},
			database.ErrContextDatabase,
		},
		{
			"Create task with generic persist error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return errors.New("database connection failed") },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       "Nova tarefa",
				Description: "Descrição válida",
			},
			errors.New("database connection failed"),
		},
		{
			"Create task with title exceeding 255 characters",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       strings.Repeat("a", 256),
				Description: "Descrição válida",
			},
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "title",
						Message: "title must not exceed 255 characters",
					},
				},
			},
		},
		{
			"Create task with title exactly 255 characters",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnCreate: func(ctx context.Context, t *taskEntity.Task) error { return nil },
				})
			},
			context.Background(),
			&taskEntity.Task{
				Title:       strings.Repeat("a", 255),
				Description: "Descrição válida",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				taskRepo.SetPersist(originalPersist)
			}()

			if tt.setup != nil {
				tt.setup()
			}

			err := Create(tt.ctx, tt.task)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("Create() error diff: %s", diff)
				return
			}
		})
	}
}

func TestRetrieveByUUID(t *testing.T) {
	originalPersist := taskRepo.Persist()

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		want     *taskEntity.Task
		wantErr  error
	}{
		{
			"Retrieve task by UUID with success",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Implementar autenticação",
							Description: "Criar sistema de autenticação JWT para a API",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			&taskEntity.Task{
				UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:       "Implementar autenticação",
				Description: "Criar sistema de autenticação JWT para a API",
				Status:      taskEntity.StatusTodo,
			},
			nil,
		},
		{
			"Retrieve task by UUID not found",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
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
			"Retrieve task by UUID with context database error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
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
			"Retrieve task by UUID with generic persist error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
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
				taskRepo.SetPersist(originalPersist)
			}()
			if tt.setup != nil {
				tt.setup()
			}

			got, err := RetrieveByUUID(tt.ctx, tt.taskUUID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("RetrieveByUUID() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("RetrieveByUUID() diff: %s", diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	originalPersist := taskRepo.Persist()

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		updates  map[string]any
		want     *taskEntity.Task
		wantErr  error
	}{
		{
			"Update task with success - title and description",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title":       "Título atualizado",
				"description": "Descrição atualizada",
			},
			&taskEntity.Task{
				UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:       "Título atualizado",
				Description: "Descrição atualizada",
				Status:      taskEntity.StatusTodo,
			},
			nil,
		},
		{
			"Update task with success - only title",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title": "Título atualizado",
			},
			&taskEntity.Task{
				UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:       "Título atualizado",
				Description: "Descrição original",
				Status:      taskEntity.StatusTodo,
			},
			nil,
		},
		{
			"Update task with success - only description",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"description": "Descrição atualizada",
			},
			&taskEntity.Task{
				UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:       "Título original",
				Description: "Descrição atualizada",
				Status:      taskEntity.StatusTodo,
			},
			nil,
		},
		{
			"Update task with success - trim whitespace",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title":       "  Título com espaços  ",
				"description": "  Descrição com espaços  ",
			},
			&taskEntity.Task{
				UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:       "Título com espaços",
				Description: "Descrição com espaços",
				Status:      taskEntity.StatusTodo,
			},
			nil,
		},
		{
			"Update task with empty title",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title": "",
			},
			nil,
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
				},
			},
		},
		{
			"Update task with empty description",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"description": "",
			},
			nil,
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
			"Update task with empty title and description",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title":       "",
				"description": "",
			},
			nil,
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Update task with only whitespace title",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title": "   ",
			},
			nil,
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
				},
			},
		},
		{
			"Update task with only whitespace description",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"description": "   ",
			},
			nil,
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
			"Update task not found",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return nil, errs.ErrNotFound
					},
				})
			},
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			map[string]any{
				"title": "Título atualizado",
			},
			nil,
			errs.ErrNotFound,
		},
		{
			"Update task with retrieve context database error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return nil, database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title": "Título atualizado",
			},
			nil,
			database.ErrContextDatabase,
		},
		{
			"Update task with update persist error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title": "Título atualizado",
			},
			nil,
			database.ErrContextDatabase,
		},
		{
			"Update task with generic persist error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return errors.New("database connection failed")
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title": "Título atualizado",
			},
			nil,
			errors.New("database connection failed"),
		},
		{
			"Update task with title exceeding 255 characters",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title": strings.Repeat("a", 256),
			},
			nil,
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "title",
						Message: "title must not exceed 255 characters",
					},
				},
			},
		},
		{
			"Update task with title exactly 255 characters",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Título original",
							Description: "Descrição original",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdate: func(ctx context.Context, taskUUID uuid.UUID, t *taskEntity.Task) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"title": strings.Repeat("a", 255),
			},
			&taskEntity.Task{
				UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:       strings.Repeat("a", 255),
				Description: "Descrição original",
				Status:      taskEntity.StatusTodo,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				taskRepo.SetPersist(originalPersist)
			}()
			if tt.setup != nil {
				tt.setup()
			}

			got, err := Update(tt.ctx, tt.taskUUID, tt.updates)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("Update() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Update() diff: %s", diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	originalPersist := taskRepo.Persist()

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		wantErr  error
	}{
		{
			"Delete task with success",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnDelete: func(ctx context.Context, taskUUID uuid.UUID) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			nil,
		},
		{
			"Delete task with context database error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnDelete: func(ctx context.Context, taskUUID uuid.UUID) error {
						return database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			database.ErrContextDatabase,
		},
		{
			"Delete task with generic persist error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnDelete: func(ctx context.Context, taskUUID uuid.UUID) error {
						return errors.New("database connection failed")
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			errors.New("database connection failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				taskRepo.SetPersist(originalPersist)
			}()
			if tt.setup != nil {
				tt.setup()
			}

			err := Delete(tt.ctx, tt.taskUUID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("Delete() error diff: %s", diff)
				return
			}
		})
	}
}

func TestListPaginated(t *testing.T) {
	originalPersist := taskRepo.Persist()

	statusTodo := taskEntity.StatusTodo
	statusInProgress := taskEntity.StatusInProgress
	statusDone := taskEntity.StatusDone
	statusCanceled := taskEntity.StatusCanceled
	invalidStatus := taskEntity.TaskStatus("invalid_status")

	tests := []struct {
		name         string
		setup        func()
		ctx          context.Context
		statusFilter *taskEntity.TaskStatus
		page         int
		limit        int
		want         *taskEntity.ListTasks
		wantErr      error
	}{
		{
			"ListPaginated all tasks with success",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return &taskEntity.ListTasks{
							Page:  1,
							Limit: 10,
							Tasks: []taskEntity.Task{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
									Title:       "Tarefa 1",
									Description: "Descrição 1",
									Status:      taskEntity.StatusTodo,
								},
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
									Title:       "Tarefa 2",
									Description: "Descrição 2",
									Status:      taskEntity.StatusInProgress,
								},
							},
							TotalItems: 2,
						}, nil
					},
				})
			},
			context.Background(),
			nil,
			1,
			10,
			&taskEntity.ListTasks{
				Page:  1,
				Limit: 10,
				Tasks: []taskEntity.Task{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Title:       "Tarefa 1",
						Description: "Descrição 1",
						Status:      taskEntity.StatusTodo,
					},
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
						Title:       "Tarefa 2",
						Description: "Descrição 2",
						Status:      taskEntity.StatusInProgress,
					},
				},
				TotalItems: 2,
			},
			nil,
		},
		{
			"ListPaginated filtered by status to_do with success",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return &taskEntity.ListTasks{
							Page:  1,
							Limit: 10,
							Tasks: []taskEntity.Task{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
									Title:       "Tarefa TODO",
									Description: "Descrição TODO",
									Status:      taskEntity.StatusTodo,
								},
							},
							TotalItems: 1,
						}, nil
					},
				})
			},
			context.Background(),
			&statusTodo,
			1,
			10,
			&taskEntity.ListTasks{
				Page:  1,
				Limit: 10,
				Tasks: []taskEntity.Task{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Title:       "Tarefa TODO",
						Description: "Descrição TODO",
						Status:      taskEntity.StatusTodo,
					},
				},
				TotalItems: 1,
			},
			nil,
		},
		{
			"ListPaginated filtered by status in_progress with success",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return &taskEntity.ListTasks{
							Page:  1,
							Limit: 10,
							Tasks: []taskEntity.Task{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
									Title:       "Tarefa In Progress",
									Description: "Descrição In Progress",
									Status:      taskEntity.StatusInProgress,
								},
							},
							TotalItems: 1,
						}, nil
					},
				})
			},
			context.Background(),
			&statusInProgress,
			1,
			10,
			&taskEntity.ListTasks{
				Page:  1,
				Limit: 10,
				Tasks: []taskEntity.Task{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
						Title:       "Tarefa In Progress",
						Description: "Descrição In Progress",
						Status:      taskEntity.StatusInProgress,
					},
				},
				TotalItems: 1,
			},
			nil,
		},
		{
			"ListPaginated filtered by status done with success",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return &taskEntity.ListTasks{
							Page:  1,
							Limit: 10,
							Tasks: []taskEntity.Task{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
									Title:       "Tarefa Done",
									Description: "Descrição Done",
									Status:      taskEntity.StatusDone,
								},
							},
							TotalItems: 1,
						}, nil
					},
				})
			},
			context.Background(),
			&statusDone,
			1,
			10,
			&taskEntity.ListTasks{
				Page:  1,
				Limit: 10,
				Tasks: []taskEntity.Task{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
						Title:       "Tarefa Done",
						Description: "Descrição Done",
						Status:      taskEntity.StatusDone,
					},
				},
				TotalItems: 1,
			},
			nil,
		},
		{
			"ListPaginated filtered by status canceled with success",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return &taskEntity.ListTasks{
							Page:  1,
							Limit: 10,
							Tasks: []taskEntity.Task{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"),
									Title:       "Tarefa Canceled",
									Description: "Descrição Canceled",
									Status:      taskEntity.StatusCanceled,
								},
							},
							TotalItems: 1,
						}, nil
					},
				})
			},
			context.Background(),
			&statusCanceled,
			1,
			10,
			&taskEntity.ListTasks{
				Page:  1,
				Limit: 10,
				Tasks: []taskEntity.Task{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"),
						Title:       "Tarefa Canceled",
						Description: "Descrição Canceled",
						Status:      taskEntity.StatusCanceled,
					},
				},
				TotalItems: 1,
			},
			nil,
		},
		{
			"ListPaginated with invalid status (validation moved to controller)",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						// Domain no longer validates status, it just passes it through
						return &taskEntity.ListTasks{
							Page:       1,
							Limit:      10,
							Tasks:      []taskEntity.Task{},
							TotalItems: 0,
						}, nil
					},
				})
			},
			context.Background(),
			&invalidStatus,
			1,
			10,
			&taskEntity.ListTasks{
				Page:       1,
				Limit:      10,
				Tasks:      []taskEntity.Task{},
				TotalItems: 0,
			},
			nil,
		},
		{
			"ListPaginated with context database error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return nil, database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			nil,
			1,
			10,
			nil,
			database.ErrContextDatabase,
		},
		{
			"ListPaginated with generic persist error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return nil, errors.New("database connection failed")
					},
				})
			},
			context.Background(),
			nil,
			1,
			10,
			nil,
			errors.New("database connection failed"),
		},
		{
			"ListPaginated with empty result",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return &taskEntity.ListTasks{
							Page:       1,
							Limit:      10,
							Tasks:      []taskEntity.Task{},
							TotalItems: 0,
						}, nil
					},
				})
			},
			context.Background(),
			nil,
			1,
			10,
			&taskEntity.ListTasks{
				Page:       1,
				Limit:      10,
				Tasks:      []taskEntity.Task{},
				TotalItems: 0,
			},
			nil,
		},
		{
			"ListPaginated with pagination - page 2, limit 2",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return &taskEntity.ListTasks{
							Page:  2,
							Limit: 2,
							Tasks: []taskEntity.Task{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
									Title:       "Tarefa 3",
									Description: "Descrição 3",
									Status:      taskEntity.StatusDone,
								},
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"),
									Title:       "Tarefa 4",
									Description: "Descrição 4",
									Status:      taskEntity.StatusCanceled,
								},
							},
							TotalItems: 4,
						}, nil
					},
				})
			},
			context.Background(),
			nil,
			2,
			2,
			&taskEntity.ListTasks{
				Page:  2,
				Limit: 2,
				Tasks: []taskEntity.Task{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
						Title:       "Tarefa 3",
						Description: "Descrição 3",
						Status:      taskEntity.StatusDone,
					},
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"),
						Title:       "Tarefa 4",
						Description: "Descrição 4",
						Status:      taskEntity.StatusCanceled,
					},
				},
				TotalItems: 4,
			},
			nil,
		},
		{
			"ListPaginated with limit 0 uses default limit",
			func() {
				Config.ListDefaultLimit = 15
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return &taskEntity.ListTasks{
							Page:  1,
							Limit: 15,
							Tasks: []taskEntity.Task{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
									Title:       "Tarefa 1",
									Description: "Descrição 1",
									Status:      taskEntity.StatusTodo,
								},
							},
							TotalItems: 1,
						}, nil
					},
				})
			},
			context.Background(),
			nil,
			1,
			0,
			&taskEntity.ListTasks{
				Page:  1,
				Limit: 15,
				Tasks: []taskEntity.Task{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Title:       "Tarefa 1",
						Description: "Descrição 1",
						Status:      taskEntity.StatusTodo,
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
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return &taskEntity.ListTasks{
							Page:  1,
							Limit: 15,
							Tasks: []taskEntity.Task{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
									Title:       "Tarefa 1",
									Description: "Descrição 1",
									Status:      taskEntity.StatusTodo,
								},
							},
							TotalItems: 1,
						}, nil
					},
				})
			},
			context.Background(),
			nil,
			1,
			-5,
			&taskEntity.ListTasks{
				Page:  1,
				Limit: 15,
				Tasks: []taskEntity.Task{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Title:       "Tarefa 1",
						Description: "Descrição 1",
						Status:      taskEntity.StatusTodo,
					},
				},
				TotalItems: 1,
			},
			nil,
		},
		{
			"ListPaginated with limit exceeding max uses max limit",
			func() {
				Config.ListDefaultLimit = 20
				Config.ListMaxLimit = 50
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnListPaginated: func(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
						return &taskEntity.ListTasks{
							Page:  1,
							Limit: 50,
							Tasks: []taskEntity.Task{
								{
									UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
									Title:       "Tarefa 1",
									Description: "Descrição 1",
									Status:      taskEntity.StatusTodo,
								},
							},
							TotalItems: 1,
						}, nil
					},
				})
			},
			context.Background(),
			nil,
			1,
			100,
			&taskEntity.ListTasks{
				Page:  1,
				Limit: 50,
				Tasks: []taskEntity.Task{
					{
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Title:       "Tarefa 1",
						Description: "Descrição 1",
						Status:      taskEntity.StatusTodo,
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
				taskRepo.SetPersist(originalPersist)
			}()

			if tt.setup != nil {
				tt.setup()
			}

			got, err := ListPaginated(tt.ctx, tt.statusFilter, tt.page, tt.limit)
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

func TestUpdateStatus(t *testing.T) {
	originalPersist := taskRepo.Persist()

	tests := []struct {
		name      string
		setup     func()
		ctx       context.Context
		taskUUID  uuid.UUID
		newStatus taskEntity.TaskStatus
		wantErr   error
	}{
		{
			"UpdateStatus with success - Todo to InProgress",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdateStatus: func(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusInProgress,
			nil,
		},
		{
			"UpdateStatus with success - Todo to Canceled",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdateStatus: func(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusCanceled,
			nil,
		},
		{
			"UpdateStatus with success - InProgress to Done",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusInProgress,
						}, nil
					},
					FnUpdateStatus: func(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusDone,
			nil,
		},
		{
			"UpdateStatus with success - InProgress to Canceled",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusInProgress,
						}, nil
					},
					FnUpdateStatus: func(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error {
						return nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusCanceled,
			nil,
		},
		{
			"UpdateStatus with invalid transition - Todo to Done",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusDone,
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"UpdateStatus with invalid transition - InProgress to Todo",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusInProgress,
						}, nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusTodo,
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"UpdateStatus with invalid transition - Done to InProgress",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusDone,
						}, nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusInProgress,
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"UpdateStatus with invalid transition - Canceled to Done",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusCanceled,
						}, nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusDone,
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"UpdateStatus with invalid status value",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.TaskStatus("invalid_status"),
			&errs.ValidationErrors{
				Errors: []errs.ValidationError{
					{
						Field:   "status",
						Message: "invalid status value",
					},
				},
			},
		},
		{
			"UpdateStatus task not found",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return nil, errs.ErrNotFound
					},
				})
			},
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			taskEntity.StatusInProgress,
			errs.ErrNotFound,
		},
		{
			"UpdateStatus with retrieve context database error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return nil, database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusInProgress,
			database.ErrContextDatabase,
		},
		{
			"UpdateStatus with update persist error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdateStatus: func(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error {
						return database.ErrContextDatabase
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusInProgress,
			database.ErrContextDatabase,
		},
		{
			"UpdateStatus with generic persist error",
			func() {
				taskRepo.SetPersist(&taskRepo.MockPersistent{
					FnRetrieveByUUID: func(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
						return &taskEntity.Task{
							UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
							Title:       "Tarefa",
							Description: "Descrição",
							Status:      taskEntity.StatusTodo,
						}, nil
					},
					FnUpdateStatus: func(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error {
						return errors.New("database connection failed")
					},
				})
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			taskEntity.StatusInProgress,
			errors.New("database connection failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				taskRepo.SetPersist(originalPersist)
			}()
			if tt.setup != nil {
				tt.setup()
			}

			err := UpdateStatus(tt.ctx, tt.taskUUID, tt.newStatus)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("UpdateStatus() error diff: %s", diff)
				return
			}
		})
	}
}
