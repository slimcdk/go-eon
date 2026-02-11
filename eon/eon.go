package eon

import (
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// client is the internal implementation that satisfies the Client interface.
type client struct {
	clientID     string
	clientSecret string
	accessToken  string
	tokenExpiry  time.Time
	resty        *resty.Client
}

// New creates and returns a new Eon client.
// Credentials are loaded from environment variables CLIENT_ID and CLIENT_SECRET.
//
// Example:
//
//	client := eon.New()
func New() Client {
	return &client{
		clientID:     os.Getenv("CLIENT_ID"),
		clientSecret: os.Getenv("CLIENT_SECRET"),
		resty:        resty.New().SetBaseURL(apiBaseURL),
	}
}

// NewWithCredentials creates an Eon client with explicit credentials.
//
// Example:
//
//	client := eon.NewWithCredentials(clientID, clientSecret)
func NewWithCredentials(clientID, clientSecret string) Client {
	return &client{
		clientID:     clientID,
		clientSecret: clientSecret,
		resty:        resty.New().SetBaseURL(apiBaseURL),
	}
}
