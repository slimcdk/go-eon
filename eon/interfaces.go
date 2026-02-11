package eon

import (
	"time"
)

// Client represents the Eon API client interface
type Client interface {
	// Authentication
	GetAccessToken() (string, error)

	// Installations
	GetInstallations(filter []string) (InstallationsWrapper, error)
	GetMeasurementSeries() (InstallationsMeasurementsWrapper, error)

	// Measurements
	GetMeasurements(id int, resolution Resolution, from, to time.Time, includeMissing bool) (MeasurementsWrapper, error)

	// Costs
	GetCosts(installationID string, from, to *time.Time) (interface{}, error)

	// Health check
	IsAlive() (bool, error)
}
