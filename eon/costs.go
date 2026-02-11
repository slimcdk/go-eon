package eon

import (
	"fmt"
	"net/http"
	"time"
)

// GetCosts retrieves cost data for a specific installation.
// Whole months are considered for the time range.
//
// The API returns different schemas based on energy type:
//   - CostsElectricityWrapper for electricity
//   - CostsProductionWrapper for production
//   - CostsHeatWrapper for heat
//   - CostsColdWrapper for cold
//   - CostsGasWrapper for gas
//
// Use type assertion to access the specific type:
//
//	result, err := client.GetCosts("installation-id", &from, &to)
//	if wrapper, ok := result.(map[string]interface{}); ok {
//	    // Process the result
//	}
//
// Example:
//
//	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//	to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
//	costs, err := client.GetCosts("installation-id", &from, &to)
func (c *client) GetCosts(installationID string, from, to *time.Time) (interface{}, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	var result interface{}

	path := fmt.Sprintf("/costs/%s", installationID)

	req := c.resty.R().
		SetAuthToken(accessToken).
		SetResult(&result)

	// Add time range parameters if provided
	if from != nil {
		req.SetQueryParam("from", from.Format(time.RFC3339))
	}
	if to != nil {
		req.SetQueryParam("to", to.Format(time.RFC3339))
	}

	res, err := req.Get(path)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == http.StatusNoContent {
		return nil, fmt.Errorf("no cost data available for installation %s", installationID)
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get costs: %s (status %d)", res.String(), res.StatusCode())
	}

	return result, nil
}
