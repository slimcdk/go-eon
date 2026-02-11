package eon

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// OAuth2TokenResponse represents the response from the token endpoint
type OAuth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// authenticate fetches an OAuth2 access token using client credentials flow
func (c *client) authenticate() error {
	var result OAuth2TokenResponse

	// Create a temporary client for token endpoint (no /api prefix)
	tokenClient := resty.New()

	// Request token using client credentials flow with form data
	res, err := tokenClient.R().
		SetFormData(map[string]string{
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
			"grant_type":    "client_credentials",
			"scope":         "navigator",
		}).
		SetResult(&result).
		Post(tokenEndpoint)

	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("authentication failed with status %d: %s", res.StatusCode(), res.String())
	}

	// Set access token and expiry on client
	c.accessToken = result.AccessToken
	// Subtract 60 seconds as safety margin
	c.tokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn-60) * time.Second)

	return nil
}

// GetAccessToken returns a valid access token, authenticating if necessary
func (c *client) GetAccessToken() (string, error) {
	if c.accessToken == "" || time.Now().After(c.tokenExpiry) {
		if err := c.authenticate(); err != nil {
			return "", err
		}
	}
	return c.accessToken, nil
}

// IsAlive checks if the API is reachable (health check)
func (c *client) IsAlive() (bool, error) {
	// Note: Eon API may not have a dedicated health endpoint
	// Using the token endpoint as a simple connectivity check
	res, err := c.resty.R().Get("/")
	if err != nil {
		return false, err
	}
	// Any response (even 404) means the API is reachable
	return res.StatusCode() > 0, nil
}
