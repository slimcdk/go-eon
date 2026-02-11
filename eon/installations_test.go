package eon

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetInstallations(t *testing.T) {
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

	t.Run("successfully retrieves installations without filter", func(t *testing.T) {
		response := InstallationsWrapper{
			Installations: []InstallationDto{
				{
					ID:       "inst-1",
					Active:   true,
					Name:     "Test Installation",
					Address:  "Test Street 1",
					City:     "Stockholm",
					GridArea: "SE3",
				},
			},
		}

		httpmock.RegisterResponder("GET", "/installations",
			httpmock.NewJsonResponderOrPanic(200, response))

		result, err := c.GetInstallations(nil)

		assert.NoError(t, err)
		assert.Len(t, result.Installations, 1)
		assert.Equal(t, "inst-1", result.Installations[0].ID)
		assert.Equal(t, "Test Installation", result.Installations[0].Name)
	})

	t.Run("successfully retrieves installations with filter", func(t *testing.T) {
		httpmock.Reset()
		response := InstallationsWrapper{
			Installations: []InstallationDto{
				{ID: "inst-1", Name: "Filtered Installation"},
			},
		}

		httpmock.RegisterResponder("GET", "/installations",
			func(req *http.Request) (*http.Response, error) {
				// Verify filter query parameter
				filters := req.URL.Query()["installationFilter"]
				assert.Contains(t, filters, "inst-1")
				assert.Contains(t, filters, "inst-2")

				resp, err := httpmock.NewJsonResponse(200, response)
				return resp, err
			})

		result, err := c.GetInstallations([]string{"inst-1", "inst-2"})

		assert.NoError(t, err)
		assert.Len(t, result.Installations, 1)
	})

	t.Run("handles 204 no content", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "/installations",
			httpmock.NewStringResponder(204, ""))

		result, err := c.GetInstallations(nil)

		assert.NoError(t, err)
		assert.Empty(t, result.Installations)
	})

	t.Run("handles error response", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "/installations",
			httpmock.NewStringResponder(401, `{"error":"unauthorized"}`))

		result, err := c.GetInstallations(nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get installations")
		assert.Empty(t, result.Installations)
	})
}

func TestGetMeasurementSeries(t *testing.T) {
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

	t.Run("successfully retrieves measurement series", func(t *testing.T) {
		response := InstallationsMeasurementsWrapper{
			Installations: []InstallationMeasurementsDto{
				{
					ID: "inst-1",
					MeasurementSeries: []MeasurementSeriesDto{
						{
							ID:         12345,
							SeriesType: "electricity",
							Unit:       "kWh",
						},
					},
				},
			},
		}

		httpmock.RegisterResponder("GET", "/installations/measurement-series",
			httpmock.NewJsonResponderOrPanic(200, response))

		result, err := c.GetMeasurementSeries()

		assert.NoError(t, err)
		assert.Len(t, result.Installations, 1)
		assert.Equal(t, "inst-1", result.Installations[0].ID)
		assert.Len(t, result.Installations[0].MeasurementSeries, 1)
		assert.Equal(t, 12345, result.Installations[0].MeasurementSeries[0].ID)
	})

	t.Run("handles 204 no content", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "/installations/measurement-series",
			httpmock.NewStringResponder(204, ""))

		result, err := c.GetMeasurementSeries()

		assert.NoError(t, err)
		assert.Empty(t, result.Installations)
	})

	t.Run("handles error response", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", "/installations/measurement-series",
			httpmock.NewStringResponder(500, `{"error":"server error"}`))

		_, err := c.GetMeasurementSeries()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get measurement series")
	})
}
