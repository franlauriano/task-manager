//go:build test

package task

import (
	"context"
	"log/slog"
	"taskmanager/internal/entity/task"

	"github.com/google/uuid"
)

// MockPersistent é um mock da interface Persistent para testes
type MockPersistent struct {
	FnCreate         func(context.Context, *task.Task) error
	FnRetrieveByUUID func(context.Context, uuid.UUID) (*task.Task, error)
	FnUpdate         func(context.Context, uuid.UUID, *task.Task) error
	FnDelete         func(context.Context, uuid.UUID) error
	FnListPaginated  func(context.Context, *task.TaskStatus, int, int) (*task.ListTasks, error)
	FnUpdateStatus   func(context.Context, uuid.UUID, map[string]any) error
	FnListByTeamID   func(context.Context, uint) ([]task.Task, error)
}

// Create implementa o método Create da interface Persistent
func (m *MockPersistent) Create(ctx context.Context, t *task.Task) error {
	if m.FnCreate == nil {
		slog.Error("fnCreate is nil")
		return nil
	}
	return m.FnCreate(ctx, t)
}

// RetrieveByUUID implementa o método RetrieveByUUID da interface Persistent
func (m *MockPersistent) RetrieveByUUID(ctx context.Context, taskUUID uuid.UUID) (*task.Task, error) {
	if m.FnRetrieveByUUID == nil {
		slog.Error("fnRetrieveByUUID is nil")
		return nil, nil
	}
	return m.FnRetrieveByUUID(ctx, taskUUID)
}

// Update implementa o método Update da interface Persistent
func (m *MockPersistent) Update(ctx context.Context, taskUUID uuid.UUID, t *task.Task) error {
	if m.FnUpdate == nil {
		slog.Error("fnUpdate is nil")
		return nil
	}
	return m.FnUpdate(ctx, taskUUID, t)
}

// Delete implementa o método Delete da interface Persistent
func (m *MockPersistent) Delete(ctx context.Context, taskUUID uuid.UUID) error {
	if m.FnDelete == nil {
		slog.Error("fnDelete is nil")
		return nil
	}
	return m.FnDelete(ctx, taskUUID)
}

// ListPaginated implementa o método ListPaginated da interface Persistent
func (m *MockPersistent) ListPaginated(ctx context.Context, statusFilter *task.TaskStatus, page, limit int) (*task.ListTasks, error) {
	if m.FnListPaginated == nil {
		slog.Error("fnListPaginated is nil")
		return nil, nil
	}
	return m.FnListPaginated(ctx, statusFilter, page, limit)
}

// UpdateStatus implementa o método UpdateStatus da interface Persistent
func (m *MockPersistent) UpdateStatus(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error {
	if m.FnUpdateStatus == nil {
		slog.Error("fnUpdateStatus is nil")
		return nil
	}
	return m.FnUpdateStatus(ctx, taskUUID, updates)
}

// ListByTeamID implementa o método ListByTeamID da interface Persistent
func (m *MockPersistent) ListByTeamID(ctx context.Context, teamID uint) ([]task.Task, error) {
	if m.FnListByTeamID == nil {
		slog.Error("fnListByTeamID is nil")
		return nil, nil
	}
	return m.FnListByTeamID(ctx, teamID)
}
