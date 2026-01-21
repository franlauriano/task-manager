package dto

import (
	"time"

	"github.com/google/uuid"

	"taskmanager/internal/entity/team"
)

// TeamResponse represents the API response for a team
type TeamResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToTeamResponse converts a team.Team to TeamResponse
func ToTeamResponse(t team.Team) TeamResponse {
	return TeamResponse{
		UUID:        t.UUID,
		Name:        t.Name,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// TeamsResponse represents a list of teams
type TeamsResponse struct {
	Teams []TeamResponse `json:"teams"`
}

// PaginatedTeamsResponse represents a paginated list of teams
type PaginatedTeamsResponse struct {
	Page         int            `json:"page"`
	ItemsPerPage int            `json:"items_per_page"`
	TotalItems   int            `json:"total_items"`
	TotalPages   int            `json:"total_pages"`
	Items        []TeamResponse `json:"items"`
}

// ToPaginatedTeamsResponse converts pagination info and teams to PaginatedTeamsResponse
func ToPaginatedTeamsResponse(page, limit, totalItems int, teams []team.Team) PaginatedTeamsResponse {
	totalPages := (totalItems + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	data := make([]TeamResponse, len(teams))
	for i, t := range teams {
		data[i] = ToTeamResponse(t)
	}

	return PaginatedTeamsResponse{
		Page:         page,
		ItemsPerPage: limit,
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		Items:        data,
	}
}

// TeamWithTasksResponse represents a team with its associated tasks
type TeamWithTasksResponse struct {
	TeamResponse
	Tasks []TaskResponse `json:"tasks"`
}

// ToTeamWithTasksResponse converts a team and its tasks to TeamWithTasksResponse
func ToTeamWithTasksResponse(t team.Team) TeamWithTasksResponse {
	taskResponses := make([]TaskResponse, len(t.Tasks))
	for i, task := range t.Tasks {
		taskResponses[i] = ToTaskResponse(task)
	}

	return TeamWithTasksResponse{
		TeamResponse: ToTeamResponse(t),
		Tasks:        taskResponses,
	}
}
