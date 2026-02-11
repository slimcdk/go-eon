package eon

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("creates client from environment variables", func(t *testing.T) {
		// Set env vars
		_ = os.Setenv("CLIENT_ID", "test-client-id")
		_ = os.Setenv("CLIENT_SECRET", "test-client-secret")
		defer func() { _ = os.Unsetenv("CLIENT_ID") }()
		defer func() { _ = os.Unsetenv("CLIENT_SECRET") }()

		c := New()
		assert.NotNil(t, c)

		// Use type assertion to access the internal client struct for testing
		internalClient, ok := c.(*client)
		assert.True(t, ok, "client should be of internal type *client")

		assert.Equal(t, "test-client-id", internalClient.clientID)
		assert.Equal(t, "test-client-secret", internalClient.clientSecret)
		assert.Equal(t, apiBaseURL, internalClient.resty.BaseURL)
	})

	t.Run("creates client with explicit credentials", func(t *testing.T) {
		c := NewWithCredentials("my-client-id", "my-client-secret")
		assert.NotNil(t, c)

		// Use type assertion to access the internal client struct for testing
		internalClient, ok := c.(*client)
		assert.True(t, ok, "client should be of internal type *client")

		assert.Equal(t, "my-client-id", internalClient.clientID)
		assert.Equal(t, "my-client-secret", internalClient.clientSecret)
		assert.Equal(t, apiBaseURL, internalClient.resty.BaseURL)
	})
}
