package eon

import (
	"fmt"
	"net/http"
)

// GetInstallations retrieves all installations with metadata.
// Optional filter can be provided to get specific installations.
//
// Example:
//
//	installations, err := client.GetInstallations([]string{"installation-id-1", "installation-id-2"})
func (c *client) GetInstallations(filter []string) (InstallationsWrapper, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return InstallationsWrapper{}, err
	}

	var result InstallationsWrapper

	req := c.resty.R().
		SetAuthToken(accessToken).
		SetResult(&result)

	// Add filter query parameters if provided
	if len(filter) > 0 {
		// Use SetQueryParamsFromValues for array parameters
		params := make(map[string][]string)
		params["installationFilter"] = filter
		req = req.SetQueryParamsFromValues(params)
	}

	res, err := req.Get("/installations")
	if err != nil {
		return InstallationsWrapper{}, err
	}

	if res.StatusCode() == http.StatusNoContent {
		return InstallationsWrapper{Installations: []InstallationDto{}}, nil
	}

	if res.StatusCode() != http.StatusOK {
		return InstallationsWrapper{}, fmt.Errorf("failed to get installations: %s (status %d)", res.String(), res.StatusCode())
	}

	return result, nil
}

// GetMeasurementSeries retrieves measurement series for all installations.
//
// Example:
//
//	series, err := client.GetMeasurementSeries()
func (c *client) GetMeasurementSeries() (InstallationsMeasurementsWrapper, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return InstallationsMeasurementsWrapper{}, err
	}

	var result InstallationsMeasurementsWrapper

	res, err := c.resty.R().
		SetAuthToken(accessToken).
		SetResult(&result).
		Get("/installations/measurement-series")

	if err != nil {
		return InstallationsMeasurementsWrapper{}, err
	}

	if res.StatusCode() == http.StatusNoContent {
		return InstallationsMeasurementsWrapper{Installations: []InstallationMeasurementsDto{}}, nil
	}

	if res.StatusCode() != http.StatusOK {
		return InstallationsMeasurementsWrapper{}, fmt.Errorf("failed to get measurement series: %s (status %d)", res.String(), res.StatusCode())
	}

	return result, nil
}
