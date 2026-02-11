# Go Client for Eon Energy Navigator API

[![Tests](https://github.com/slimcdk/go-eon/workflows/Tests/badge.svg)](https://github.com/slimcdk/go-eon/actions?query=workflow%3ATests)
[![Go Report Card](https://goreportcard.com/badge/github.com/slimcdk/go-eon)](https://goreportcard.com/report/github.com/slimcdk/go-eon)
[![Go Reference](https://pkg.go.dev/badge/github.com/slimcdk/go-eon.svg)](https://pkg.go.dev/github.com/slimcdk/go-eon)
[![License](https://img.shields.io/github/license/slimcdk/go-eon)](LICENSE)

A comprehensive Go client library and CLI tool for the [Eon Energy Navigator API](https://navigator-api.eon.se/). Access electricity installations, consumption measurements, and cost data from Eon, a leading Swedish energy provider.

## Features

- **OAuth2 Authentication**: Automatic token management with client credentials flow
- **Complete API Coverage**: Access installations, measurement series, measurements, and costs
- **Well-Tested**: 89.5% test coverage with comprehensive unit tests
- **Easy to Use**: Simple CLI interface and intuitive Go library
- **Flexible**: Automatic token refresh and error handling
- **Multi-Platform**: Cross-compiled binaries for Linux, macOS, and Windows

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Getting API Credentials](#getting-api-credentials)
  - [CLI Usage](#cli-usage)
  - [Library Usage](#library-usage)
- [API Coverage](#api-coverage)
- [CLI Reference](#cli-reference)
- [Library Reference](#library-reference)
- [Examples](#examples)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Using Go Install

```bash
go install github.com/slimcdk/go-eon@latest
```

### Download Pre-built Binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/slimcdk/go-eon/releases).

### Build from Source

```bash
git clone https://github.com/slimcdk/go-eon.git
cd go-eon
go build -o eon .
```

## Quick Start

### Getting API Credentials

To use the Eon Energy Navigator API, you need OAuth2 client credentials:

1. Visit the [Eon Energy Navigator API portal](https://navigator-api.eon.se/)
2. Register your application to obtain:
   - Client ID
   - Client Secret

### CLI Usage

```bash
# Set credentials as environment variables
export CLIENT_ID="your-client-id"
export CLIENT_SECRET="your-client-secret"

# Get your installations
eon installations

# Get measurement series for installations
eon measurement-series

# Get measurement data (requires series ID from measurement-series command)
eon measurements 737605 \
  --from=2024-01-01 \
  --to=2024-01-31 \
  --resolution=hour

# Get costs for an installation
eon costs 735999163005019944 \
  --from=2024-01-01 \
  --to=2024-12-31
```

Alternatively, provide credentials via command-line flags:

```bash
eon --client-id="your-id" --client-secret="your-secret" installations
```

### Library Usage

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/slimcdk/go-eon/eon"
)

func main() {
    // Create client with credentials from environment variables
    // CLIENT_ID and CLIENT_SECRET
    client := eon.New()

    // Or create with explicit credentials
    // client := eon.NewWithCredentials("client-id", "client-secret")

    // Get installations
    installations, err := client.GetInstallations(nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, inst := range installations.Installations {
        fmt.Printf("Installation: %s - %s\n", inst.ID, inst.Name)
        fmt.Printf("Address: %s, %s\n", inst.Address, inst.City)
        fmt.Printf("Active: %v\n", inst.Active)
    }

    // Get measurement series
    series, err := client.GetMeasurementSeries()
    if err != nil {
        log.Fatal(err)
    }

    for _, inst := range series.Installations {
        fmt.Printf("Installation %s has %d measurement series\n",
            inst.ID, len(inst.MeasurementSeries))

        for _, ms := range inst.MeasurementSeries {
            fmt.Printf("  Series %d: %s (%s)\n",
                ms.ID, ms.SeriesType, ms.Unit)
        }
    }

    // Get measurements for a series
    from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    to := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

    measurements, err := client.GetMeasurements(
        737605,           // measurement series ID
        eon.Hour,         // resolution
        from,
        to,
        false,            // include missing values
    )
    if err != nil {
        log.Fatal(err)
    }

    for _, m := range measurements.Measurements {
        if m.Value != nil {
            fmt.Printf("%s: %.3f\n", m.TimeStamp.Time.Format(time.RFC3339), *m.Value)
        }
    }
}
```

## API Coverage

### Eon Energy Navigator API

| Endpoint | CLI Command | Library Method | Description |
|----------|-------------|----------------|-------------|
| `/connect/token` | - | `GetAccessToken()` | OAuth2 authentication |
| `/api/installations` | `installations` | `GetInstallations()` | List installations with optional filter |
| `/api/installations/measurement-series` | `measurement-series` | `GetMeasurementSeries()` | List measurement series for installations |
| `/api/measurements/{id}/resolution/{resolution}` | `measurements` | `GetMeasurements()` | Get measurement data with time range |
| `/api/costs/{id}` | `costs` | `GetCosts()` | Get cost data for installation |
| `/api/` | - | `IsAlive()` | Health check |

## CLI Reference

### Global Flags

```
--client-id string       Eon API client ID (env: CLIENT_ID)
--client-secret string   Eon API client secret (env: CLIENT_SECRET)
```

### Commands

```bash
# Get all installations
eon installations

# Get installations filtered by ID
eon installations --filter=735999163005019944

# Get measurement series for all installations
eon measurement-series

# Get measurements for a specific series
eon measurements <series-id> \
  --from=YYYY-MM-DD \
  --to=YYYY-MM-DD \
  --resolution=hour \           # quarter, hour, day, month
  --include-missing             # Fill in missing values

# Get costs for an installation
eon costs <installation-id> \
  --from=YYYY-MM-DD \
  --to=YYYY-MM-DD
```

### Resolution Options

- **quarter**: 15-minute intervals (requires from/to, max 3 months)
- **hour**: Hourly values (requires from/to, max 1 year)
- **day**: Daily values
- **month**: Monthly values

## Library Reference

### Creating Clients

```go
import "github.com/slimcdk/go-eon/eon"

// Create client using environment variables CLIENT_ID and CLIENT_SECRET
client := eon.New()

// Create client with explicit credentials
client := eon.NewWithCredentials("client-id", "client-secret")
```

### Resolution Types

```go
eon.Quarter  // 15-minute intervals
eon.Hour     // Hourly aggregation
eon.Day      // Daily aggregation
eon.Month    // Monthly aggregation
```

### Client Interface

```go
type Client interface {
    GetAccessToken() (string, error)
    GetInstallations(filter []string) (InstallationsWrapper, error)
    GetMeasurementSeries() (InstallationsMeasurementsWrapper, error)
    GetMeasurements(id int, resolution Resolution, from, to time.Time, includeMissing bool) (MeasurementsWrapper, error)
    GetCosts(installationID string, from, to *time.Time) (interface{}, error)
    IsAlive() (bool, error)
}
```

## Examples

### Two-Step Workflow: Get Measurements

The Eon API requires a two-step process to get measurements:

1. **First, get the measurement series ID:**

```bash
$ eon measurement-series
{
  "installations": [
    {
      "id": "735999163005019944",
      "measurementSeries": [
        {
          "id": 737605,
          "seriesType": "ElectricActive",
          "unit": "KWH",
          "lastUpdate": "2024-02-11T10:00:00"
        }
      ]
    }
  ]
}
```

2. **Then, use the series ID (737605) to get measurements:**

```bash
$ eon measurements 737605 --from=2024-01-01 --to=2024-01-31 --resolution=hour
```

### Filter Installations

```go
// Get specific installations by ID
installations, err := client.GetInstallations([]string{
    "735999163005019944",
    "735999163005019945",
})
if err != nil {
    log.Fatal(err)
}
```

### Handle Missing Data

```go
// Include missing values (filled with null)
measurements, err := client.GetMeasurements(
    737605,
    eon.Hour,
    from,
    to,
    true, // includeMissing
)

for _, m := range measurements.Measurements {
    if m.Value == nil {
        fmt.Printf("%s: <no data>\n", m.TimeStamp.Time.Format(time.RFC3339))
    } else {
        fmt.Printf("%s: %.3f\n", m.TimeStamp.Time.Format(time.RFC3339), *m.Value)
    }
}
```

### Get Costs with Time Range

```go
from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

costs, err := client.GetCosts("735999163005019944", &from, &to)
if err != nil {
    log.Fatal(err)
}

// Note: costs is interface{} and can be different types
// (CostsElectricityWrapper, CostsProductionWrapper, etc.)
```

### Error Handling

```go
measurements, err := client.GetMeasurements(737605, eon.Hour, from, to, false)
if err != nil {
    // Check for specific error types
    if strings.Contains(err.Error(), "authentication failed") {
        log.Fatal("Invalid credentials")
    }
    if strings.Contains(err.Error(), "status 404") {
        log.Fatal("Measurement series not found")
    }
    log.Fatal(err)
}
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./eon/...
go test -v ./cmd/...
```

### Running Linter

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run --timeout=5m
```

### Building

```bash
# Build for current platform
go build -o eon .

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o eon-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o eon-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o eon-windows-amd64.exe .
```

### Testing with Act (GitHub Actions locally)

```bash
# Install act: https://github.com/nektos/act
# Run workflow locally
act -j test
act -j lint
act -j build
```

## Project Structure

```
.
├── cmd/                    # CLI command implementations
│   ├── costs.go           # Costs commands
│   ├── installations.go   # Installations and measurement-series commands
│   ├── measurements.go    # Measurements commands
│   └── root.go            # Root command and initialization
├── eon/                   # Library implementation
│   ├── auth.go            # OAuth2 authentication
│   ├── constvars.go       # Constants and resolutions
│   ├── costs.go           # Costs endpoints
│   ├── eon.go             # Client initialization
│   ├── errors.go          # Error handling
│   ├── installations.go   # Installations endpoints
│   ├── interfaces.go      # Client interface
│   ├── measurements.go    # Measurements endpoints
│   ├── models.go          # Data models
│   ├── utils.go           # Utilities
│   └── *_test.go          # Unit tests
├── .github/
│   └── workflows/
│       └── test.yml       # CI/CD pipeline
├── .golangci.yml          # Linter configuration
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
└── main.go                # CLI entry point
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Run linter (`golangci-lint run`)
6. Format code (`go fmt ./...`)
7. Commit your changes (`git commit -m 'Add some amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Maintain or improve test coverage (currently 89.5%)
- Follow existing code style
- Update documentation for API changes
- Add examples for new functionality
- Use idiomatic Go (proper error bubbling, interface design)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Eon](https://www.eon.se/) for providing the Energy Navigator API
- [Eon Energy Navigator API](https://navigator-api.eon.se/) platform
- All contributors who have helped improve this project

## Support

- [Eon Website](https://www.eon.se/)
- [API Platform](https://navigator-api.eon.se/)
- [Report Issues](https://github.com/slimcdk/go-eon/issues)
- [Discussions](https://github.com/slimcdk/go-eon/discussions)

## Related Resources

- [Eon Sweden](https://www.eon.se/) - Main website for Eon's Swedish operations
- [Energy Navigator API Documentation](https://navigator-api.eon.se/) - API platform and documentation

---

Made for the Swedish energy community 🇸🇪
