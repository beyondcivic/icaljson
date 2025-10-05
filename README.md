# icaljson

[![Version](https://img.shields.io/badge/version-v0.8.2-blue)](https://github.com/beyondcivic/icaljson/releases/tag/v0.8.2)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/doc/devel/release.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/beyondcivic/icaljson.svg)](https://pkg.go.dev/github.com/beyondcivic/icaljson)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A Go implementation for parsing iCalendar format (.ics) files and converting them into structured JSON format. This library simplifies working with calendar data by providing both a command-line interface and a Go library for iCalendar processing.

## Overview

iCalendar is a standard format for exchanging calendar and scheduling information. This tool streamlines the process of converting iCalendar data to JSON format by:

- **Parsing standard iCalendar files** (.ics format) with full RFC 5545 compliance
- **Converting to clean JSON structure** with organized event data
- **Preserving all calendar metadata** including timezones, locations, and descriptions
- **Handling complex event properties** like geographic coordinates and recurring events
- **Providing both CLI and library interfaces** for different use cases

This project provides both a command-line interface and a Go library for working with iCalendar data.

## Key Features

- ✅ **iCalendar Parsing**: Full RFC 5545 compliant .ics file parsing
- ✅ **JSON Conversion**: Clean, structured JSON output format
- ✅ **Timezone Support**: Proper handling of timezone information
- ✅ **Event Properties**: Complete support for all standard event properties
- ✅ **Geographic Data**: Parse and convert GEO coordinates
- ✅ **CLI & Library**: Both command-line tool and Go library interfaces
- ✅ **Cross-platform**: Works on Linux, macOS, and Windows

## Getting Started

### Prerequisites

- Go 1.24 or later
- Nix 2.25.4 or later (optional but recommended)
- PowerShell v7.5.1 or later (for building)

### Installation

#### Option 1: Install from Source

1. Clone the repository:

```bash
git clone https://github.com/beyondcivic/icaljson.git
cd icaljson
```

2. Build the application:

```bash
go build -o icaljson .
```

#### Option 2: Using Nix (Recommended)

1. Clone the repository:

```bash
git clone https://github.com/beyondcivic/icaljson.git
cd icaljson
```

2. Prepare the environment using Nix flakes:

```bash
nix develop
```

3. Build the application:

```bash
./build.ps1
```

#### Option 3: Go Install

```bash
go install github.com/beyondcivic/icaljson@latest
```

## Quick Start

### Command Line Interface

The `icaljson` tool provides commands for converting iCalendar files to JSON:

```bash
# Convert iCalendar to JSON
icaljson generate calendar.ics -o calendar.json

# Show version information
icaljson version
```

### Go Library Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/beyondcivic/icaljson/pkg/icaljson"
)

func main() {
	// Convert iCalendar to JSON
	calendar, err := icaljson.Generate("calendar.ics", "calendar.json")
	if err != nil {
		log.Fatalf("Error converting calendar: %v", err)
	}

	fmt.Printf("Converted calendar with %d events\n", len(calendar.Events))

	// Access event data
	for _, event := range calendar.Events {
		fmt.Printf("Event: %s at %s\n", event.Summary, event.Location)
	}
}
```

## Detailed Command Reference

### `generate` - Convert iCalendar to JSON

Convert an iCalendar (.ics) file to structured JSON format.

```bash
icaljson generate [ICS_FILE] [OPTIONS]
```

**Options:**

- `-o, --output`: Output file path (default: `[filename].json`)

**Examples:**

```bash
# Basic conversion
icaljson generate events.ics

# With custom output path
icaljson generate events.ics -o my-events.json
```

### `version` - Show Version Information

Display version, build information, and system details.

```bash
icaljson version
```

## JSON Output Format

The tool converts iCalendar data to a clean JSON structure:

```json
{
  "prodid": "-//Calendar Producer//Calendar Product//EN",
  "version": "2.0",
  "calscale": "GREGORIAN",
  "events": [
    {
      "uid": "unique-event-id",
      "start": "2025-10-04T09:00:00Z",
      "end": "2025-10-04T10:00:00Z",
      "summary": "Meeting Title",
      "description": "Event description",
      "location": "Meeting Room A",
      "url": "https://example.com/event",
      "geo": {
        "latitude": 47.378177,
        "longitude": 8.540192
      },
      "comment": "Additional notes"
    }
  ]
}
```

## Supported iCalendar Properties

The parser supports standard iCalendar properties:

| iCalendar Property | JSON Field    | Description                 |
| ------------------ | ------------- | --------------------------- |
| `SUMMARY`          | `summary`     | Event title                 |
| `DESCRIPTION`      | `description` | Event description           |
| `DTSTART`          | `start`       | Event start time (ISO 8601) |
| `DTEND`            | `end`         | Event end time (ISO 8601)   |
| `LOCATION`         | `location`    | Event location              |
| `GEO`              | `geo`         | Geographic coordinates      |
| `URL`              | `url`         | Associated URL              |
| `UID`              | `uid`         | Unique identifier           |
| `COMMENT`          | `comment`     | Additional comments         |

## Examples

### Example 1: Basic Calendar Conversion

```bash
# Convert a simple calendar file
$ icaljson generate my-calendar.ics -o calendar.json

Generating JSON file for 'my-calendar.ics'...
✓ JSON file generated successfully and saved to: calendar.json
```

### Example 2: Processing Multiple Events

Given an iCalendar file with multiple events, the tool will create a JSON array containing all events with their properties properly converted and formatted.

## API Reference

### Core Functions

#### `Generate(icsPath, outputPath string) (*Calendar, error)`

Converts an iCalendar file to JSON format.

**Parameters:**

- `icsPath`: Path to the input .ics file
- `outputPath`: Path for the output JSON file

**Returns:**

- `*Calendar`: Parsed calendar structure
- `error`: Any error that occurred during processing

### Data Structures

#### `Calendar`

Represents the complete calendar structure:

```go
type Calendar struct {
    Prodid   string  `json:"prodid"`
    Version  string  `json:"version"`
    Calscale string  `json:"calscale"`
    Events   []Event `json:"events"`
}
```

#### `Event`

Represents an individual calendar event:

```go
type Event struct {
    UID         string     `json:"uid"`
    Start       string     `json:"start"`
    End         string     `json:"end"`
    Summary     string     `json:"summary"`
    Description string     `json:"description,omitempty"`
    Location    string     `json:"location,omitempty"`
    URL         string     `json:"url,omitempty"`
    Geo         *GeoPoint  `json:"geo,omitempty"`
    Comment     string     `json:"comment,omitempty"`
}
```

#### `GeoPoint`

Represents geographic coordinates:

```go
type GeoPoint struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}
```

## Architecture

The library is organized into several key components:

### Core Package (`pkg/icaljson`)

- **Parsing**: iCalendar file parsing and validation
- **Conversion**: iCalendar to JSON structure conversion
- **Data Types**: Calendar and event data structures
- **Utilities**: Helper functions for file handling and validation

### Command Line Interface (`cmd/icaljson`)

- **Cobra-based CLI** with subcommands for each major function
- **Comprehensive help system** with detailed usage examples
- **Flexible output options** and error handling

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes and add tests
4. Ensure all tests pass: `go test ./...`
5. Commit your changes: `git commit -am 'Add new feature'`
6. Push to the branch: `git push origin feature/new-feature`
7. Submit a pull request

### Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Build Environment

### Using Nix (Recommended)

Use Nix flakes to set up the build environment:

```bash
nix develop
```

### Manual Build

Check the build arguments in `build.ps1`:

```bash
# Build static binary with version information
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
```

Then run:

```bash
./build.ps1
```

Or build manually:

```bash
go build -o icaljson .
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
