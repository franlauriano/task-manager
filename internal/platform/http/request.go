package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	apperrors "taskmanager/internal/platform/errors"
)

// DecodeJSONBody decodes JSON request body into the provided struct
// Returns BadRequestError for JSON syntax errors or type mismatches
// Returns other errors for internal server errors
func DecodeJSONBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err == nil {
		return nil
	}

	// Check for JSON syntax errors (malformed JSON)
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return &apperrors.BadRequestError{
			Message: "invalid JSON syntax",
			Field:   "",
		}
	}

	// Check for type mismatch errors (e.g., expecting int but got string)
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return &apperrors.BadRequestError{
			Message: fmt.Sprintf("invalid type for field %q: expected %s, got %s", typeErr.Field, typeErr.Type, typeErr.Value),
			Field:   typeErr.Field,
		}
	}

	return err
}

// QueryParam extracts a query parameter value by key from the request
func QueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
