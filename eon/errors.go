package eon

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// Common API errors
	ErrorUnauthorized    error = errors.New("unauthorized - check client credentials")
	ErrorNotFound        error = errors.New("resource not found")
	ErrorBadRequest      error = errors.New("bad request - check parameters")
	ErrorTooManyRequests error = errors.New("too many requests - rate limit exceeded")
	ErrorServerError     error = errors.New("server error")
	ErrorNoContent       error = errors.New("no content available")
)

// apiError maps HTTP status codes to errors
func apiError(statusCode int) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusNoContent:
		return ErrorNoContent
	case http.StatusBadRequest:
		return ErrorBadRequest
	case http.StatusUnauthorized:
		return ErrorUnauthorized
	case http.StatusNotFound:
		return ErrorNotFound
	case http.StatusTooManyRequests:
		return ErrorTooManyRequests
	case http.StatusInternalServerError:
		return ErrorServerError
	default:
		return fmt.Errorf("unexpected status code: %d", statusCode)
	}
}
