package dto

import (
	"taskmanager/internal/entity/task"
	"taskmanager/internal/platform/errors"
)

// CreateTaskRequest represents the payload for creating a new task
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// ToTask converts CreateTaskRequest to task.Task
func (r *CreateTaskRequest) ToTask() *task.Task {
	return &task.Task{
		Title:       r.Title,
		Description: r.Description,
	}
}

// UpdateTaskRequest represents the payload for updating a task
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// ToTaskStatus converts a status filter string to *task.TaskStatus
// Returns nil if the string is empty
// Returns an error if the status is invalid
func ToTaskStatus(status string) (*task.TaskStatus, error) {
	if status == "" {
		return nil, nil
	}
	taskStatus := task.TaskStatus(status)
	if !isValidTaskStatus(taskStatus) {
		return nil, &errors.BadRequestError{
			Message: "invalid status value",
			Field:   "status",
		}
	}
	return &taskStatus, nil
}

// isValidTaskStatus validates if a TaskStatus is valid
// Returns true if valid, false otherwise
func isValidTaskStatus(taskStatus task.TaskStatus) bool {
	validStatuses := []task.TaskStatus{
		task.StatusTodo,
		task.StatusInProgress,
		task.StatusCanceled,
		task.StatusDone,
	}
	for _, vs := range validStatuses {
		if taskStatus == vs {
			return true
		}
	}
	return false
}
