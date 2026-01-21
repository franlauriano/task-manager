package http

import (
	"encoding/json"
	"net/http"

	"taskmanager/internal/platform/errors"
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
	if validErr, ok := err.(*errors.ValidationErrors); ok {
		return writeResponse(http.StatusUnprocessableEntity, validErr)
	}
	if err == errors.ErrNotFound {
		return writeResponse(http.StatusNotFound, nil)
	}
	if badReqErr, ok := err.(*errors.BadRequestError); ok {
		return writeResponse(http.StatusBadRequest, badReqErr)
	}
	return writeResponse(http.StatusInternalServerError, nil)
}

// BadRequest returns a bad request error response with HTTP status code and body
func BadRequest(message, field string) (int, []byte) {
	resp := errors.BadRequestError{
		Message: message,
	}
	if field != "" {
		resp.Field = field
	}
	return writeResponse(http.StatusBadRequest, resp)
}
