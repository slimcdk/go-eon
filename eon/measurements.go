package eon

import (
	"fmt"
	"net/http"
	"time"
)

// GetMeasurements retrieves measurement data for a specific measurement series.
//
// Parameters:
//   - id: The measurement series ID (from GetMeasurementSeries)
//   - resolution: quarter, hour, day, or month
//   - from, to: Time range (required for quarter and hour resolutions)
//   - includeMissing: Whether to fill in missing values for the given resolution
//
// Example:
//
//	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//	to := time.Date(2024, 1, 31, 23, 59, 0, 0, time.UTC)
//	measurements, err := client.GetMeasurements(12345, eon.Hour, from, to, false)
func (c *client) GetMeasurements(id int, resolution Resolution, from, to time.Time, includeMissing bool) (MeasurementsWrapper, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return MeasurementsWrapper{}, err
	}

	var result MeasurementsWrapper

	path := fmt.Sprintf("/measurements/%d/resolution/%s", id, resolution)

	req := c.resty.R().
		SetAuthToken(accessToken).
		SetResult(&result)

	// Add time range parameters if provided
	if !from.IsZero() {
		// Format as RFC3339 with milliseconds and Z suffix
		req.SetQueryParam("from", from.Format("2006-01-02T15:04:05.000Z"))
	}
	if !to.IsZero() {
		req.SetQueryParam("to", to.Format("2006-01-02T15:04:05.000Z"))
	}

	// Add includeMissing parameter
	if includeMissing {
		req.SetQueryParam("includeMissing", "true")
	} else {
		req.SetQueryParam("includeMissing", "false")
	}

	res, err := req.Get(path)
	if err != nil {
		return MeasurementsWrapper{}, err
	}

	if res.StatusCode() == http.StatusNoContent {
		return MeasurementsWrapper{Measurements: []MeasurementDto{}}, nil
	}

	if res.StatusCode() != http.StatusOK {
		return MeasurementsWrapper{}, fmt.Errorf("failed to get measurements: %s (status %d)", res.String(), res.StatusCode())
	}

	return result, nil
}
