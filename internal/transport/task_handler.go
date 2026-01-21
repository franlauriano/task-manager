package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	taskEntity "taskmanager/internal/entity/task"
	httputil "taskmanager/internal/platform/http"
	"taskmanager/internal/transport/dto"
	"taskmanager/internal/usecase/task"
)

// CreateTask creates a new task
func CreateTask(w http.ResponseWriter, r *http.Request) (int, []byte) {
	var req dto.CreateTaskRequest
	if err := httputil.DecodeJSONBody(r, &req); err != nil {
		slog.Error("error decoding JSON body for create task", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	t := req.ToTask()
	if err := task.Create(r.Context(), t); err != nil {
		slog.Error("error creating task", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return httputil.HandleErrorResponse(nil, dto.ToTaskResponse(*t))
}

// RetrieveByUUID retrieves a task by UUID
func RetrieveByUUID(w http.ResponseWriter, r *http.Request) (int, []byte) {
	taskUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		slog.Error("error parsing UUID from path for retrieve task", "error", err)
		return httputil.BadRequest("invalid uuid format", "uuid")
	}

	t, err := task.RetrieveByUUID(r.Context(), taskUUID)
	if err != nil {
		slog.Error("error retrieving task", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return httputil.HandleErrorResponse(nil, dto.ToTaskResponse(*t))
}

// UpdateTask updates an existing task
func UpdateTask(w http.ResponseWriter, r *http.Request) (int, []byte) {
	taskUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		slog.Error("error parsing UUID from path for update task", "error", err)
		return httputil.BadRequest("invalid uuid format", "uuid")
	}

	var req dto.UpdateTaskRequest
	if err := httputil.DecodeJSONBody(r, &req); err != nil {
		slog.Error("error decoding JSON body for update task", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	updates := map[string]any{
		"title":       req.Title,
		"description": req.Description,
	}

	t, err := task.Update(r.Context(), taskUUID, updates)
	if err != nil {
		slog.Error("error updating task", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return httputil.HandleErrorResponse(nil, dto.ToTaskResponse(*t))
}

// DeleteTask deletes a task (soft delete)
func DeleteTask(w http.ResponseWriter, r *http.Request) (int, []byte) {
	taskUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		slog.Error("error parsing UUID from path for delete task", "error", err)
		return httputil.BadRequest("invalid uuid format", "uuid")
	}

	if err := task.Delete(r.Context(), taskUUID); err != nil {
		slog.Error("error deleting task", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return http.StatusOK, []byte{}
}

// ListTasks lists all tasks with optional status filter and pagination
func ListTasks(w http.ResponseWriter, r *http.Request) (int, []byte) {
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

	statusFilter := httputil.QueryParam(r, "status")
	status, err := dto.ToTaskStatus(statusFilter)
	if err != nil {
		slog.Error("error listing tasks", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	result, err := task.ListPaginated(r.Context(), status, page, limit)
	if err != nil {
		slog.Error("error listing tasks", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return httputil.HandleErrorResponse(nil, dto.ToPaginatedTasksResponse(result.Page, result.Limit, result.TotalItems, result.Tasks))
}

// UpdateTaskStatus updates the status of a task
func UpdateTaskStatus(w http.ResponseWriter, r *http.Request) (int, []byte) {
	taskUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		slog.Error("error parsing UUID from path for update task status", "error", err)
		return httputil.BadRequest("invalid uuid format", "uuid")
	}

	var req dto.StatusUpdateRequest
	if err := httputil.DecodeJSONBody(r, &req); err != nil {
		slog.Error("error decoding JSON body for update task status", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	newStatus := taskEntity.TaskStatus(req.Status)

	if err := task.UpdateStatus(r.Context(), taskUUID, newStatus); err != nil {
		slog.Error("error updating task status", "error", err)
		return httputil.HandleErrorResponse(err, nil)
	}

	return http.StatusOK, []byte{}
}
