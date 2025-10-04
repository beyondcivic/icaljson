package icaljson

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Generate generates JSON file from a ICS file with automatic type inference.
func Generate(icsPath string, outputPath string) (*Calendar, error) {
	// Get file information
	_, err := os.Stat(icsPath)
	if err != nil {
		return nil, AppError{Message: "failed to get file info", Value: err}
	}

	// Read and parse ICS file
	file, err := os.Open(icsPath)
	if err != nil {
		return nil, AppError{Message: "failed to open ICS file", Value: err}
	}
	defer file.Close()

	calendar, err := parseICS(file)
	if err != nil {
		return nil, AppError{Message: "failed to parse ICS file", Value: err}
	}

	// Write to file if output path is provided
	if outputPath != "" {
		// Marshal calendar to JSON with proper indentation
		metadataJSON, err := json.MarshalIndent(calendar, "", "  ")
		if err != nil {
			return nil, AppError{Message: "failed to marshal JSON", Value: err}
		}

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(outputPath), 0750); err != nil {
			return nil, AppError{Message: "failed to create directory", Value: err}
		}

		// Write metadata to file
		if err := os.WriteFile(outputPath, metadataJSON, 0600); err != nil {
			return nil, AppError{Message: "failed to write file", Value: err}
		}
	}

	return calendar, nil
}

// parseICS parses an ICS file according to RFC 5545
func parseICS(file *os.File) (*Calendar, error) {
	scanner := bufio.NewScanner(file)
	calendar := &Calendar{}
	var currentEvent *Event
	var lines []string

	// First pass: unfold lines (handle line continuation)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, AppError{Message: "failed to read file", Value: err}
	}

	// Unfold lines (lines starting with space or tab continue previous line)
	unfoldedLines := unfoldLines(lines)

	// Parse the unfolded lines
	for _, line := range unfoldedLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse component boundaries
		if strings.HasPrefix(line, "BEGIN:") {
			component := strings.TrimPrefix(line, "BEGIN:")
			if component == "VEVENT" {
				currentEvent = &Event{}
			}
			continue
		}

		if strings.HasPrefix(line, "END:") {
			component := strings.TrimPrefix(line, "END:")
			if component == "VEVENT" && currentEvent != nil {
				calendar.Events = append(calendar.Events, *currentEvent)
				currentEvent = nil
			}
			continue
		}

		// Parse properties
		name, value, tzid, err := parsePropertyWithTZ(line)
		if err != nil {
			continue // Skip malformed lines
		}

		// Handle calendar-level properties
		if currentEvent == nil {
			switch name {
			case "PRODID":
				calendar.ProdID = value
			case "VERSION":
				calendar.Version = value
			case "CALSCALE":
				calendar.CalScale = value
			case "METHOD":
				calendar.Method = value
			}
		} else {
			// Handle event properties
			switch name {
			// Required properties
			case "UID":
				currentEvent.UID = value

			// Date/Time properties - parse to ISO8601 UTC
			case "DTSTART", "DTSTART;TZID":
				if parsed := parseICalDateTimeWithTZ(value, tzid); parsed != "" {
					currentEvent.Start = parsed
				} else {
					currentEvent.Start = value // Fallback to raw value if parsing fails
				}
			case "DTEND", "DTEND;TZID":
				if parsed := parseICalDateTimeWithTZ(value, tzid); parsed != "" {
					currentEvent.End = parsed
				} else {
					currentEvent.End = value // Fallback to raw value if parsing fails
				}
			case "DURATION":
				currentEvent.Duration = value

			// Core descriptive properties
			case "SUMMARY":
				currentEvent.Summary = value
			case "DESCRIPTION":
				currentEvent.Description = unescapeText(value)
			case "LOCATION":
				currentEvent.Location = unescapeText(value)

			// Optional commonly used properties
			case "URL":
				currentEvent.URL = value
			case "STATUS":
				currentEvent.Status = strings.ToUpper(value)
			case "CATEGORIES":
				if value != "" {
					currentEvent.Categories = strings.Split(value, ",")
					for i := range currentEvent.Categories {
						currentEvent.Categories[i] = strings.TrimSpace(currentEvent.Categories[i])
					}
				}

			// Classification and access
			case "CLASS":
				currentEvent.Class = strings.ToUpper(value)
			case "TRANSP":
				currentEvent.Transp = strings.ToUpper(value)

			// Organizational properties
			case "ORGANIZER":
				currentEvent.Organizer = value
			case "ATTENDEE":
				currentEvent.Attendees = append(currentEvent.Attendees, value)

			// Scheduling properties
			case "PRIORITY":
				// PRIORITY is 0-9 integer
				if priority := parseInt(value); priority >= 0 {
					currentEvent.Priority = priority
				}
			case "SEQUENCE":
				if sequence := parseInt(value); sequence >= 0 {
					currentEvent.Sequence = sequence
				}

			// Date/Time metadata
			case "CREATED":
				currentEvent.Created = value
			case "LAST-MODIFIED":
				currentEvent.LastModified = value

			// Recurrence properties
			case "RRULE":
				currentEvent.RRule = value
			case "RECURRENCE-ID", "RECURRENCE-ID;TZID":
				currentEvent.RecurrenceID = value
			case "EXDATE", "EXDATE;TZID":
				currentEvent.ExDates = append(currentEvent.ExDates, value)
			case "RDATE", "RDATE;TZID":
				currentEvent.RDates = append(currentEvent.RDates, value)

			case "GEO":
				// Other properties
				if value != "" {
					parts := strings.Split(value, ";")
					if len(parts) == 2 {
						lat, err1 := strconv.ParseFloat(parts[0], 64)
						lon, err2 := strconv.ParseFloat(parts[1], 64)
						if err1 == nil && err2 == nil {
							currentEvent.Geo = Geolocation{
								Latitude:  lat,
								Longitude: lon,
							}
						}
					}
				}

			case "RESOURCES":
				if value != "" {
					resources := strings.Split(value, ",")
					for _, r := range resources {
						currentEvent.Resources = append(currentEvent.Resources, strings.TrimSpace(r))
					}
				}
			case "CONTACT":
				currentEvent.Contact = value
			case "RELATED-TO":
				currentEvent.RelatedTo = value
			case "COMMENT":
				currentEvent.Comment = unescapeText(value)
			}
		}
	}

	return calendar, nil
}

// parseICalDateTimeWithTZ converts iCalendar datetime format to ISO8601 UTC
// Supports formats like:
// - 20251004T090000Z (UTC)
// - 20251004T090000 (local time, will be converted to UTC if tzid provided)
// - 20251004 (date only)
func parseICalDateTimeWithTZ(icalDate string, tzid string) string {
	if icalDate == "" {
		return ""
	}

	// Try parsing as UTC datetime (with Z suffix) - already in UTC
	if t, err := time.Parse("20060102T150405Z", icalDate); err == nil {
		return t.Format(time.RFC3339)
	}

	// Try parsing as local datetime with timezone conversion
	if t, err := time.Parse("20060102T150405", icalDate); err == nil {
		// If timezone ID is provided, load the timezone and convert to UTC
		if tzid != "" {
			loc, err := time.LoadLocation(tzid)
			if err == nil {
				// Parse the time in the specified timezone
				localTime := time.Date(t.Year(), t.Month(), t.Day(),
					t.Hour(), t.Minute(), t.Second(), 0, loc)
				// Convert to UTC and format as RFC3339
				return localTime.UTC().Format(time.RFC3339)
			}
		}
		// No timezone or timezone load failed - return as-is without timezone info
		return t.Format("2006-01-02T15:04:05")
	}

	// Try parsing as date only
	if t, err := time.Parse("20060102", icalDate); err == nil {
		return t.Format("2006-01-02")
	}

	// Return empty string if parsing fails
	return ""
}

// unfoldLines handles RFC 5545 line folding where lines starting with space or tab
// are continuations of the previous line
func unfoldLines(lines []string) []string {
	var result []string
	var currentLine string

	for _, line := range lines {
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			// This is a continuation of the previous line
			currentLine += line[1:] // Remove the leading space/tab
		} else {
			// This is a new line
			if currentLine != "" {
				result = append(result, currentLine)
			}
			currentLine = line
		}
	}

	// Don't forget the last line
	if currentLine != "" {
		result = append(result, currentLine)
	}

	return result
}

// parsePropertyWithTZ parses a property line into name, value, and timezone ID
// Handles properties with parameters like "DTSTART;TZID=Europe/Zurich:20251004T090000"
func parsePropertyWithTZ(line string) (string, string, string, error) {
	// Find the colon that separates name from value
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return "", "", "", AppError{Message: "invalid property line: no colon found"}
	}

	nameWithParams := line[:colonIndex]
	value := line[colonIndex+1:]
	tzid := ""

	// Extract the property name and timezone if present
	// For "DTSTART;TZID=Europe/Zurich", we want "DTSTART", "Europe/Zurich"
	name := nameWithParams
	if semicolonIndex := strings.Index(nameWithParams, ";"); semicolonIndex != -1 {
		baseName := nameWithParams[:semicolonIndex]
		params := nameWithParams[semicolonIndex+1:]

		// Extract TZID parameter if present
		if strings.HasPrefix(params, "TZID=") {
			tzid = strings.TrimPrefix(params, "TZID=")
			// For DTSTART and DTEND with timezone, we create a composite key
			name = baseName + ";TZID"
		} else {
			name = baseName
		}
	}

	return name, value, tzid, nil
}

// unescapeText unescapes special characters in text values according to RFC 5545
// \n -> newline, \, -> comma, \; -> semicolon, \\ -> backslash
func unescapeText(text string) string {
	text = strings.ReplaceAll(text, "\\n", "\n")
	text = strings.ReplaceAll(text, "\\N", "\n")
	text = strings.ReplaceAll(text, "\\,", ",")
	text = strings.ReplaceAll(text, "\\;", ";")
	text = strings.ReplaceAll(text, "\\\\", "\\")
	return text
}

// parseInt safely parses a string to int, returning -1 on error
func parseInt(s string) int {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return -1
	}
	return result
}
