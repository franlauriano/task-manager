//go:build test

package task

import (
	"context"
	"testing"
	"time"

	"taskmanager/internal/entity/task"
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
		testenv.WithContainerDatabaseAlias(databaseAliasTaskTest, databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name    string
		setup   func()
		ctx     context.Context
		task    *task.Task
		wantErr error
	}{
		{
			"Create task with success",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			&task.Task{
				Title:       "Nova tarefa",
				Description: "Descrição da nova tarefa",
				Status:      task.StatusTodo,
			},
			nil,
		},
		{
			"Create task with context nil",
			func() {
				resetWithMinimalData()
			},
			nil,
			&task.Task{
				Title:       "Nova tarefa",
				Description: "Descrição da nova tarefa",
				Status:      task.StatusTodo,
			},
			database.ErrContextDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx, databaseAliasTaskTest)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{alias: databaseAliasTaskTest}
			err := p.Create(ctx, tt.task)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("datasource.Create() error diff: %s", diff)
				return
			}
		})
	}
}

func Test_datasource_RetrieveByUUID(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithContainerDatabaseAlias(databaseAliasTaskTest, databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		want     *task.Task
		wantErr  error
	}{
		{
			"Retrieve task by UUID with success",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			&task.Task{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
					UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
				},
				UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:       "Implementar autenticação",
				Description: "Criar sistema de autenticação JWT para a API",
				Status:      task.StatusTodo,
			},
			nil,
		},
		{
			"Retrieve task by UUID not found",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			nil,
			errs.ErrNotFound,
		},
		{
			"Retrieve task by UUID with context nil",
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
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx, databaseAliasTaskTest)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{alias: databaseAliasTaskTest}
			got, err := p.RetrieveByUUID(ctx, tt.taskUUID)
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

func Test_datasource_Update(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithContainerDatabaseAlias(databaseAliasTaskTest, databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		task     *task.Task
		wantErr  error
	}{
		{
			"Update task with success",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			&task.Task{
				Title:       "Título atualizado",
				Description: "Descrição atualizada",
				Status:      task.StatusInProgress,
			},
			nil,
		},
		{
			"Update task not found",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			&task.Task{
				Title:       "Título atualizado",
				Description: "Descrição atualizada",
				Status:      task.StatusInProgress,
			},
			errs.ErrNotFound,
		},
		{
			"Update task with context nil",
			nil,
			nil,
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			&task.Task{
				Title:       "Título atualizado",
				Description: "Descrição atualizada",
				Status:      task.StatusInProgress,
			},
			database.ErrContextDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx, databaseAliasTaskTest)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{alias: databaseAliasTaskTest}
			err := p.Update(ctx, tt.taskUUID, tt.task)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("datasource.Update() error diff: %s", diff)
				return
			}
		})
	}
}

func Test_datasource_Delete(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithContainerDatabaseAlias(databaseAliasTaskTest, databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

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
				resetWithMinimalData()
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			nil,
		},
		{
			"Delete task not found",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			errs.ErrNotFound,
		},
		{
			"Delete task with context nil",
			nil,
			nil,
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			database.ErrContextDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx, databaseAliasTaskTest)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{alias: databaseAliasTaskTest}
			err := p.Delete(ctx, tt.taskUUID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("datasource.Delete() error diff: %s", diff)
				return
			}
		})
	}
}

func Test_datasource_ListPaginated(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithContainerDatabaseAlias(databaseAliasTaskTest, databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	statusTodo := task.StatusTodo
	statusInProgress := task.StatusInProgress
	statusDone := task.StatusDone
	statusCanceled := task.StatusCanceled

	tests := []struct {
		name         string
		setup        func()
		ctx          context.Context
		statusFilter *task.TaskStatus
		page         int
		limit        int
		want         *task.ListTasks
		wantErr      error
	}{
		{
			"ListPaginated all tasks - page 1, limit 3",
			resetWithMinimalData,
			context.Background(),
			nil,
			1,
			3,
			&task.ListTasks{
				Page:  1,
				Limit: 3,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        6,
							CreatedAt: time.Date(2025, 12, 1, 19, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 19, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
						Title:       "Implementar feature de notificações",
						Description: "Criar sistema de notificações em tempo real",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        12,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("423e4567-e89b-12d3-a456-426614174000"),
						Title:       "Criar testes de integração",
						Description: "Desenvolver suite completa de testes de integração",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(3); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        10,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("323e4567-e89b-12d3-a456-426614174000"),
						Title:       "Configurar monitoramento de logs",
						Description: "Implementar sistema centralizado de logs com ELK Stack",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(2); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
				},
				TotalItems: 14,
			},
			nil,
		},
		{
			"ListPaginated all tasks - page 2, limit 3",
			resetWithMinimalData,
			context.Background(),
			nil,
			2,
			3,
			&task.ListTasks{
				Page:  2,
				Limit: 3,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        5,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
						Title:       "Refatorar módulo de autenticação",
						Description: "Melhorar estrutura e organização do código de autenticação",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Title:       "Implementar autenticação",
						Description: "Criar sistema de autenticação JWT para a API",
						Status:      task.StatusTodo,
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        3,
							CreatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174004"),
						Title:       "Adicionar testes unitários",
						Description: "Escrever testes unitários para todas as funções principais",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
				},
				TotalItems: 14,
			},
			nil,
		},
		{
			"ListPaginated all tasks - page 3, limit 3 (last page)",
			resetWithMinimalData,
			context.Background(),
			nil,
			3,
			3,
			&task.ListTasks{
				Page:  3,
				Limit: 3,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        13,
							CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("423e4567-e89b-12d3-a456-426614174001"),
						Title:       "Executar testes de carga",
						Description: "Realizar testes de performance e carga na aplicação",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(3); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        11,
							CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("323e4567-e89b-12d3-a456-426614174001"),
						Title:       "Otimizar configuração do Docker",
						Description: "Melhorar Dockerfile e docker-compose para produção",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(2); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        4,
							CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174005"),
						Title:       "Otimizar queries do banco",
						Description: "Analisar e otimizar queries lentas do banco de dados",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
				},
				TotalItems: 14,
			},
			nil,
		},
		{
			"ListPaginated all tasks - page 1, limit 10 (all items)",
			resetWithMinimalData,
			context.Background(),
			nil,
			1,
			10,
			&task.ListTasks{
				Page:  1,
				Limit: 10,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        6,
							CreatedAt: time.Date(2025, 12, 1, 19, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 19, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
						Title:       "Implementar feature de notificações",
						Description: "Criar sistema de notificações em tempo real",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        12,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("423e4567-e89b-12d3-a456-426614174000"),
						Title:       "Criar testes de integração",
						Description: "Desenvolver suite completa de testes de integração",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(3); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        10,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("323e4567-e89b-12d3-a456-426614174000"),
						Title:       "Configurar monitoramento de logs",
						Description: "Implementar sistema centralizado de logs com ELK Stack",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(2); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        5,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
						Title:       "Refatorar módulo de autenticação",
						Description: "Melhorar estrutura e organização do código de autenticação",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Title:       "Implementar autenticação",
						Description: "Criar sistema de autenticação JWT para a API",
						Status:      task.StatusTodo,
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        3,
							CreatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174004"),
						Title:       "Adicionar testes unitários",
						Description: "Escrever testes unitários para todas as funções principais",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        13,
							CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("423e4567-e89b-12d3-a456-426614174001"),
						Title:       "Executar testes de carga",
						Description: "Realizar testes de performance e carga na aplicação",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(3); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        11,
							CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("323e4567-e89b-12d3-a456-426614174001"),
						Title:       "Otimizar configuração do Docker",
						Description: "Melhorar Dockerfile e docker-compose para produção",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(2); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        4,
							CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174005"),
						Title:       "Otimizar queries do banco",
						Description: "Analisar e otimizar queries lentas do banco de dados",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        2,
							CreatedAt: time.Date(2025, 11, 28, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
						Title:       "Criar documentação da API",
						Description: "Documentar todos os endpoints da API usando Swagger",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
				},
				TotalItems: 14,
			},
			nil,
		},
		{
			"ListPaginated filtered by status to_do - page 1, limit 10",
			resetWithMinimalData,
			context.Background(),
			&statusTodo,
			1,
			10,
			&task.ListTasks{
				Page:  1,
				Limit: 10,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        12,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("423e4567-e89b-12d3-a456-426614174000"),
						Title:       "Criar testes de integração",
						Description: "Desenvolver suite completa de testes de integração",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(3); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        10,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("323e4567-e89b-12d3-a456-426614174000"),
						Title:       "Configurar monitoramento de logs",
						Description: "Implementar sistema centralizado de logs com ELK Stack",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(2); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        5,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
						Title:       "Refatorar módulo de autenticação",
						Description: "Melhorar estrutura e organização do código de autenticação",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Title:       "Implementar autenticação",
						Description: "Criar sistema de autenticação JWT para a API",
						Status:      task.StatusTodo,
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        3,
							CreatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174004"),
						Title:       "Adicionar testes unitários",
						Description: "Escrever testes unitários para todas as funções principais",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
				},
				TotalItems: 5,
			},
			nil,
		},
		{
			"ListPaginated filtered by status in_progress - page 1, limit 10",
			resetWithMinimalData,
			context.Background(),
			&statusInProgress,
			1,
			10,
			&task.ListTasks{
				Page:  1,
				Limit: 10,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        6,
							CreatedAt: time.Date(2025, 12, 1, 19, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 19, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
						Title:       "Implementar feature de notificações",
						Description: "Criar sistema de notificações em tempo real",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        13,
							CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("423e4567-e89b-12d3-a456-426614174001"),
						Title:       "Executar testes de carga",
						Description: "Realizar testes de performance e carga na aplicação",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(3); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        11,
							CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("323e4567-e89b-12d3-a456-426614174001"),
						Title:       "Otimizar configuração do Docker",
						Description: "Melhorar Dockerfile e docker-compose para produção",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(2); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        4,
							CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174005"),
						Title:       "Otimizar queries do banco",
						Description: "Analisar e otimizar queries lentas do banco de dados",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
					{
						Model: gorm.Model{
							ID:        2,
							CreatedAt: time.Date(2025, 11, 28, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
						Title:       "Criar documentação da API",
						Description: "Documentar todos os endpoints da API usando Swagger",
						Status:      task.StatusInProgress,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  nil,
					},
				},
				TotalItems: 5,
			},
			nil,
		},
		{
			"ListPaginated filtered by status done - page 1, limit 10",
			resetWithMinimalData,
			context.Background(),
			&statusDone,
			1,
			10,
			&task.ListTasks{
				Page:  1,
				Limit: 10,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        14,
							CreatedAt: time.Date(2025, 11, 26, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("423e4567-e89b-12d3-a456-426614174002"),
						Title:       "Revisar cobertura de testes",
						Description: "Auditar e melhorar cobertura de testes do projeto",
						Status:      task.StatusDone,
						TeamID:      func() *uint { id := uint(3); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 28, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
					},
					{
						Model: gorm.Model{
							ID:        7,
							CreatedAt: time.Date(2025, 11, 24, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
						Title:       "Configurar CI/CD",
						Description: "Configurar pipeline de CI/CD usando GitHub Actions",
						Status:      task.StatusDone,
						TeamID:      func() *uint { id := uint(2); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 26, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
					},
					{
						Model: gorm.Model{
							ID:        9,
							CreatedAt: time.Date(2025, 11, 19, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174006"),
						Title:       "Criar dashboard de métricas",
						Description: "Implementar dashboard para visualizar métricas da aplicação",
						Status:      task.StatusDone,
						TeamID:      func() *uint { id := uint(2); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 21, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  func() *time.Time { t := time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC); return &t }(),
					},
				},
				TotalItems: 3,
			},
			nil,
		},
		{
			"ListPaginated filtered by status canceled - page 1, limit 10",
			resetWithMinimalData,
			context.Background(),
			&statusCanceled,
			1,
			10,
			&task.ListTasks{
				Page:  1,
				Limit: 10,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        8,
							CreatedAt: time.Date(2025, 11, 25, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 11, 28, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"),
						Title:       "Implementar cache Redis",
						Description: "Adicionar cache Redis para melhorar performance",
						Status:      task.StatusCanceled,
						TeamID:      func() *uint { id := uint(2); return &id }(),
						StartedAt:   func() *time.Time { t := time.Date(2025, 11, 27, 18, 21, 6, 0, time.UTC); return &t }(),
						FinishedAt:  func() *time.Time { t := time.Date(2025, 11, 28, 18, 21, 6, 0, time.UTC); return &t }(),
					},
				},
				TotalItems: 1,
			},
			nil,
		},
		{
			"ListPaginated filtered by status to_do - page 1, limit 1",
			resetWithMinimalData,
			context.Background(),
			&statusTodo,
			1,
			1,
			&task.ListTasks{
				Page:  1,
				Limit: 1,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        12,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("423e4567-e89b-12d3-a456-426614174000"),
						Title:       "Criar testes de integração",
						Description: "Desenvolver suite completa de testes de integração",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(3); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
				},
				TotalItems: 5,
			},
			nil,
		},
		{
			"ListPaginated filtered by status to_do - page 2, limit 1",
			resetWithMinimalData,
			context.Background(),
			&statusTodo,
			2,
			1,
			&task.ListTasks{
				Page:  2,
				Limit: 1,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        10,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("323e4567-e89b-12d3-a456-426614174000"),
						Title:       "Configurar monitoramento de logs",
						Description: "Implementar sistema centralizado de logs com ELK Stack",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(2); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
				},
				TotalItems: 5,
			},
			nil,
		},
		{
			"ListPaginated filtered by status to_do - page 3, limit 1",
			resetWithMinimalData,
			context.Background(),
			&statusTodo,
			3,
			1,
			&task.ListTasks{
				Page:  3,
				Limit: 1,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        5,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
						Title:       "Refatorar módulo de autenticação",
						Description: "Melhorar estrutura e organização do código de autenticação",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
				},
				TotalItems: 5,
			},
			nil,
		},
		{
			"ListPaginated filtered by status to_do - page 4, limit 1",
			resetWithMinimalData,
			context.Background(),
			&statusTodo,
			4,
			1,
			&task.ListTasks{
				Page:  4,
				Limit: 1,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						Title:       "Implementar autenticação",
						Description: "Criar sistema de autenticação JWT para a API",
						Status:      task.StatusTodo,
						StartedAt:   nil,
						FinishedAt:  nil,
					},
				},
				TotalItems: 5,
			},
			nil,
		},
		{
			"ListPaginated filtered by status to_do - page 5, limit 1 (empty page)",
			resetWithMinimalData,
			context.Background(),
			&statusTodo,
			5,
			1,
			&task.ListTasks{
				Page:  5,
				Limit: 1,
				Tasks: []task.Task{
					{
						Model: gorm.Model{
							ID:        3,
							CreatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
							UpdatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
						},
						UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174004"),
						Title:       "Adicionar testes unitários",
						Description: "Escrever testes unitários para todas as funções principais",
						Status:      task.StatusTodo,
						TeamID:      func() *uint { id := uint(1); return &id }(),
						StartedAt:   nil,
						FinishedAt:  nil,
					},
				},
				TotalItems: 5,
			},
			nil,
		},
		{
			"ListPaginated filtered by status to_do - page 6, limit 1 (empty page)",
			resetWithMinimalData,
			context.Background(),
			&statusTodo,
			6,
			1,
			&task.ListTasks{
				Page:       6,
				Limit:      1,
				Tasks:      []task.Task{},
				TotalItems: 5,
			},
			nil,
		},
		{
			"ListPaginated with context nil",
			nil,
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
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx, databaseAliasTaskTest)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{alias: databaseAliasTaskTest}
			got, err := p.ListPaginated(ctx, tt.statusFilter, tt.page, tt.limit)
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

func Test_datasource_UpdateStatus(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithContainerDatabaseAlias(databaseAliasTaskTest, databaseTest),
	)

	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		updates  map[string]any
		wantErr  error
	}{
		{
			"UpdateStatus with success",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"status": task.StatusInProgress,
			},
			nil,
		},
		{
			"UpdateStatus with multiple fields",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			map[string]any{
				"status":      task.StatusDone,
				"started_at":  time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
				"finished_at": time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
			},
			nil,
		},
		{
			"UpdateStatus task not found",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			map[string]any{
				"status": task.StatusInProgress,
			},
			errs.ErrNotFound,
		},
		{
			"UpdateStatus with context nil",
			nil,
			nil,
			uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			map[string]any{
				"status": task.StatusInProgress,
			},
			database.ErrContextDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx, databaseAliasTaskTest)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{alias: databaseAliasTaskTest}
			err := p.UpdateStatus(ctx, tt.taskUUID, tt.updates)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("datasource.UpdateStatus() error diff: %s", diff)
				return
			}
		})
	}
}

func Test_datasource_ListByTeamID(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithContainerDatabaseAlias(databaseAliasTaskTest, databaseTest),
	)
	resetWithMinimalData := func() {
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	teamID1 := uint(1)
	teamID2 := uint(2)
	teamID3 := uint(999)

	tests := []struct {
		name    string
		setup   func()
		ctx     context.Context
		teamID  uint
		want    []task.Task
		wantErr error
	}{
		{
			"ListByTeamID with tasks associated to team",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			teamID1,
			[]task.Task{
				{
					Model: gorm.Model{
						ID:        6,
						CreatedAt: time.Date(2025, 12, 1, 19, 21, 6, 0, time.UTC),
						UpdatedAt: time.Date(2025, 12, 1, 19, 21, 6, 0, time.UTC),
					},
					UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
					Title:       "Implementar feature de notificações",
					Description: "Criar sistema de notificações em tempo real",
					Status:      task.StatusInProgress,
					TeamID:      &teamID1,
					StartedAt:   nil,
					FinishedAt:  nil,
				},
				{
					Model: gorm.Model{
						ID:        5,
						CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
					},
					UUID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
					Title:       "Refatorar módulo de autenticação",
					Description: "Melhorar estrutura e organização do código de autenticação",
					Status:      task.StatusTodo,
					TeamID:      &teamID1,
					StartedAt:   nil,
					FinishedAt:  nil,
				},
				{
					Model: gorm.Model{
						ID:        3,
						CreatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
						UpdatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
					},
					UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174004"),
					Title:       "Adicionar testes unitários",
					Description: "Escrever testes unitários para todas as funções principais",
					Status:      task.StatusTodo,
					TeamID:      &teamID1,
					StartedAt:   nil,
					FinishedAt:  nil,
				},
				{
					Model: gorm.Model{
						ID:        4,
						CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
						UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
					},
					UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174005"),
					Title:       "Otimizar queries do banco",
					Description: "Analisar e otimizar queries lentas do banco de dados",
					Status:      task.StatusInProgress,
					TeamID:      &teamID1,
					StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
					FinishedAt:  nil,
				},
				{
					Model: gorm.Model{
						ID:        2,
						CreatedAt: time.Date(2025, 11, 28, 18, 21, 6, 0, time.UTC),
						UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
					},
					UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
					Title:       "Criar documentação da API",
					Description: "Documentar todos os endpoints da API usando Swagger",
					Status:      task.StatusInProgress,
					TeamID:      &teamID1,
					StartedAt:   func() *time.Time { t := time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC); return &t }(),
					FinishedAt:  nil,
				},
			},
			nil,
		},
		{
			"ListByTeamID with no tasks associated to team",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			teamID3,
			[]task.Task{},
			nil,
		},
		{
			"ListByTeamID filters correctly by team_id",
			func() {
				resetWithMinimalData()
			},
			context.Background(),
			teamID2,
			[]task.Task{
				{
					Model: gorm.Model{
						ID:        10,
						CreatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
						UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
					},
					UUID:        uuid.MustParse("323e4567-e89b-12d3-a456-426614174000"),
					Title:       "Configurar monitoramento de logs",
					Description: "Implementar sistema centralizado de logs com ELK Stack",
					Status:      task.StatusTodo,
					TeamID:      &teamID2,
					StartedAt:   nil,
					FinishedAt:  nil,
				},
				{
					Model: gorm.Model{
						ID:        11,
						CreatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
						UpdatedAt: time.Date(2025, 12, 1, 18, 21, 6, 0, time.UTC),
					},
					UUID:        uuid.MustParse("323e4567-e89b-12d3-a456-426614174001"),
					Title:       "Otimizar configuração do Docker",
					Description: "Melhorar Dockerfile e docker-compose para produção",
					Status:      task.StatusInProgress,
					TeamID:      &teamID2,
					StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
					FinishedAt:  nil,
				},
				{
					Model: gorm.Model{
						ID:        8,
						CreatedAt: time.Date(2025, 11, 25, 18, 21, 6, 0, time.UTC),
						UpdatedAt: time.Date(2025, 11, 28, 18, 21, 6, 0, time.UTC),
					},
					UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"),
					Title:       "Implementar cache Redis",
					Description: "Adicionar cache Redis para melhorar performance",
					Status:      task.StatusCanceled,
					TeamID:      &teamID2,
					StartedAt:   func() *time.Time { t := time.Date(2025, 11, 27, 18, 21, 6, 0, time.UTC); return &t }(),
					FinishedAt:  func() *time.Time { t := time.Date(2025, 11, 28, 18, 21, 6, 0, time.UTC); return &t }(),
				},
				{
					Model: gorm.Model{
						ID:        7,
						CreatedAt: time.Date(2025, 11, 24, 18, 21, 6, 0, time.UTC),
						UpdatedAt: time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC),
					},
					UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
					Title:       "Configurar CI/CD",
					Description: "Configurar pipeline de CI/CD usando GitHub Actions",
					Status:      task.StatusDone,
					TeamID:      &teamID2,
					StartedAt:   func() *time.Time { t := time.Date(2025, 11, 26, 18, 21, 6, 0, time.UTC); return &t }(),
					FinishedAt:  func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
				},
				{
					Model: gorm.Model{
						ID:        9,
						CreatedAt: time.Date(2025, 11, 19, 18, 21, 6, 0, time.UTC),
						UpdatedAt: time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC),
					},
					UUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174006"),
					Title:       "Criar dashboard de métricas",
					Description: "Implementar dashboard para visualizar métricas da aplicação",
					Status:      task.StatusDone,
					TeamID:      &teamID2,
					StartedAt:   func() *time.Time { t := time.Date(2025, 11, 21, 18, 21, 6, 0, time.UTC); return &t }(),
					FinishedAt:  func() *time.Time { t := time.Date(2025, 11, 29, 18, 21, 6, 0, time.UTC); return &t }(),
				},
			},
			nil,
		},
		{
			"ListByTeamID with context nil",
			nil,
			nil,
			teamID1,
			nil,
			database.ErrContextDatabase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithTransaction(t, tt.ctx, databaseAliasTaskTest)
			}

			if tt.setup != nil {
				tt.setup()
			}

			p := &datasource{alias: databaseAliasTaskTest}
			got, err := p.ListByTeamID(ctx, tt.teamID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("datasource.ListByTeamID() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("datasource.ListByTeamID() diff: %s", diff)
			}
		})
	}
}
