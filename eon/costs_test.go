package eon

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetCosts(t *testing.T) {
	mockResty := resty.New()
	httpmock.ActivateNonDefault(mockResty.GetClient())
	defer httpmock.DeactivateAndReset()

	c := &client{
		clientID:     "test-client-id",
		clientSecret: "test-client-secret",
		accessToken:  "fake-token",
		tokenExpiry:  time.Now().Add(1 * time.Hour), // Set valid token expiry
		resty:        mockResty,
	}

	t.Run("successfully retrieves costs without time range", func(t *testing.T) {
		response := map[string]interface{}{
			"installationId": "inst-1",
			"totalCost":      1250.50,
			"currency":       "SEK",
		}

		httpmock.RegisterResponder("GET", "/costs/inst-1",
			func(req *http.Request) (*http.Response, error) {
				// Verify no time parameters when nil
				assert.Empty(t, req.URL.Query().Get("from"))
				assert.Empty(t, req.URL.Query().Get("to"))

				resp, err := httpmock.NewJsonResponse(200, response)
				return resp, err
			})

		result, err := c.GetCosts("inst-1", nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("successfully retrieves costs with time range", func(t *testing.T) {
		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

		response := map[string]interface{}{
			"installationId": "inst-1",
			"totalCost":      1250.50,
			"currency":       "SEK",
			"from":           "2024-01-01T00:00:00",
			"to":             "2024-12-31T00:00:00",
		}

		httpmock.RegisterResponder("GET", "/costs/inst-1",
			func(req *http.Request) (*http.Response, error) {
				// Verify time parameters (with Z suffix for UTC)
				assert.Equal(t, "2024-01-01T00:00:00Z", req.URL.Query().Get("from"))
				assert.Equal(t, "2024-12-31T00:00:00Z", req.URL.Query().Get("to"))

				resp, err := httpmock.NewJsonResponse(200, response)
				return resp, err
			})

		result, err := c.GetCosts("inst-1", &from, &to)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("handles electricity costs response", func(t *testing.T) {
		httpmock.Reset()
		retailCost := 75.50
		retailCostVAT := 19.50
		response := CostsElectricityWrapper{
			CostsWrapper: CostsWrapper{
				EnergyClass:  "electricity",
				Installation: "inst-1",
			},
			Costs: []CostElectricityProductionDto{
				{
					RetailCost:    &retailCost,
					RetailCostVAT: &retailCostVAT,
				},
			},
		}

		httpmock.RegisterResponder("GET", "/costs/inst-1",
			httpmock.NewJsonResponderOrPanic(200, response))

		result, err := c.GetCosts("inst-1", nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("handles 204 no content", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "/costs/inst-1",
			httpmock.NewStringResponder(204, ""))

		result, err := c.GetCosts("inst-1", nil, nil)

		assert.Error(t, err) // 204 returns an error per the implementation
		assert.Contains(t, err.Error(), "no cost data available")
		assert.Nil(t, result)
	})

	t.Run("handles error response", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "/costs/inst-1",
			httpmock.NewStringResponder(404, `{"error":"not found"}`))

		result, err := c.GetCosts("inst-1", nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get costs")
		assert.Nil(t, result)
	})
}
