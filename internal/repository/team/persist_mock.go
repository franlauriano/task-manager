//go:build test

package team

import (
	"context"
	"log/slog"
	"taskmanager/internal/entity/team"

	"github.com/google/uuid"
)

// MockPersistent é um mock da interface Persistent para testes
type MockPersistent struct {
	FnCreate             func(context.Context, *team.Team) error
	FnRetrieveByUUID     func(context.Context, uuid.UUID) (*team.Team, error)
	FnListPaginated      func(context.Context, int, int) (*team.ListTeams, error)
	FnRetrieveTaskTeamID func(context.Context, uuid.UUID) (*uint, error)
	FnUpdateTaskTeamID   func(context.Context, uuid.UUID, *uint) error
}

// Create implementa o método Create da interface Persistent
func (m *MockPersistent) Create(ctx context.Context, t *team.Team) error {
	if m.FnCreate == nil {
		slog.Error("fnCreate is nil")
		return nil
	}
	return m.FnCreate(ctx, t)
}

// RetrieveByUUID implementa o método RetrieveByUUID da interface Persistent
func (m *MockPersistent) RetrieveByUUID(ctx context.Context, teamUUID uuid.UUID) (*team.Team, error) {
	if m.FnRetrieveByUUID == nil {
		slog.Error("fnRetrieveByUUID is nil")
		return nil, nil
	}
	return m.FnRetrieveByUUID(ctx, teamUUID)
}

// ListPaginated implementa o método ListPaginated da interface Persistent
func (m *MockPersistent) ListPaginated(ctx context.Context, page, limit int) (*team.ListTeams, error) {
	if m.FnListPaginated == nil {
		slog.Error("fnListPaginated is nil")
		return nil, nil
	}
	return m.FnListPaginated(ctx, page, limit)
}

// RetrieveTaskTeamID implementa o método RetrieveTaskTeamID da interface Persistent
func (m *MockPersistent) RetrieveTaskTeamID(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
	if m.FnRetrieveTaskTeamID == nil {
		slog.Error("fnRetrieveTaskTeamID is nil")
		return nil, nil
	}
	return m.FnRetrieveTaskTeamID(ctx, taskUUID)
}

// UpdateTaskTeamID implementa o método UpdateTaskTeamID da interface Persistent
func (m *MockPersistent) UpdateTaskTeamID(ctx context.Context, taskUUID uuid.UUID, teamID *uint) error {
	if m.FnUpdateTaskTeamID == nil {
		slog.Error("fnUpdateTaskTeamID is nil")
		return nil
	}
	return m.FnUpdateTaskTeamID(ctx, taskUUID, teamID)
}
