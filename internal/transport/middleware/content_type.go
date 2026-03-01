package middleware

import "net/http"

// RequireContentTypeJSON returns 415 when Content-Type is not application/json.
func RequireContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(`{"message":"Content-Type must be application/json"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
