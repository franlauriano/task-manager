//go:build test

package task

import (
	"context"
	"testing"
	"time"

	"taskmanager/internal/entity/task"
	"taskmanager/internal/paths"
	"taskmanager/internal/platform/cache"
	errs "taskmanager/internal/platform/errors"
	"taskmanager/internal/platform/testing/assert"
	"taskmanager/internal/platform/testing/dbtest"
	"taskmanager/internal/platform/testing/testenv"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Test_cachedDatasource_ListPaginated(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
		testenv.WithRedis(redisTest),
	)

	resetWithMinimalData := func() {
		env.FlushRedis()
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	statusTodo := task.StatusTodo

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
			"Cache miss - all tasks page 1 limit 3",
			resetWithMinimalData,
			context.Background(),
			nil, 1, 3,
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
					},
				},
				TotalItems: 14,
			},
			nil,
		},
		{
			"Cache miss - filtered by status to_do",
			resetWithMinimalData,
			context.Background(),
			&statusTodo, 1, 10,
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
					},
				},
				TotalItems: 5,
			},
			nil,
		},
		{
			"Cache miss - empty page beyond total",
			resetWithMinimalData,
			context.Background(),
			nil, 100, 10,
			&task.ListTasks{
				Page:       100,
				Limit:      10,
				Tasks:      []task.Task{},
				TotalItems: 14,
			},
			nil,
		},
		{
			"Cache hit - returns cached data instead of querying database",
			func() {
				env.FlushRedis()
				dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
				_ = cache.Set(context.Background(), env.Redis, listCacheKey(nil, 1, 10), &task.ListTasks{
					Page:       1,
					Limit:      10,
					TotalItems: 999,
				}, 5*time.Minute)
			},
			context.Background(),
			nil, 1, 10,
			&task.ListTasks{
				Page:       1,
				Limit:      10,
				TotalItems: 999,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithoutTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			cached := NewCachedPersist(&datasource{}, env.Redis, 5*time.Minute)
			got, err := cached.ListPaginated(ctx, tt.statusFilter, tt.page, tt.limit)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("cachedDatasource.ListPaginated() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("cachedDatasource.ListPaginated() diff: %s", diff)
			}
		})
	}
}

func Test_cachedDatasource_Create(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
		testenv.WithRedis(redisTest),
	)

	statusTodo := task.StatusTodo

	populateCacheAndReset := func() {
		env.FlushRedis()
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
		_ = cache.Set(context.Background(), env.Redis, listCacheKey(nil, 1, 10), &task.ListTasks{TotalItems: 14}, 5*time.Minute)
		_ = cache.Set(context.Background(), env.Redis, listCacheKey(&statusTodo, 1, 10), &task.ListTasks{TotalItems: 5}, 5*time.Minute)
	}

	tests := []struct {
		name    string
		setup   func()
		ctx     context.Context
		task    *task.Task
		wantErr error
	}{
		{
			"Create task with success and invalidate cache",
			populateCacheAndReset,
			context.Background(),
			&task.Task{
				Title:       "Nova tarefa via cache",
				Description: "Descrição da nova tarefa",
				Status:      task.StatusTodo,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithoutTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			cached := NewCachedPersist(&datasource{}, env.Redis, 5*time.Minute)
			err := cached.Create(ctx, tt.task)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("cachedDatasource.Create() error diff: %s", diff)
				return
			}

			if tt.wantErr == nil {
				keyAll := listCacheKey(nil, 1, 10)
				keyTodo := listCacheKey(&statusTodo, 1, 10)
				afterAll, _ := cache.Get[task.ListTasks](ctx, env.Redis, keyAll)
				afterTodo, _ := cache.Get[task.ListTasks](ctx, env.Redis, keyTodo)
				if afterAll != nil {
					t.Error("expected 'all' cache entry to be invalidated after Create")
				}
				if afterTodo != nil {
					t.Error("expected 'to_do' cache entry to be invalidated after Create")
				}
			}
		})
	}
}

func Test_cachedDatasource_Update(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
		testenv.WithRedis(redisTest),
	)

	populateCacheAndReset := func() {
		env.FlushRedis()
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
		_ = cache.Set(context.Background(), env.Redis, listCacheKey(nil, 1, 10), &task.ListTasks{TotalItems: 14}, 5*time.Minute)
	}

	existingTaskUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		task     *task.Task
		wantErr  error
	}{
		{
			"Update task with success and invalidate cache",
			populateCacheAndReset,
			context.Background(),
			existingTaskUUID,
			&task.Task{
				Title:       "Título atualizado",
				Description: "Descrição atualizada",
				Status:      task.StatusInProgress,
			},
			nil,
		},
		{
			"Update task not found preserves cache",
			populateCacheAndReset,
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			&task.Task{
				Title:       "Título atualizado",
				Description: "Descrição atualizada",
			},
			errs.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithoutTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			cached := NewCachedPersist(&datasource{}, env.Redis, 5*time.Minute)
			err := cached.Update(ctx, tt.taskUUID, tt.task)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("cachedDatasource.Update() error diff: %s", diff)
				return
			}

			key := listCacheKey(nil, 1, 10)
			after, _ := cache.Get[task.ListTasks](ctx, env.Redis, key)
			if tt.wantErr == nil && after != nil {
				t.Error("expected list cache to be invalidated after successful Update")
			}
			if tt.wantErr != nil && after == nil {
				t.Error("expected list cache to remain after failed Update")
			}
		})
	}
}

func Test_cachedDatasource_Delete(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
		testenv.WithRedis(redisTest),
	)

	populateCacheAndReset := func() {
		env.FlushRedis()
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
		_ = cache.Set(context.Background(), env.Redis, listCacheKey(nil, 1, 10), &task.ListTasks{TotalItems: 14}, 5*time.Minute)
	}

	existingTaskUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		wantErr  error
	}{
		{
			"Delete task with success and invalidate cache",
			populateCacheAndReset,
			context.Background(),
			existingTaskUUID,
			nil,
		},
		{
			"Delete task not found preserves cache",
			populateCacheAndReset,
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			errs.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithoutTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			cached := NewCachedPersist(&datasource{}, env.Redis, 5*time.Minute)
			err := cached.Delete(ctx, tt.taskUUID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("cachedDatasource.Delete() error diff: %s", diff)
				return
			}

			key := listCacheKey(nil, 1, 10)
			after, _ := cache.Get[task.ListTasks](ctx, env.Redis, key)
			if tt.wantErr == nil && after != nil {
				t.Error("expected list cache to be invalidated after successful Delete")
			}
			if tt.wantErr != nil && after == nil {
				t.Error("expected list cache to remain after failed Delete")
			}
		})
	}
}

func Test_cachedDatasource_UpdateStatus(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
		testenv.WithRedis(redisTest),
	)

	populateCacheAndReset := func() {
		env.FlushRedis()
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
		_ = cache.Set(context.Background(), env.Redis, listCacheKey(nil, 1, 10), &task.ListTasks{TotalItems: 14}, 5*time.Minute)
	}

	existingTaskUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name     string
		setup    func()
		ctx      context.Context
		taskUUID uuid.UUID
		updates  map[string]any
		wantErr  error
	}{
		{
			"UpdateStatus with success and invalidate cache",
			populateCacheAndReset,
			context.Background(),
			existingTaskUUID,
			map[string]any{
				"status": task.StatusInProgress,
			},
			nil,
		},
		{
			"UpdateStatus task not found preserves cache",
			populateCacheAndReset,
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			map[string]any{
				"status": task.StatusInProgress,
			},
			errs.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithoutTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			cached := NewCachedPersist(&datasource{}, env.Redis, 5*time.Minute)
			err := cached.UpdateStatus(ctx, tt.taskUUID, tt.updates)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("cachedDatasource.UpdateStatus() error diff: %s", diff)
				return
			}

			key := listCacheKey(nil, 1, 10)
			after, _ := cache.Get[task.ListTasks](ctx, env.Redis, key)
			if tt.wantErr == nil && after != nil {
				t.Error("expected list cache to be invalidated after successful UpdateStatus")
			}
			if tt.wantErr != nil && after == nil {
				t.Error("expected list cache to remain after failed UpdateStatus")
			}
		})
	}
}

func Test_cachedDatasource_RetrieveByUUID(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
		testenv.WithRedis(redisTest),
	)

	resetWithMinimalData := func() {
		env.FlushRedis()
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
			"RetrieveByUUID delegates to real datasource",
			resetWithMinimalData,
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
			"RetrieveByUUID not found",
			resetWithMinimalData,
			context.Background(),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			nil,
			errs.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithoutTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			cached := NewCachedPersist(&datasource{}, env.Redis, 5*time.Minute)
			got, err := cached.RetrieveByUUID(ctx, tt.taskUUID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("cachedDatasource.RetrieveByUUID() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("cachedDatasource.RetrieveByUUID() diff: %s", diff)
			}
		})
	}
}

func Test_cachedDatasource_ListByTeamID(t *testing.T) {
	env := testenv.Setup(t,
		testenv.WithDatabase(databaseTest),
		testenv.WithRedis(redisTest),
	)

	resetWithMinimalData := func() {
		env.FlushRedis()
		dbtest.ResetWithFixtures(env.DB, paths.FixtureDir(), "tasks_minimal.sql")
	}

	teamID3 := uint(3)

	tests := []struct {
		name    string
		setup   func()
		ctx     context.Context
		teamID  uint
		want    []task.Task
		wantErr error
	}{
		{
			"ListByTeamID delegates to real datasource",
			resetWithMinimalData,
			context.Background(),
			teamID3,
			[]task.Task{
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
					TeamID:      &teamID3,
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
					TeamID:      &teamID3,
					StartedAt:   func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
				},
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
					TeamID:      &teamID3,
					StartedAt:   func() *time.Time { t := time.Date(2025, 11, 28, 18, 21, 6, 0, time.UTC); return &t }(),
					FinishedAt:  func() *time.Time { t := time.Date(2025, 11, 30, 18, 21, 6, 0, time.UTC); return &t }(),
				},
			},
			nil,
		},
		{
			"ListByTeamID with no tasks",
			resetWithMinimalData,
			context.Background(),
			uint(999),
			[]task.Task{},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if tt.ctx != nil {
				ctx = dbtest.SetupDBWithoutTransaction(t, tt.ctx)
			}

			if tt.setup != nil {
				tt.setup()
			}

			cached := NewCachedPersist(&datasource{}, env.Redis, 5*time.Minute)
			got, err := cached.ListByTeamID(ctx, tt.teamID)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("cachedDatasource.ListByTeamID() error diff: %s", diff)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("cachedDatasource.ListByTeamID() diff: %s", diff)
			}
		})
	}
}

func Test_listCacheKey(t *testing.T) {
	statusTodo := task.StatusTodo

	tests := []struct {
		name         string
		statusFilter *task.TaskStatus
		page         int
		limit        int
		want         string
	}{
		{"without filter", nil, 1, 10, "tasks:list:status=all:page=1:limit=10"},
		{"with status filter", &statusTodo, 2, 20, "tasks:list:status=to_do:page=2:limit=20"},
		{"different page", nil, 3, 5, "tasks:list:status=all:page=3:limit=5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := listCacheKey(tt.statusFilter, tt.page, tt.limit)
			if got != tt.want {
				t.Errorf("listCacheKey() = %q, want %q", got, tt.want)
			}
		})
	}
}
