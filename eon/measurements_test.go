package eon

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetMeasurements(t *testing.T) {
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

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	t.Run("successfully retrieves measurements", func(t *testing.T) {
		value1 := 10.5
		value2 := 12.3
		response := MeasurementsWrapper{
			ID:         12345,
			Resolution: "hour",
			Measurements: []MeasurementDto{
				{
					TimeStamp: FlexibleTime{Time: from},
					Value:     &value1,
				},
				{
					TimeStamp: FlexibleTime{Time: from.Add(time.Hour)},
					Value:     &value2,
				},
			},
		}

		httpmock.RegisterResponder("GET", "/measurements/12345/resolution/hour",
			func(req *http.Request) (*http.Response, error) {
				// Verify query parameters
				assert.Equal(t, "2024-01-01T00:00:00.000Z", req.URL.Query().Get("from"))
				assert.Equal(t, "2024-01-02T00:00:00.000Z", req.URL.Query().Get("to"))
				assert.Equal(t, "false", req.URL.Query().Get("includeMissing"))

				resp, err := httpmock.NewJsonResponse(200, response)
				return resp, err
			})

		result, err := c.GetMeasurements(12345, Hour, from, to, false)

		assert.NoError(t, err)
		assert.Equal(t, 12345, result.ID)
		assert.Equal(t, "hour", result.Resolution)
		assert.Len(t, result.Measurements, 2)
		assert.Equal(t, 10.5, *result.Measurements[0].Value)
	})

	t.Run("handles measurements with includeMissing=true", func(t *testing.T) {
		httpmock.Reset()
		response := MeasurementsWrapper{
			ID:           12345,
			Resolution:   "day",
			Measurements: []MeasurementDto{},
		}

		httpmock.RegisterResponder("GET", "/measurements/12345/resolution/day",
			func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "true", req.URL.Query().Get("includeMissing"))
				resp, err := httpmock.NewJsonResponse(200, response)
				return resp, err
			})

		result, err := c.GetMeasurements(12345, Day, from, to, true)

		assert.NoError(t, err)
		assert.Equal(t, 12345, result.ID)
	})

	t.Run("handles measurements with null values", func(t *testing.T) {
		httpmock.Reset()
		response := MeasurementsWrapper{
			ID:         12345,
			Resolution: "hour",
			Measurements: []MeasurementDto{
				{
					TimeStamp: FlexibleTime{Time: from},
					Value:     nil, // Missing measurement
				},
			},
		}

		httpmock.RegisterResponder("GET", "/measurements/12345/resolution/hour",
			httpmock.NewJsonResponderOrPanic(200, response))

		result, err := c.GetMeasurements(12345, Hour, from, to, false)

		assert.NoError(t, err)
		assert.Len(t, result.Measurements, 1)
		assert.Nil(t, result.Measurements[0].Value)
	})

	t.Run("handles 204 no content", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "/measurements/12345/resolution/hour",
			httpmock.NewStringResponder(204, ""))

		result, err := c.GetMeasurements(12345, Hour, from, to, false)

		assert.NoError(t, err)
		assert.Empty(t, result.Measurements)
	})

	t.Run("handles error response", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "/measurements/12345/resolution/hour",
			httpmock.NewStringResponder(400, `{"error":"bad request"}`))

		result, err := c.GetMeasurements(12345, Hour, from, to, false)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get measurements")
		assert.Empty(t, result.Measurements)
	})

	t.Run("handles different resolutions", func(t *testing.T) {
		resolutions := []Resolution{Quarter, Hour, Day, Month}

		for _, resolution := range resolutions {
			httpmock.Reset()
			response := MeasurementsWrapper{
				ID:           12345,
				Resolution:   string(resolution),
				Measurements: []MeasurementDto{},
			}

			httpmock.RegisterResponder("GET", "/measurements/12345/resolution/"+string(resolution),
				httpmock.NewJsonResponderOrPanic(200, response))

			result, err := c.GetMeasurements(12345, resolution, from, to, false)

			assert.NoError(t, err)
			assert.Equal(t, string(resolution), result.Resolution)
		}
	})
}
