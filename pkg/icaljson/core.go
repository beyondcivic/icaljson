package icaljson

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		name, value, err := parseProperty(line)
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
			case "DTSTAMP":
				currentEvent.DTStamp = value

			// Date/Time properties
			case "DTSTART", "DTSTART;TZID":
				currentEvent.Start = value
			case "DTEND", "DTEND;TZID":
				currentEvent.End = value
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

			// Other properties
			case "GEO":
				currentEvent.Geo = value
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

// parseProperty parses a property line into name and value
// Handles properties with parameters like "DTSTART;TZID=Europe/Zurich:20251004T090000"
func parseProperty(line string) (string, string, error) {
	// Find the colon that separates name from value
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return "", "", AppError{Message: "invalid property line: no colon found"}
	}

	nameWithParams := line[:colonIndex]
	value := line[colonIndex+1:]

	// Extract just the property name (before any semicolon)
	// For "DTSTART;TZID=Europe/Zurich", we want "DTSTART"
	// But we'll keep the full nameWithParams for matching specific cases
	name := nameWithParams
	if semicolonIndex := strings.Index(nameWithParams, ";"); semicolonIndex != -1 {
		// For properties with parameters, we'll use a simplified name
		baseName := nameWithParams[:semicolonIndex]
		// For DTSTART and DTEND with timezone, we create a composite key
		if baseName == "DTSTART" || baseName == "DTEND" {
			name = baseName + ";TZID"
		} else {
			name = baseName
		}
	}

	return name, value, nil
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
