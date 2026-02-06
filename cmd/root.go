package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "groovekit",
	Short: "Monitor cron jobs and APIs from your terminal",
	Long: `GrooveKit CLI - Monitor your cron jobs and API endpoints before users notice.

Verify your services are working correctly with heartbeat monitoring,
JSON Schema validation, GraphQL support, and instant alerts.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
}
