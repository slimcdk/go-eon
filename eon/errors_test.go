package eon

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiError(t *testing.T) {
	t.Run("returns nil for 200 OK", func(t *testing.T) {
		err := apiError(http.StatusOK)
		assert.NoError(t, err)
	})

	t.Run("returns ErrorNoContent for 204", func(t *testing.T) {
		err := apiError(http.StatusNoContent)
		assert.Equal(t, ErrorNoContent, err)
	})

	t.Run("returns ErrorBadRequest for 400", func(t *testing.T) {
		err := apiError(http.StatusBadRequest)
		assert.Equal(t, ErrorBadRequest, err)
	})

	t.Run("returns ErrorUnauthorized for 401", func(t *testing.T) {
		err := apiError(http.StatusUnauthorized)
		assert.Equal(t, ErrorUnauthorized, err)
	})

	t.Run("returns ErrorNotFound for 404", func(t *testing.T) {
		err := apiError(http.StatusNotFound)
		assert.Equal(t, ErrorNotFound, err)
	})

	t.Run("returns ErrorTooManyRequests for 429", func(t *testing.T) {
		err := apiError(http.StatusTooManyRequests)
		assert.Equal(t, ErrorTooManyRequests, err)
	})

	t.Run("returns ErrorServerError for 500", func(t *testing.T) {
		err := apiError(http.StatusInternalServerError)
		assert.Equal(t, ErrorServerError, err)
	})

	t.Run("returns formatted error for unknown status code", func(t *testing.T) {
		err := apiError(418) // I'm a teapot
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 418")
	})
}
