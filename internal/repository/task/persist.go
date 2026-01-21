package task

import (
	"context"
	"errors"

	"taskmanager/internal/entity/task"
	"taskmanager/internal/platform/database"
	errs "taskmanager/internal/platform/errors"
	"taskmanager/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Persistent defines the interface for task persistence
type Persistent interface {
	Create(ctx context.Context, t *task.Task) error
	RetrieveByUUID(ctx context.Context, taskUUID uuid.UUID) (*task.Task, error)
	Update(ctx context.Context, taskUUID uuid.UUID, t *task.Task) error
	Delete(ctx context.Context, taskUUID uuid.UUID) error
	ListPaginated(ctx context.Context, statusFilter *task.TaskStatus, page, limit int) (*task.ListTasks, error)
	UpdateStatus(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error
	ListByTeamID(ctx context.Context, teamID uint) ([]task.Task, error)
}

// datasource implements the persistent interface using PostgreSQL
type datasource struct {
	alias string
}

// persist is the global persistent implementation
var persist Persistent = &datasource{alias: repository.DatabaseAlias}

// SetPersist sets the persistent implementation
func SetPersist(p Persistent) {
	persist = p
}

// Persist returns the current persistent implementation
func Persist() Persistent {
	return persist
}

// Create saves a new task to the datasource
func (p *datasource) Create(ctx context.Context, t *task.Task) error {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return err
	}

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// RetrieveByUUID retrieves a task by UUID from the datasource
func (p *datasource) RetrieveByUUID(ctx context.Context, taskUUID uuid.UUID) (*task.Task, error) {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return nil, err
	}

	var t task.Task
	if err := db.Where("uuid = ?", taskUUID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &t, nil
}

// Update updates an existing task in the datasource
func (p *datasource) Update(ctx context.Context, taskUUID uuid.UUID, t *task.Task) error {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return err
	}

	result := db.Model(&task.Task{}).
		Where("uuid = ?", taskUUID).
		Updates(t)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// Delete performs a soft delete of a task in the datasource
func (p *datasource) Delete(ctx context.Context, taskUUID uuid.UUID) error {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return err
	}

	result := db.Where("uuid = ?", taskUUID).Delete(&task.Task{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// ListPaginated lists tasks with pagination and optional filters from the datasource
func (p *datasource) ListPaginated(ctx context.Context, statusFilter *task.TaskStatus, page, limit int) (*task.ListTasks, error) {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return nil, err
	}

	var tasks []task.Task
	var totalItems int64

	query := db.Model(&task.Task{})

	if statusFilter != nil {
		query = query.Where("status = ?", *statusFilter)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Order("id DESC").Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return &task.ListTasks{
		Limit:      limit,
		Page:       page,
		Tasks:      tasks,
		TotalItems: int(totalItems),
	}, nil
}

// UpdateStatus updates only the status and timestamps in the datasource
func (p *datasource) UpdateStatus(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) error {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return err
	}

	result := db.Model(&task.Task{}).
		Where("uuid = ?", taskUUID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// ListByTeamID lists all tasks associated with a team
func (p *datasource) ListByTeamID(ctx context.Context, teamID uint) ([]task.Task, error) {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return nil, err
	}

	var tasks []task.Task
	if err := db.Where("team_id = ?", teamID).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
