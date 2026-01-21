package dto

import (
	"time"

	"github.com/google/uuid"

	"taskmanager/internal/entity/task"
)

// TaskResponse represents the API response for a task
type TaskResponse struct {
	UUID        uuid.UUID  `json:"uuid"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ToTaskResponse converts a task.Task to TaskResponse
func ToTaskResponse(t task.Task) TaskResponse {
	return TaskResponse{
		UUID:        t.UUID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		FinishedAt:  t.FinishedAt,
		StartedAt:   t.StartedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// TasksResponse represents a list of tasks
type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

// PaginatedTasksResponse represents a paginated list of tasks
type PaginatedTasksResponse struct {
	Page         int            `json:"page"`
	ItemsPerPage int            `json:"items_per_page"`
	TotalItems   int            `json:"total_items"`
	TotalPages   int            `json:"total_pages"`
	Items        []TaskResponse `json:"items"`
}

// ToPaginatedTasksResponse converts pagination info and tasks to PaginatedTasksResponse
func ToPaginatedTasksResponse(page, limit, totalItems int, tasks []task.Task) PaginatedTasksResponse {
	totalPages := (totalItems + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	data := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		data[i] = ToTaskResponse(t)
	}

	return PaginatedTasksResponse{
		Page:         page,
		ItemsPerPage: limit,
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		Items:        data,
	}
}
