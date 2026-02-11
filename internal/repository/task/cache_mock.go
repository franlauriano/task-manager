//go:build test

package task

import (
	"context"

	"taskmanager/internal/entity/task"

	"github.com/google/uuid"
)

// MockCachedPersistent decorates a Persistent with in-memory cache simulation.
// It mirrors the cache-aside pattern of cachedDatasource without requiring Redis.
type MockCachedPersistent struct {
	Next  Persistent
	store map[string]*task.ListTasks
}

// NewMockCachedPersist creates an in-memory cache-aside decorator for testing.
func NewMockCachedPersist(next Persistent) *MockCachedPersistent {
	return &MockCachedPersistent{
		Next:  next,
		store: make(map[string]*task.ListTasks),
	}
}

// ListPaginated checks the in-memory cache first; on miss, queries the next and caches the result.
func (m *MockCachedPersistent) ListPaginated(ctx context.Context, statusFilter *task.TaskStatus, page, limit int) (*task.ListTasks, error) {
	key := listCacheKey(statusFilter, page, limit)
	if cached, ok := m.store[key]; ok {
		return cached, nil
	}

	result, err := m.Next.ListPaginated(ctx, statusFilter, page, limit)
	if err != nil {
		return nil, err
	}

	m.store[key] = result
	return result, nil
}

// Create delegates to the next implementation and invalidates the list cache.
func (m *MockCachedPersistent) Create(ctx context.Context, t *task.Task) error {
	if err := m.Next.Create(ctx, t); err != nil {
		return err
	}
	m.invalidate()
	return nil
}

// Update delegates to the next implementation and invalidates the list cache.
func (m *MockCachedPersistent) Update(ctx context.Context, taskUUID uuid.UUID, t *task.Task) error {
	if err := m.Next.Update(ctx, taskUUID, t); err != nil {
		return err
	}
	m.invalidate()
	return nil
}

// Delete delegates to the next implementation and invalidates the list cache.
func (m *MockCachedPersistent) Delete(ctx context.Context, taskUUID uuid.UUID) error {
	if err := m.Next.Delete(ctx, taskUUID); err != nil {
		return err
	}
	m.invalidate()
	return nil
}

// UpdateStatus delegates to the next implementation and invalidates the list cache.
func (m *MockCachedPersistent) UpdateStatus(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error {
	if err := m.Next.UpdateStatus(ctx, taskUUID, updates); err != nil {
		return err
	}
	m.invalidate()
	return nil
}

// RetrieveByUUID delegates directly to the next implementation (no cache).
func (m *MockCachedPersistent) RetrieveByUUID(ctx context.Context, taskUUID uuid.UUID) (*task.Task, error) {
	return m.Next.RetrieveByUUID(ctx, taskUUID)
}

// ListByTeamID delegates directly to the next implementation (no cache).
func (m *MockCachedPersistent) ListByTeamID(ctx context.Context, teamID uint) ([]task.Task, error) {
	return m.Next.ListByTeamID(ctx, teamID)
}

// invalidate removes all cached list entries.
func (m *MockCachedPersistent) invalidate() {
	m.store = make(map[string]*task.ListTasks)
}
