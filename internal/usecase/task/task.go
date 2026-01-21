package task

import (
	"context"
	"strings"
	taskEntity "taskmanager/internal/entity/task"
	taskRepo "taskmanager/internal/repository/task"
	"time"

	"github.com/google/uuid"
)

// Create creates a new task with business rules
func Create(ctx context.Context, t *taskEntity.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.Status = taskEntity.StatusTodo
	t.Title = strings.TrimSpace(t.Title)
	t.Description = strings.TrimSpace(t.Description)

	return taskRepo.Persist().Create(ctx, t)
}

// RetrieveByUUID retrieves a task by UUID
func RetrieveByUUID(ctx context.Context, taskUUID uuid.UUID) (*taskEntity.Task, error) {
	return taskRepo.Persist().RetrieveByUUID(ctx, taskUUID)
}

// Update updates an existing task
func Update(ctx context.Context, taskUUID uuid.UUID, updates map[string]any) (*taskEntity.Task, error) {
	t, err := taskRepo.Persist().RetrieveByUUID(ctx, taskUUID)
	if err != nil {
		return nil, err
	}

	if title, ok := updates["title"].(string); ok {
		t.Title = strings.TrimSpace(title)
	}
	if description, ok := updates["description"].(string); ok {
		t.Description = strings.TrimSpace(description)
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	err = taskRepo.Persist().Update(ctx, taskUUID, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Delete performs a soft delete of a task
func Delete(ctx context.Context, taskUUID uuid.UUID) error {
	return taskRepo.Persist().Delete(ctx, taskUUID)
}

// ListPaginated lists tasks with pagination and optional filters
func ListPaginated(ctx context.Context, statusFilter *taskEntity.TaskStatus, page, limit int) (*taskEntity.ListTasks, error) {
	if limit <= 0 {
		limit = Config.ListDefaultLimit
	}

	if limit > Config.ListMaxLimit {
		limit = Config.ListMaxLimit
	}

	return taskRepo.Persist().ListPaginated(ctx, statusFilter, page, limit)
}

// UpdateStatus updates the status of a task with transition validation
func UpdateStatus(ctx context.Context, taskUUID uuid.UUID, newStatus taskEntity.TaskStatus) error {
	task, err := taskRepo.Persist().RetrieveByUUID(ctx, taskUUID)
	if err != nil {
		return err
	}

	if err := task.Status.ValidateTransitionTo(newStatus); err != nil {
		return err
	}

	timestamp := time.Now()
	task.EnsureTimestampsForStatus(newStatus, &timestamp)

	updates := map[string]interface{}{
		"status":      newStatus,
		"started_at":  task.StartedAt,
		"finished_at": task.FinishedAt,
	}

	return taskRepo.Persist().UpdateStatus(ctx, taskUUID, updates)
}
