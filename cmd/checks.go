package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/scookdev/groovekit-cli/internal/api"
	"github.com/scookdev/groovekit-cli/internal/output"
	"github.com/spf13/cobra"
)

var checksCmd = &cobra.Command{
	Use:   "checks",
	Short: "View check and ping history",
	Long:  "View recent health check results for monitors and job heartbeat pings",
}

// checks list
var checksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent checks",
	Long:  "List recent health checks for a monitor or pings for a job",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		monitorID, _ := cmd.Flags().GetString("monitor")
		jobID, _ := cmd.Flags().GetString("job")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		// Must specify either --monitor or --job
		if monitorID == "" && jobID == "" {
			return fmt.Errorf("must specify either --monitor or --job")
		}

		if monitorID != "" && jobID != "" {
			return fmt.Errorf("cannot specify both --monitor and --job")
		}

		if monitorID != "" {
			return listMonitorChecks(client, monitorID, jsonOutput)
		}

		return listJobPings(client, jobID, jsonOutput)
	},
}

func listMonitorChecks(client *api.Client, monitorID string, jsonOutput bool) error {
	// Resolve short ID to full ID
	fullID, err := resolveMonitorID(client, monitorID)
	if err != nil {
		return err
	}

	var s *spinner.Spinner
	if !jsonOutput {
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
	}

	checks, err := client.ListMonitorChecks(fullID)

	if s != nil {
		s.Stop()
	}

	if err != nil {
		return fmt.Errorf("failed to list checks: %w", err)
	}

	if jsonOutput {
		return outputJSON(checks)
	}

	if len(checks) == 0 {
		output.InfoMessage("No checks found")
		return nil
	}

	// Create table
	table := output.NewTable([]string{"TIME", "STATUS", "RESPONSE", "SUCCESS"})
	table.Render()

	// Add rows
	for _, check := range checks {
		statusCode := fmt.Sprintf("%d", check.StatusCode)
		responseTime := fmt.Sprintf("%.2fms", check.ResponseTime)

		success := output.Green("✓")
		if !check.Success {
			success = output.Red("✗")
		}

		table.Append([]string{
			check.CreatedAt,
			statusCode,
			responseTime,
			success,
		})
	}

	table.Flush()
	fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d check(s)", len(checks))))
	return nil
}

func listJobPings(client *api.Client, jobID string, jsonOutput bool) error {
	// Resolve short ID to full ID
	fullID, err := resolveJobID(client, jobID)
	if err != nil {
		return err
	}

	var s *spinner.Spinner
	if !jsonOutput {
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
	}

	pings, err := client.ListJobPings(fullID)

	if s != nil {
		s.Stop()
	}

	if err != nil {
		return fmt.Errorf("failed to list pings: %w", err)
	}

	if jsonOutput {
		return outputJSON(pings)
	}

	if len(pings) == 0 {
		output.InfoMessage("No pings found")
		return nil
	}

	// Create table
	table := output.NewTable([]string{"TIME", "TYPE", "DURATION"})
	table.Render()

	// Add rows
	for _, ping := range pings {
		pingType := ping.PingType
		if pingType == "" {
			pingType = "heartbeat"
		}

		duration := "-"
		if ping.Duration != nil && *ping.Duration != "" {
			// Parse duration string (in seconds) and convert to milliseconds
			if durationFloat, err := strconv.ParseFloat(*ping.Duration, 64); err == nil {
				durationMs := durationFloat * 1000
				duration = fmt.Sprintf("%.0fms", durationMs)
			} else {
				duration = *ping.Duration
			}
		}

		table.Append([]string{
			ping.CreatedAt,
			pingType,
			duration,
		})
	}

	table.Flush()
	fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d ping(s)", len(pings))))
	return nil
}

func init() {
	// Add flags to list command
	checksListCmd.Flags().StringP("monitor", "m", "", "Monitor ID to view checks for")
	checksListCmd.Flags().StringP("job", "j", "", "Job ID to view pings for")
	checksListCmd.Flags().Bool("json", false, "Output as JSON")

	// Add subcommands
	checksCmd.AddCommand(checksListCmd)

	// Add checks command to root
	rootCmd.AddCommand(checksCmd)
}
