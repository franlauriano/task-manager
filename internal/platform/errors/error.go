package errors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("not found")
)

// ValidationError represents a single validation error for a specific field.
type ValidationError struct {
	Field   string         `json:"field"`            // Name of the field that failed validation
	Type    string         `json:"type"`             // Error type (e.g., "required", "invalid", "format")
	Code    string         `json:"code"`             // Unique error code for programmatic identification
	Message string         `json:"message"`          // Human-readable error message (fallback in English or default language)
	Params  map[string]any `json:"params,omitempty"` // Additional error parameters (optional)
}

// ValidationErrors represents a collection of validation errors.
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements Go's error interface, returning a string
// with all error messages separated by semicolons.
func (v *ValidationErrors) Error() string {
	errors := []string{}
	for _, err := range v.Errors {
		errors = append(errors, err.Message)
	}
	return strings.Join(errors, ";")
}

// BadRequestError represents a bad request error response
type BadRequestError struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// Error implements Go's error interface, returning a string with the error message.
func (b *BadRequestError) Error() string {
	return fmt.Sprintf("field: %s, error: %s", b.Field, b.Message)
}
