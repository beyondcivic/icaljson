// commands.go
// Contains cobra command definitions
//
//nolint:funlen,mnd
package cmd

import (
	"fmt"
	"os"

	"github.com/beyondcivic/icaljson/pkg/icaljson"
	"github.com/beyondcivic/icaljson/pkg/version"
	"github.com/spf13/cobra"
)

// Version Command.
// Displays tool version and build information
func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version, git hash, and build time information of the icaljson tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s\n", version.AppName, version.Version)
			stamp := version.RetrieveStamp()
			fmt.Printf("  Built with %s on %s\n", stamp.InfoGoCompiler, stamp.InfoBuildTime)
			fmt.Printf("  Git ref: %s\n", stamp.VCSRevision)
			fmt.Printf("  Go version %s, GOOS %s, GOARCH %s\n", stamp.InfoGoVersion, stamp.InfoGOOS, stamp.InfoGOARCH)
		},
	}
}

// Generate command
func generateCmd() *cobra.Command {
	var generateCmd = &cobra.Command{
		Use:   "generate [icsPath]",
		Short: "Generate JSON from a ICS file",
		Long:  `Generate JSON from a ICS file, automatically inferring data types.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			icsPath := args[0]
			flagOutputPath, _ := cmd.Flags().GetString("output")

			// Validate input file
			if !fileExists(icsPath) {
				fmt.Printf("Error: ICS file '%s' does not exist.\n", icsPath)
				os.Exit(1)
			}

			if !isICalFile(icsPath) {
				fmt.Printf("Error: File '%s' does not appear to be a ICS file.\n", icsPath)
				os.Exit(1)
			}

			// Determine output path
			outputPath := determineOutputPath(flagOutputPath, icsPath)

			// Validate output path
			if err := icaljson.ValidateOutputPath(outputPath); err != nil {
				fmt.Printf("Error: Invalid output path: %v\n", err)
				os.Exit(1)
			}

			// Generate metadata
			fmt.Printf("Generating JSON file for '%s'...\n", icsPath)
			_, err := icaljson.Generate(icsPath, outputPath)
			if err != nil {
				fmt.Printf("Error generating metadata: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✓ JSON file generated successfully")
			if outputPath != "" {
				fmt.Printf(" and saved to: %s\n", outputPath)
			}

		},
	}
	generateCmd.Flags().StringP("output", "o", "", "Output path for the JSON file")

	return generateCmd
}
