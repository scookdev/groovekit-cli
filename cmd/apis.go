package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/scookdev/groovekit-cli/internal/api"
	"github.com/scookdev/groovekit-cli/internal/output"
	"github.com/spf13/cobra"
)

var apisCmd = &cobra.Command{
	Use:   "apis",
	Short: "Manage API endpoint monitors",
	Long:  "List, create, show, and delete API endpoint monitors",
}

// apis list
var apisListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all api monitors",
	Long:  "List all API endpoint monitors for your account",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Check for --json flag first
		jsonOutput, _ := cmd.Flags().GetBool("json")

		var s *spinner.Spinner
		if !jsonOutput {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Start()
		}

		result, err := client.ListApis()

		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to list API monitors: %w", err)
		}
		if jsonOutput {
			return outputJSON(result)
		}

		if len(result.APIMonitors) == 0 {
			output.InfoMessage("No API monitors found")
			fmt.Println("\nCreate your first API monitor:")
			fmt.Println("  groovekit apis create --name 'Production API' --url https://api.example.com/health --interval 60")
			return nil
		}

		// Create table
		table := output.NewTable([]string{"ID", "NAME", "URL", "INTERVAL", "STATUS", "HEALTH"})
		table.Render()

		// Add rows
		for _, monitor := range result.APIMonitors {
			status := monitor.Status
			if monitor.Status == "active" {
				status = output.Green(status)
			}

			health := output.Green("✓ Up")
			if monitor.Down {
				health = output.Red("✗ Down")
			}

			// Truncate ID to first 8 chars (like Docker)
			shortID := monitor.ID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}

			table.Append([]string{
				output.Cyan(shortID),
				monitor.Name,
				truncate(monitor.URL, 40),
				output.FormatDuration(monitor.Interval),
				status,
				health,
			})
		}

		table.Flush()
		fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d API monitor(s)", len(result.APIMonitors))))
		return nil
	},
}

// apis show <id>
var apisShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show API monitor details",
	Long:  "Display detailed information about a specific API monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveMonitorID(client, args[0])
		if err != nil {
			return err
		}

		// Check for --json flag first
		jsonOutput, _ := cmd.Flags().GetBool("json")

		var s *spinner.Spinner
		if !jsonOutput {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Start()
		}

		monitor, err := client.GetApi(fullID)

		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to get API monitor: %w", err)
		}
		if jsonOutput {
			return outputJSON(monitor)
		}

		// Print monitor details
		fmt.Printf("ID:               %s\n", output.Cyan(monitor.ID))
		fmt.Printf("Name:             %s\n", output.Bold(monitor.Name))
		fmt.Printf("URL:              %s\n", monitor.URL)
		fmt.Printf("HTTP Method:      %s\n", monitor.HTTPMethod)
		fmt.Printf("Status:           %s\n", monitor.Status)
		fmt.Printf("Interval:         %s\n", output.FormatDuration(monitor.Interval))
		fmt.Printf("Timeout:          %d seconds\n", monitor.Timeout)
		fmt.Printf("Grace Period:     %s\n", output.FormatDuration(monitor.GracePeriod))
		fmt.Printf("Down:             %t\n", monitor.Down)

		if len(monitor.ExpectedStatusCodes) > 0 {
			fmt.Printf("Expected Status:  %v\n", monitor.ExpectedStatusCodes)
		}

		if monitor.LastCheckAt != nil {
			fmt.Printf("Last Check:       %s\n", *monitor.LastCheckAt)
		} else {
			fmt.Printf("Last Check:       Never\n")
		}

		if monitor.UptimePercentage != nil {
			fmt.Printf("Uptime (30d):     %.2f%%\n", *monitor.UptimePercentage)
		}

		if monitor.AverageResponseTime != nil {
			fmt.Printf("Avg Response:     %.0fms\n", *monitor.AverageResponseTime)
		}

		if len(monitor.ValidateResponsePaths) > 0 {
			fmt.Printf("\nJSON Path Validation:\n")
			for _, path := range monitor.ValidateResponsePaths {
				fmt.Printf("  - %s\n", path)
			}
		}

		return nil
	},
}

// apis create
var apisCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API monitor",
	Long:  "Create a new API endpoint monitor",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		url, _ := cmd.Flags().GetString("url")
		interval, _ := cmd.Flags().GetInt("interval")
		method, _ := cmd.Flags().GetString("method")

		if name == "" {
			return fmt.Errorf("--name is required")
		}
		if url == "" {
			return fmt.Errorf("--url is required")
		}
		if interval <= 0 {
			return fmt.Errorf("--interval must be greater than 0")
		}

		req := &api.CreateApiRequest{
			Name:       name,
			URL:        url,
			Interval:   interval,
			HTTPMethod: method,
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		monitor, err := client.CreateApi(req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to create API monitor: %w", err)
		}

		output.SuccessMessage("API monitor created successfully\n")
		fmt.Printf("ID:          %s\n", output.Cyan(monitor.ID))
		fmt.Printf("Name:        %s\n", output.Bold(monitor.Name))
		fmt.Printf("URL:         %s\n", monitor.URL)
		fmt.Printf("Interval:    %s\n", fmt.Sprintf("%d minutes", monitor.Interval))

		return nil
	},
}

// apis update <id>
var apisUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an API monitor",
	Long:  "Update an existing API endpoint monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveMonitorID(client, args[0])
		if err != nil {
			return err
		}

		// Build update request with only provided flags
		req := &api.UpdateApiRequest{}
		hasUpdates := false

		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			req.Name = &name
			hasUpdates = true
		}

		if cmd.Flags().Changed("url") {
			url, _ := cmd.Flags().GetString("url")
			req.URL = &url
			hasUpdates = true
		}

		if cmd.Flags().Changed("http-method") {
			method, _ := cmd.Flags().GetString("http-method")
			req.HTTPMethod = &method
			hasUpdates = true
		}

		if cmd.Flags().Changed("interval") {
			interval, _ := cmd.Flags().GetInt("interval")
			req.Interval = &interval
			hasUpdates = true
		}

		if cmd.Flags().Changed("timeout") {
			timeout, _ := cmd.Flags().GetInt("timeout")
			req.Timeout = &timeout
			hasUpdates = true
		}

		if cmd.Flags().Changed("grace-period") {
			gracePeriod, _ := cmd.Flags().GetInt("grace-period")
			req.GracePeriod = &gracePeriod
			hasUpdates = true
		}

		if cmd.Flags().Changed("status") {
			status, _ := cmd.Flags().GetString("status")
			req.Status = &status
			hasUpdates = true
		}

		if cmd.Flags().Changed("expected-status-codes") {
			codes, _ := cmd.Flags().GetIntSlice("expected-status-codes")
			req.ExpectedStatusCodes = &codes
			hasUpdates = true
		}

		if !hasUpdates {
			return fmt.Errorf("no fields to update. Use --name, --url, --http-method, --interval, --timeout, --grace-period, --status, or --expected-status-codes")
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		monitor, err := client.UpdateApi(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to update API monitor: %w", err)
		}

		output.SuccessMessage("API monitor updated successfully\n")
		fmt.Printf("ID:       %s\n", output.Cyan(monitor.ID))
		fmt.Printf("Name:     %s\n", output.Bold(monitor.Name))
		fmt.Printf("URL:      %s\n", monitor.URL)
		fmt.Printf("Method:   %s\n", monitor.HTTPMethod)
		fmt.Printf("Interval: %s\n", output.FormatDuration(monitor.Interval))
		fmt.Printf("Status:   %s\n", monitor.Status)

		return nil
	},
}

// apis pause <id>
var apisPauseCmd = &cobra.Command{
	Use:   "pause <id>",
	Short: "Pause an API monitor",
	Long:  "Pause an API endpoint monitor (sets status to paused)",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveMonitorID(client, args[0])
		if err != nil {
			return err
		}

		// Update status to paused
		status := "paused"
		req := &api.UpdateApiRequest{Status: &status}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		_, err = client.UpdateApi(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to pause API monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("API monitor %s paused successfully", args[0]))
		return nil
	},
}

// apis resume <id>
var apisResumeCmd = &cobra.Command{
	Use:   "resume <id>",
	Short: "Resume an API monitor",
	Long:  "Resume a paused API endpoint monitor (sets status to active)",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveMonitorID(client, args[0])
		if err != nil {
			return err
		}

		// Update status to active
		status := "active"
		req := &api.UpdateApiRequest{Status: &status}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		_, err = client.UpdateApi(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to resume API monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("API monitor %s resumed successfully", args[0]))
		return nil
	},
}

// apis incidents <id>
var apisIncidentsCmd = &cobra.Command{
	Use:   "incidents <id>",
	Short: "Show incident history",
	Long:  "Display incident history (downtime periods) for an API monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveMonitorID(client, args[0])
		if err != nil {
			return err
		}

		// Check for --json flag
		jsonOutput, _ := cmd.Flags().GetBool("json")

		var s *spinner.Spinner
		if !jsonOutput {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Start()
		}

		incidents, err := client.ListApiIncidents(fullID)

		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to get incidents: %w", err)
		}

		if jsonOutput {
			return outputJSON(incidents)
		}

		if len(incidents) == 0 {
			output.InfoMessage("No incidents found - this API monitor has been running smoothly!")
			return nil
		}

		// Create table
		table := output.NewTable([]string{"STARTED", "ENDED", "DURATION", "STATUS", "ERROR"})
		table.Render()

		// Add rows
		for _, incident := range incidents {
			status := output.Red("Ongoing")
			ended := output.Yellow("Still down")

			if incident.EndedAt != nil {
				status = output.Green("Recovered")
				ended = *incident.EndedAt
			}

			// Format duration
			duration := formatIncidentDuration(incident.Duration)

			// Truncate error message
			errorMsg := "-"
			if incident.ErrorMessage != nil {
				errorMsg = truncate(*incident.ErrorMessage, 40)
			}

			table.Append([]string{
				incident.StartedAt,
				ended,
				duration,
				status,
				errorMsg,
			})
		}

		table.Flush()
		fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d incident(s)", len(incidents))))
		return nil
	},
}

// apis delete <id>
var apisDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an API monitor",
	Long:  "Delete an API endpoint monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveMonitorID(client, args[0])
		if err != nil {
			return err
		}

		// Confirm deletion
		confirm, _ := cmd.Flags().GetBool("force")
		if !confirm {
			fmt.Printf("Are you sure you want to delete API monitor %s? (y/N): ", args[0])
			var response string
			_, _ = fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cancelled")
				return nil
			}
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		err = client.DeleteApi(fullID)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to delete API monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("API monitor %s deleted successfully", args[0]))
		return nil
	},
}

// Helper function to resolve a short monitor ID to a full ID
func resolveMonitorID(client *api.Client, shortID string) (string, error) {
	// If it looks like a full UUID, use it as-is
	if len(shortID) >= 32 {
		return shortID, nil
	}

	// Otherwise, fetch all monitors and match by prefix
	result, err := client.ListApis()
	if err != nil {
		return "", fmt.Errorf("failed to list API monitors: %w", err)
	}

	var matches []string
	for _, monitor := range result.APIMonitors {
		if len(monitor.ID) >= len(shortID) && monitor.ID[:len(shortID)] == shortID {
			matches = append(matches, monitor.ID)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no API monitor found with ID prefix '%s'", shortID)
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("ambiguous ID prefix '%s' matches multiple API monitors", shortID)
	}

	return matches[0], nil
}

func init() {
	// Add flags to list command
	apisListCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to show command
	apisShowCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to create command
	apisCreateCmd.Flags().String("name", "", "Monitor name (required)")
	apisCreateCmd.Flags().String("url", "", "URL to monitor (required)")
	apisCreateCmd.Flags().Int("interval", 60, "Check interval in minutes")
	apisCreateCmd.Flags().String("method", "GET", "HTTP method")
	_ = apisCreateCmd.MarkFlagRequired("name")
	_ = apisCreateCmd.MarkFlagRequired("url")

	// Add flags to update command
	apisUpdateCmd.Flags().String("name", "", "Monitor name")
	apisUpdateCmd.Flags().String("url", "", "URL to monitor")
	apisUpdateCmd.Flags().String("http-method", "", "HTTP method (GET, POST, etc)")
	apisUpdateCmd.Flags().Int("interval", 0, "Check interval in minutes")
	apisUpdateCmd.Flags().Int("timeout", 0, "Request timeout in seconds")
	apisUpdateCmd.Flags().Int("grace-period", 0, "Grace period in minutes")
	apisUpdateCmd.Flags().String("status", "", "Monitor status (active, inactive, paused)")
	apisUpdateCmd.Flags().IntSlice("expected-status-codes", nil, "Expected HTTP status codes (comma-separated)")

	// Add flags to incidents command
	apisIncidentsCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to delete command
	apisDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	// Add subcommands
	apisCmd.AddCommand(apisListCmd)
	apisCmd.AddCommand(apisShowCmd)
	apisCmd.AddCommand(apisCreateCmd)
	apisCmd.AddCommand(apisUpdateCmd)
	apisCmd.AddCommand(apisPauseCmd)
	apisCmd.AddCommand(apisResumeCmd)
	apisCmd.AddCommand(apisIncidentsCmd)
	apisCmd.AddCommand(apisDeleteCmd)

	// Add apis command to root
	rootCmd.AddCommand(apisCmd)
}
