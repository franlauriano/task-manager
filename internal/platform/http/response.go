package http

import (
	"encoding/json"
	"errors"
	"net/http"

	appErrors "taskmanager/internal/platform/errors"
)

// writeResponse marshals a response object and returns HTTP status code and body
func writeResponse(statusCode int, v interface{}) (int, []byte) {
	body, err := json.Marshal(v)
	if err != nil {
		return http.StatusInternalServerError, []byte{}
	}
	return statusCode, body
}

// HandleErrorResponse handles errors and responses, returning appropriate HTTP status code and body
func HandleErrorResponse(err error, resp any) (int, []byte) {
	if err == nil {
		return writeResponse(http.StatusOK, resp)
	}
	if validErr, ok := err.(*appErrors.ValidationErrors); ok {
		return writeResponse(http.StatusUnprocessableEntity, validErr)
	}
	if errors.Is(err, appErrors.ErrNotFound) {
		return writeResponse(http.StatusNotFound, nil)
	}
	if badReqErr, ok := err.(*appErrors.BadRequestError); ok {
		return writeResponse(http.StatusBadRequest, badReqErr)
	}
	return writeResponse(http.StatusInternalServerError, nil)
}

// BadRequest returns a bad request error response with HTTP status code and body
func BadRequest(message, field string) (int, []byte) {
	resp := appErrors.BadRequestError{
		Message: message,
	}
	if field != "" {
		resp.Field = field
	}
	return writeResponse(http.StatusBadRequest, resp)
}
