package transport

import (
	"net/http"

	"taskmanager/internal/transport/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// Routes defines the routes for the application
func Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestLogger(middleware.NewJSONLogFormatter(nil)))
	r.Route("/api", func(r chi.Router) {
		// Task routes
		r.With(middleware.RequireContentTypeJSON).Post("/tasks", middleware.DatabaseWithTransaction(CreateTask))
		r.Get("/tasks/{uuid}", middleware.DatabaseWithoutTransaction(RetrieveByUUID))
		r.With(middleware.RequireContentTypeJSON).Put("/tasks/{uuid}", middleware.DatabaseWithTransaction(UpdateTask))
		r.With(middleware.RequireContentTypeJSON).Delete("/tasks/{uuid}", middleware.DatabaseWithTransaction(DeleteTask))
		r.Get("/tasks", middleware.DatabaseWithoutTransaction(ListTasks))
		r.With(middleware.RequireContentTypeJSON).Post("/tasks/{uuid}/status", middleware.DatabaseWithTransaction(UpdateTaskStatus))

		// Team routes
		r.With(middleware.RequireContentTypeJSON).Post("/teams", middleware.DatabaseWithTransaction(CreateTeam))
		r.Get("/teams", middleware.DatabaseWithoutTransaction(ListTeams))
		r.Get("/teams/{uuid}", middleware.DatabaseWithoutTransaction(RetrieveTeamByUUID))
		r.With(middleware.RequireContentTypeJSON).Post("/teams/{uuid}/tasks", middleware.DatabaseWithTransaction(AssociateTaskToTeam))
		r.With(middleware.RequireContentTypeJSON).Delete("/teams/{uuid}/tasks/{task_uuid}", middleware.DatabaseWithTransaction(DisassociateTaskFromTeam))
	})
	return r
}
