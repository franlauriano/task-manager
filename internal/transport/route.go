package transport

import (
	"net/http"

	"taskmanager/internal/platform/database"
	"taskmanager/internal/transport/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// Routes defines the routes for the application
func Routes(dbConnector database.Connector) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestLogger(middleware.NewJSONLogFormatter(nil)))
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	dbTx := middleware.DatabaseWithTransaction(dbConnector)
	dbNoTx := middleware.DatabaseWithoutTransaction(dbConnector)

	r.Route("/api", func(r chi.Router) {
		// Task routes
		r.With(middleware.RequireContentTypeJSON).Post("/tasks", dbTx(CreateTask))
		r.Get("/tasks/{uuid}", dbNoTx(RetrieveByUUID))
		r.With(middleware.RequireContentTypeJSON).Put("/tasks/{uuid}", dbTx(UpdateTask))
		r.With(middleware.RequireContentTypeJSON).Delete("/tasks/{uuid}", dbTx(DeleteTask))
		r.Get("/tasks", dbNoTx(ListTasks))
		r.With(middleware.RequireContentTypeJSON).Post("/tasks/{uuid}/status", dbTx(UpdateTaskStatus))

		// Team routes
		r.With(middleware.RequireContentTypeJSON).Post("/teams", dbTx(CreateTeam))
		r.Get("/teams", dbNoTx(ListTeams))
		r.Get("/teams/{uuid}", dbNoTx(RetrieveTeamByUUID))
		r.With(middleware.RequireContentTypeJSON).Post("/teams/{uuid}/tasks", dbTx(AssociateTaskToTeam))
		r.With(middleware.RequireContentTypeJSON).Delete("/teams/{uuid}/tasks/{task_uuid}", dbTx(DisassociateTaskFromTeam))
	})
	return r
}
