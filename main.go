// icaljson is a command-line tool and Go library for working with the iCalendar format.
//
// iCalendar is a standardized way to describe calendar data using a text-based format.
// This tool simplifies the creation of structured JSON from iCalendar files.
//
// # Installation
//
// Install the latest version:
//
//	go install github.com/beyondcivic/icaljson@latest
//
// # Usage
//
// Generate JSON from a ICS file:
//
//	icaljson generate myevents.ics -o myevents.json
//
// For detailed usage information, run:
//
//	icaljson --help
package main

import (
	cmd "github.com/beyondcivic/icaljson/cmd/icaljson"
)

func main() {
	cmd.Init()
	cmd.Execute()
}
