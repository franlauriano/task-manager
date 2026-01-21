package middleware

import (
	"log/slog"
	"net/http"

	"taskmanager/internal/platform/database"
)

type handlerFunc func(http.ResponseWriter, *http.Request) (int, []byte)

// DatabaseWithTransaction inserts a database transaction into the request context.
// Execute rollback in case of error and commit in case of success
func DatabaseWithTransaction(next handlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := database.InjectDBsIntoContext(r.Context(), database.WithDBTransaction())
		if err != nil {
			slog.Error("Failed to inject database transaction", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal server error"}`))
			return
		}
		r = r.WithContext(ctx)

		statusCode, body := next(w, r)

		// Consider all success status codes (200-299) for commit
		if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
			if err := database.Commit(r.Context()); err != nil {
				slog.Error("Error on commit transaction", "error", err)
			}
		} else {
			if err := database.Rollback(r.Context()); err != nil {
				slog.Error("Error on rollback transaction", "error", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(body)
	})
}

// DatabaseWithoutTransaction inserts a database without transaction into the request context
func DatabaseWithoutTransaction(next handlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := database.InjectDBsIntoContext(r.Context(), database.WithDBWithoutTransaction())
		if err != nil {
			slog.Error("Failed to inject database without transaction", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal server error"}`))
			return
		}
		r = r.WithContext(ctx)

		statusCode, body := next(w, r)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(body)
	})
}
