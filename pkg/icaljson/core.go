package icaljson

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Generate generates JSON file from a ICS file with automatic type inference.
func Generate(icsPath string, outputPath string) (*Calendar, error) {
	// Get file information
	// fileName := filepath.Base(icsPath)
	_, err := os.Stat(icsPath)
	if err != nil {
		return nil, AppError{Message: "failed to get file info", Value: err}
	}
	// fileSize := fileInfo.Size()

	// calendarFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// Parse ICS file
	var calendar Calendar

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

	return &calendar, nil
}
