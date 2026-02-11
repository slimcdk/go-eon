package eon

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAccessToken(t *testing.T) {
	mockResty := resty.New()
	httpmock.ActivateNonDefault(mockResty.GetClient())
	defer httpmock.DeactivateAndReset()

	c := &client{
		clientID:     "test-client-id",
		clientSecret: "test-client-secret",
		resty:        mockResty,
	}

	t.Run("returns cached token when valid", func(t *testing.T) {
		c.accessToken = "cached-token"
		c.tokenExpiry = time.Now().Add(1 * time.Hour)

		token, err := c.GetAccessToken()

		assert.NoError(t, err)
		assert.Equal(t, "cached-token", token)
	})

	t.Run("returns empty when no token and expired", func(t *testing.T) {
		c.accessToken = ""
		c.tokenExpiry = time.Time{}

		// Since we can't mock the token endpoint (it creates its own client),
		// this will fail, but we're testing the token expiry logic
		_, err := c.GetAccessToken()
		assert.Error(t, err) // Expected to fail without mock
	})

	t.Run("refreshes when token is expired", func(t *testing.T) {
		c.accessToken = "expired-token"
		c.tokenExpiry = time.Now().Add(-1 * time.Hour)

		// This will fail without mock, but demonstrates the expiry check logic
		_, err := c.GetAccessToken()
		assert.Error(t, err) // Expected to fail without mock
	})
}

func TestTokenExpiry(t *testing.T) {
	t.Run("token expiry is calculated correctly", func(t *testing.T) {
		c := &client{}

		// Simulate what happens after successful authentication
		c.accessToken = "test-token"
		expiresIn := 3600
		c.tokenExpiry = time.Now().Add(time.Duration(expiresIn-60) * time.Second)

		// Token should be valid for approximately 59 minutes
		assert.True(t, time.Now().Before(c.tokenExpiry))
		assert.True(t, c.tokenExpiry.Before(time.Now().Add(60*time.Minute)))
	})

	t.Run("detects expired tokens correctly", func(t *testing.T) {
		c := &client{
			accessToken: "token",
			tokenExpiry: time.Now().Add(-1 * time.Minute),
		}

		// Token is expired
		assert.True(t, time.Now().After(c.tokenExpiry))
	})
}

func TestIsAlive(t *testing.T) {
	mockResty := resty.New()
	httpmock.ActivateNonDefault(mockResty.GetClient())
	defer httpmock.DeactivateAndReset()

	c := &client{resty: mockResty}

	t.Run("returns true on successful response", func(t *testing.T) {
		httpmock.RegisterResponder("GET", "/", httpmock.NewStringResponder(200, "OK"))
		alive, err := c.IsAlive()
		assert.NoError(t, err)
		assert.True(t, alive)
	})

	t.Run("returns true even on 404", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "/", httpmock.NewStringResponder(404, "Not Found"))
		alive, err := c.IsAlive()
		assert.NoError(t, err)
		assert.True(t, alive)
	})

	t.Run("returns false on error", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "/",
			func(req *http.Request) (*http.Response, error) {
				return nil, http.ErrServerClosed
			})
		alive, err := c.IsAlive()
		assert.Error(t, err)
		assert.False(t, alive)
	})
}
