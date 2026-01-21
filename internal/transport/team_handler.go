package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	httputil "taskmanager/internal/platform/http"
	"taskmanager/internal/transport/dto"
	"taskmanager/internal/usecase/team"
)

// CreateTeam creates a new team
func CreateTeam(w http.ResponseWriter, r *http.Request) (int, []byte) {
	var req dto.CreateTeamRequest
	if err := httputil.DecodeJSONBody(r, &req); err != nil {
		slog.Error("error decoding JSON body for create team", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	t := req.ToTeam()
	if err := team.Create(r.Context(), t); err != nil {
		slog.Error("error creating team", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return httputil.HandleErrorResponse(nil, dto.ToTeamResponse(*t))
}

// RetrieveTeamByUUID retrieves a team by UUID with its tasks
func RetrieveTeamByUUID(w http.ResponseWriter, r *http.Request) (int, []byte) {
	teamUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		slog.Error("error parsing UUID from path for retrieve team", "error", err)
		return httputil.BadRequest("invalid uuid format", "uuid")
	}

	t, err := team.RetrieveByUUIDWithTasks(r.Context(), teamUUID)
	if err != nil {
		slog.Error("error retrieving team with tasks", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return httputil.HandleErrorResponse(nil, dto.ToTeamWithTasksResponse(*t))
}

// ListTeams lists all teams with pagination
func ListTeams(w http.ResponseWriter, r *http.Request) (int, []byte) {
	pageParam := httputil.QueryParam(r, "page")
	page := 1
	if pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limitParam := r.URL.Query().Get("limit")
	limit := 0
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil {
			limit = parsedLimit
		}
	}

	result, err := team.ListPaginated(r.Context(), page, limit)
	if err != nil {
		slog.Error("error listing teams", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return httputil.HandleErrorResponse(nil, dto.ToPaginatedTeamsResponse(result.Page, result.Limit, result.TotalItems, result.Teams))
}

// AssociateTaskToTeam associates a task with a team
func AssociateTaskToTeam(w http.ResponseWriter, r *http.Request) (int, []byte) {
	teamUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		slog.Error("error parsing UUID from path for associate task to team", "error", err)
		return httputil.BadRequest("invalid uuid format", "uuid")
	}

	var req struct {
		TaskUUID string `json:"task_uuid"`
	}
	if err := httputil.DecodeJSONBody(r, &req); err != nil {
		slog.Error("error decoding JSON body for associate task to team", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	taskUUID, err := uuid.Parse(req.TaskUUID)
	if err != nil {
		slog.Error("error parsing task UUID for associate task to team", "error", err)
		return httputil.BadRequest("invalid task_uuid format", "task_uuid")
	}

	if err := team.AssociateTask(r.Context(), teamUUID, taskUUID); err != nil {
		slog.Error("error associating task to team", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return http.StatusOK, []byte{}
}

// DisassociateTaskFromTeam disassociates a task from a team
func DisassociateTaskFromTeam(w http.ResponseWriter, r *http.Request) (int, []byte) {
	teamUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		slog.Error("error parsing UUID from path for disassociate task from team", "error", err)
		return httputil.BadRequest("invalid uuid format", "uuid")
	}
	taskUUIDStr := chi.URLParam(r, "task_uuid")
	if taskUUIDStr == "" {
		return httputil.BadRequest("task_uuid is required", "task_uuid")
	}

	taskUUID, err := uuid.Parse(taskUUIDStr)
	if err != nil {
		slog.Error("error parsing task UUID for disassociate task from team", "error", err)
		return httputil.BadRequest("invalid task_uuid format", "task_uuid")
	}

	if err := team.DisassociateTask(r.Context(), teamUUID, taskUUID); err != nil {
		slog.Error("error disassociating task from team", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return http.StatusOK, []byte{}
}
