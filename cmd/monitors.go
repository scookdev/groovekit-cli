package cmd

import (
	"fmt"

	"github.com/scookdev/groovekit-cli/internal/api"
	"github.com/scookdev/groovekit-cli/internal/output"
	"github.com/spf13/cobra"
)

var monitorsCmd = &cobra.Command{
	Use:   "monitors",
	Short: "Manage API endpoint monitors",
	Long:  "List, create, show, and delete API endpoint monitors",
}

// monitors list
var monitorsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all monitors",
	Long:  "List all API endpoint monitors for your account",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		result, err := client.ListMonitors()
		if err != nil {
			return fmt.Errorf("failed to list monitors: %w", err)
		}

		// Check for --json flag
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			return outputJSON(result)
		}

		if len(result.APIMonitors) == 0 {
			output.InfoMessage("No monitors found")
			fmt.Println("\nCreate your first monitor:")
			fmt.Println("  groovekit monitors create --name 'Production API' --url https://api.example.com/health --interval 60")
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
		fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d monitor(s)", len(result.APIMonitors))))
		return nil
	},
}

// monitors show <id>
var monitorsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show monitor details",
	Long:  "Display detailed information about a specific monitor",
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

		monitor, err := client.GetMonitor(fullID)
		if err != nil {
			return fmt.Errorf("failed to get monitor: %w", err)
		}

		// Check for --json flag
		jsonOutput, _ := cmd.Flags().GetBool("json")
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

// monitors create
var monitorsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new monitor",
	Long:  "Create a new API endpoint monitor",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		req := &api.CreateMonitorRequest{
			Name:       name,
			URL:        url,
			Interval:   interval,
			HTTPMethod: method,
		}

		monitor, err := client.CreateMonitor(req)
		if err != nil {
			return fmt.Errorf("failed to create monitor: %w", err)
		}

		output.SuccessMessage("Monitor created successfully\n")
		fmt.Printf("ID:          %s\n", output.Cyan(monitor.ID))
		fmt.Printf("Name:        %s\n", output.Bold(monitor.Name))
		fmt.Printf("URL:         %s\n", monitor.URL)
		fmt.Printf("Interval:    %s\n", fmt.Sprintf("%d minutes", monitor.Interval))

		return nil
	},
}

// monitors update <id>
var monitorsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a monitor",
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
		req := &api.UpdateMonitorRequest{}
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

		monitor, err := client.UpdateMonitor(fullID, req)
		if err != nil {
			return fmt.Errorf("failed to update monitor: %w", err)
		}

		output.SuccessMessage("Monitor updated successfully\n")
		fmt.Printf("ID:       %s\n", output.Cyan(monitor.ID))
		fmt.Printf("Name:     %s\n", output.Bold(monitor.Name))
		fmt.Printf("URL:      %s\n", monitor.URL)
		fmt.Printf("Method:   %s\n", monitor.HTTPMethod)
		fmt.Printf("Interval: %s\n", output.FormatDuration(monitor.Interval))
		fmt.Printf("Status:   %s\n", monitor.Status)

		return nil
	},
}

// monitors pause <id>
var monitorsPauseCmd = &cobra.Command{
	Use:   "pause <id>",
	Short: "Pause a monitor",
	Long:  "Pause an API endpoint monitor (sets status to paused)",
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

		// Update status to paused
		status := "paused"
		req := &api.UpdateMonitorRequest{Status: &status}

		if _, err := client.UpdateMonitor(fullID, req); err != nil {
			return fmt.Errorf("failed to pause monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Monitor %s paused successfully", args[0]))
		return nil
	},
}

// monitors resume <id>
var monitorsResumeCmd = &cobra.Command{
	Use:   "resume <id>",
	Short: "Resume a monitor",
	Long:  "Resume a paused API endpoint monitor (sets status to active)",
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

		// Update status to active
		status := "active"
		req := &api.UpdateMonitorRequest{Status: &status}

		if _, err := client.UpdateMonitor(fullID, req); err != nil {
			return fmt.Errorf("failed to resume monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Monitor %s resumed successfully", args[0]))
		return nil
	},
}

// monitors incidents <id>
var monitorsIncidentsCmd = &cobra.Command{
	Use:   "incidents <id>",
	Short: "Show incident history",
	Long:  "Display incident history (downtime periods) for a monitor",
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

		incidents, err := client.ListMonitorIncidents(fullID)
		if err != nil {
			return fmt.Errorf("failed to get incidents: %w", err)
		}

		if jsonOutput {
			return outputJSON(incidents)
		}

		if len(incidents) == 0 {
			output.InfoMessage("No incidents found - this monitor has been running smoothly!")
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

// monitors delete <id>
var monitorsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a monitor",
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
			fmt.Printf("Are you sure you want to delete monitor %s? (y/N): ", args[0])
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cancelled")
				return nil
			}
		}

		if err := client.DeleteMonitor(fullID); err != nil {
			return fmt.Errorf("failed to delete monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Monitor %s deleted successfully", args[0]))
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
	result, err := client.ListMonitors()
	if err != nil {
		return "", fmt.Errorf("failed to list monitors: %w", err)
	}

	var matches []string
	for _, monitor := range result.APIMonitors {
		if len(monitor.ID) >= len(shortID) && monitor.ID[:len(shortID)] == shortID {
			matches = append(matches, monitor.ID)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no monitor found with ID prefix '%s'", shortID)
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("ambiguous ID prefix '%s' matches multiple monitors", shortID)
	}

	return matches[0], nil
}

func init() {
	// Add flags to list command
	monitorsListCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to show command
	monitorsShowCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to create command
	monitorsCreateCmd.Flags().String("name", "", "Monitor name (required)")
	monitorsCreateCmd.Flags().String("url", "", "URL to monitor (required)")
	monitorsCreateCmd.Flags().Int("interval", 60, "Check interval in minutes")
	monitorsCreateCmd.Flags().String("method", "GET", "HTTP method")
	monitorsCreateCmd.MarkFlagRequired("name")
	monitorsCreateCmd.MarkFlagRequired("url")

	// Add flags to update command
	monitorsUpdateCmd.Flags().String("name", "", "Monitor name")
	monitorsUpdateCmd.Flags().String("url", "", "URL to monitor")
	monitorsUpdateCmd.Flags().String("http-method", "", "HTTP method (GET, POST, etc)")
	monitorsUpdateCmd.Flags().Int("interval", 0, "Check interval in minutes")
	monitorsUpdateCmd.Flags().Int("timeout", 0, "Request timeout in seconds")
	monitorsUpdateCmd.Flags().Int("grace-period", 0, "Grace period in minutes")
	monitorsUpdateCmd.Flags().String("status", "", "Monitor status (active, inactive, paused)")
	monitorsUpdateCmd.Flags().IntSlice("expected-status-codes", nil, "Expected HTTP status codes (comma-separated)")

	// Add flags to incidents command
	monitorsIncidentsCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to delete command
	monitorsDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	// Add subcommands
	monitorsCmd.AddCommand(monitorsListCmd)
	monitorsCmd.AddCommand(monitorsShowCmd)
	monitorsCmd.AddCommand(monitorsCreateCmd)
	monitorsCmd.AddCommand(monitorsUpdateCmd)
	monitorsCmd.AddCommand(monitorsPauseCmd)
	monitorsCmd.AddCommand(monitorsResumeCmd)
	monitorsCmd.AddCommand(monitorsIncidentsCmd)
	monitorsCmd.AddCommand(monitorsDeleteCmd)

	// Add monitors command to root
	rootCmd.AddCommand(monitorsCmd)
}
