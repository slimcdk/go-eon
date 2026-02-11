package eon

import "time"

type Resolution string

const (
	// Eon API endpoints
	tokenEndpoint = "https://navigator-api.eon.se/connect/token"
	apiBaseURL    = "https://navigator-api.eon.se/api"
)

var (
	// Timezone for Swedish/Copenhagen operations
	stockholmTZ, _ = time.LoadLocation("Europe/Stockholm")
	cph, _         = time.LoadLocation("Europe/Copenhagen")
)

// Resolution types supported by Eon API
const (
	Quarter Resolution = "quarter" // 15-minute intervals
	Hour    Resolution = "hour"
	Day     Resolution = "day"
	Month   Resolution = "month"
)

const (
	MaximumDayRequestLeap  int           = 730
	MaximumRequestDuration time.Duration = time.Hour * 24 * 730
)
