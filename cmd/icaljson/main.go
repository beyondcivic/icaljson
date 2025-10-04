// Package cmd provides the command-line interface for icaljson.
//
// icaljson is a Go implementation for working with iCalendar format (RFC 5545).
//
// The command-line tool provides functionality to:
//   - Generate JSON from iCal files with automatic type inference
//   - Display version and build information
//
// # Command Reference
//
// Generate json with default output path:
//
//	icaljson generate caledar.ics
//
// Show version information:
//
//	icaljson version
//
// # Features
//
// Metadata Generation:
//   - Automatic data type inference from ICS file
//   - Configurable output paths and validation options
//   - Support for environment variable configuration
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/beyondcivic/icaljson/pkg/icaljson"
	"github.com/beyondcivic/icaljson/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Root cobra command.
// Call Init() once to initialize child commands.
// Global so it can be picked up by docs/doc-gen.go.
// nolint:gochecknoglobals
var RootCmd = &cobra.Command{
	Use:   "icaljson",
	Short: "iCalendar tools",
	Long: `A Go implementation for working with the iCalendar format.
iCalendar is a standardized way to describe calendar data using a text-based format.`,
	Version: version.Version,
}

// Call Once.
func Init() {
	// Initialize viper for configuration
	viper.SetEnvPrefix("ICALJSON")
	viper.AutomaticEnv()

	// Add child commands
	RootCmd.AddCommand(versionCmd())
	RootCmd.AddCommand(generateCmd())
}

func Execute() {
	// Execute the command
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Helper functions

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isICalFile(filename string) bool {
	return icaljson.IsICalFile(filename)
}

func determineOutputPath(providedPath, csvPath string) string {
	if providedPath != "" {
		return providedPath
	}

	// Check environment variable
	envOutputPath := os.Getenv("ICALJSON_OUTPUT_PATH")
	if envOutputPath != "" {
		return envOutputPath
	}

	// Generate default path based on CSV filename
	baseName := strings.TrimSuffix(filepath.Base(csvPath), filepath.Ext(csvPath))
	return baseName + "_parsed.json"
}
