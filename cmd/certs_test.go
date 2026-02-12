package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/scookdev/groovekit-cli/internal/api"
	"github.com/scookdev/groovekit-cli/internal/config"
	"github.com/scookdev/groovekit-cli/internal/output"
	"github.com/spf13/cobra"
)

var certsCmd = &cobra.Command{
	Use:   "certs",
	Short: "Manage SSL certificate monitors",
	Long:  "List, create, show, and delete SSL certificate monitors",
}

// certs list
var certsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all certs",
	Long:  "List all cert monitors for your account",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Check for --json flag first (don't show spinner for JSON output)
		jsonOutput, _ := cmd.Flags().GetBool("json")

		// Start spinner
		var s *spinner.Spinner
		if !jsonOutput {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Start()
		}

		result, err := client.ListCerts()

		// Stop spinner
		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to list certs: %w", err)
		}
		if jsonOutput {
			return outputJSON(result)
		}

		if len(result.Certs) == 0 {
			output.InfoMessage("No certs found")
			fmt.Println("\nCreate your first cert:")
			fmt.Println("  groovekit certs create --name 'Daily Backup' --interval 1440")
			return nil
		}

		// Create table
		table := output.NewTable([]string{"ID", "NAME", "INTERVAL", "STATUS", "HEALTH"})
		table.Render()

		// Add rows
		for _, cert := range result.Certs {
			status := cert.Status
			if cert.Status == "active" {
				status = output.Green(status)
			}

			health := output.Green("✓ Up")
			if cert.Down {
				health = output.Red("✗ Down")
			}

			// Truncate ID to first 8 chars (like Docker)
			shortID := cert.ID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}

			table.Append([]string{
				output.Cyan(shortID),
				cert.Name,
				output.FormatDuration(cert.Interval),
				status,
				health,
			})
		}

		table.Flush()
		fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d cert(s)", result.TotalCount)))
		return nil
	},
}

// certs show <id>
var certsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show cert details",
	Long:  "Display detailed information about a specific cert",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveCertID(client, args[0])
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

		cert, err := client.GetCert(fullID)

		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to get cert: %w", err)
		}
		if jsonOutput {
			return outputJSON(cert)
		}

		// Print cert details
		fmt.Printf("ID:            %s\n", cert.ID)
		fmt.Printf("Name:          %s\n", cert.Name)
		fmt.Printf("Status:        %s\n", cert.Status)
		fmt.Printf("Interval:      %s\n", output.FormatDuration(cert.Interval))
		fmt.Printf("Grace Period:  %s\n", output.FormatDuration(cert.GracePeriod))
		fmt.Printf("Down:          %t\n", cert.Down)

		if cert.LastPingAt != nil {
			fmt.Printf("Last Ping:     %s\n", *cert.LastPingAt)
		} else {
			fmt.Printf("Last Ping:     Never\n")
		}

		if cert.LastRunAt != nil {
			fmt.Printf("Last Run:      %s\n", *cert.LastRunAt)
		}

		fmt.Printf("\nPing URL:\n")
		fmt.Printf("  curl https://api.groovekit.io/pings/%s\n", cert.PingToken)

		if len(cert.AllowedIPs) > 0 {
			fmt.Printf("\nAllowed IPs:   %v\n", cert.AllowedIPs)
		}

		if cert.WebhookURL != "" {
			fmt.Printf("\nWebhook URL:   %s\n", cert.WebhookURL)
		}

		return nil
	},
}

// certs create
var certsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cert",
	Long:  "Create a new cert monitor",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		interval, _ := cmd.Flags().GetInt("interval")
		gracePeriod, _ := cmd.Flags().GetInt("grace-period")

		if name == "" {
			return fmt.Errorf("--name is required")
		}
		if interval <= 0 {
			return fmt.Errorf("--interval must be greater than 0")
		}

		req := &api.CreateCertRequest{
			Name:        name,
			Interval:    interval,
			GracePeriod: gracePeriod,
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		cert, err := client.CreateCert(req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to create cert: %w", err)
		}

		output.SuccessMessage("Cert created successfully\n")
		fmt.Printf("ID:           %s\n", output.Cyan(cert.ID))
		fmt.Printf("Name:         %s\n", output.Bold(cert.Name))
		fmt.Printf("Interval:     %s\n", fmt.Sprintf("%d minutes", cert.Interval))
		fmt.Printf("Grace Period: %s\n", fmt.Sprintf("%d minutes", cert.GracePeriod))
		fmt.Printf("\n%s\n", output.Bold("Ping URL:"))
		fmt.Printf("  %s\n", output.Cyan(fmt.Sprintf("curl https://api.groovekit.io/pings/%s", cert.PingToken)))

		return nil
	},
}

// certs update <id>
var certsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a cert",
	Long:  "Update an existing cert monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveCertID(client, args[0])
		if err != nil {
			return err
		}

		// Build update request with only provided flags
		req := &api.UpdateCertRequest{}
		hasUpdates := false

		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			req.Name = &name
			hasUpdates = true
		}

		if cmd.Flags().Changed("interval") {
			interval, _ := cmd.Flags().GetInt("interval")
			req.Interval = &interval
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

		if cmd.Flags().Changed("webhook-url") {
			webhookURL, _ := cmd.Flags().GetString("webhook-url")
			req.WebhookURL = &webhookURL
			hasUpdates = true
		}

		if cmd.Flags().Changed("webhook-secret") {
			webhookSecret, _ := cmd.Flags().GetString("webhook-secret")
			req.WebhookSecret = &webhookSecret
			hasUpdates = true
		}

		if !hasUpdates {
			return fmt.Errorf("no fields to update. Use --name, --interval, --grace-period, --status, --webhook-url, or --webhook-secret")
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		cert, err := client.UpdateCert(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to update cert: %w", err)
		}

		output.SuccessMessage("Cert updated successfully\n")
		fmt.Printf("ID:           %s\n", output.Cyan(cert.ID))
		fmt.Printf("Name:         %s\n", output.Bold(cert.Name))
		fmt.Printf("Interval:     %s\n", output.FormatDuration(cert.Interval))
		fmt.Printf("Grace Period: %s\n", output.FormatDuration(cert.GracePeriod))
		fmt.Printf("Status:       %s\n", cert.Status)

		return nil
	},
}

// certs pause <id>
var certsPauseCmd = &cobra.Command{
	Use:   "pause <id>",
	Short: "Pause a cert",
	Long:  "Pause a cert monitor (sets status to paused)",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveCertID(client, args[0])
		if err != nil {
			return err
		}

		// Update status to paused
		status := "paused"
		req := &api.UpdateCertRequest{Status: &status}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		_, err = client.UpdateCert(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to pause cert: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Cert %s paused successfully", args[0]))
		return nil
	},
}

// certs resume <id>
var certsResumeCmd = &cobra.Command{
	Use:   "resume <id>",
	Short: "Resume a cert",
	Long:  "Resume a paused cert monitor (sets status to active)",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveCertID(client, args[0])
		if err != nil {
			return err
		}

		// Update status to active
		status := "active"
		req := &api.UpdateCertRequest{Status: &status}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		_, err = client.UpdateCert(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to resume cert: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Cert %s resumed successfully", args[0]))
		return nil
	},
}

// certs incidents <id>
var certsIncidentsCmd = &cobra.Command{
	Use:   "incidents <id>",
	Short: "Show incident history",
	Long:  "Display incident history (downtime periods) for a cert",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveCertID(client, args[0])
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

		incidents, err := client.ListCertIncidents(fullID)

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
			output.InfoMessage("No incidents found - this cert has been running smoothly!")
			return nil
		}

		// Create table
		table := output.NewTable([]string{"STARTED", "ENDED", "DURATION", "STATUS"})
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

			table.Append([]string{
				incident.StartedAt,
				ended,
				duration,
				status,
			})
		}

		table.Flush()
		fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d incident(s)", len(incidents))))
		return nil
	},
}

// certs delete <id>
var certsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a cert",
	Long:  "Delete a cert monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveCertID(client, args[0])
		if err != nil {
			return err
		}

		// Confirm deletion
		confirm, _ := cmd.Flags().GetBool("force")
		if !confirm {
			fmt.Printf("Are you sure you want to delete cert %s? (y/N): ", args[0])
			var response string
			_, _ = fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cancelled")
				return nil
			}
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		err = client.DeleteCert(fullID)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to delete cert: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Cert %s deleted successfully", args[0]))
		return nil
	},
}

// Helper function to get authenticated client
func getAuthenticatedClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.IsAuthenticated() {
		return nil, fmt.Errorf("not logged in. Run 'groovekit auth login' first")
	}

	return api.NewClient(cfg), nil
}

// Helper function to truncate strings
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Helper function to output JSON
func outputJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// Helper function to resolve a short ID to a full ID
func resolveCertID(client *api.Client, shortID string) (string, error) {
	// If it looks like a full UUID, use it as-is
	if len(shortID) >= 32 {
		return shortID, nil
	}

	// Otherwise, fetch all certs and match by prefix
	result, err := client.ListCerts()
	if err != nil {
		return "", fmt.Errorf("failed to list certs: %w", err)
	}

	var matches []string
	for _, cert := range result.Certs {
		if len(cert.ID) >= len(shortID) && cert.ID[:len(shortID)] == shortID {
			matches = append(matches, cert.ID)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no cert found with ID prefix '%s'", shortID)
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("ambiguous ID prefix '%s' matches multiple certs", shortID)
	}

	return matches[0], nil
}

// Helper function to format incident duration (seconds to human readable)
func formatIncidentDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.0fs", seconds)
	}
	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%.0fm", minutes)
	}
	hours := minutes / 60
	if hours < 24 {
		return fmt.Sprintf("%.1fh", hours)
	}
	days := hours / 24
	return fmt.Sprintf("%.1fd", days)
}

func init() {
	// Add flags to list command
	certsListCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to show command
	certsShowCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to create command
	certsCreateCmd.Flags().String("name", "", "Cert name (required)")
	certsCreateCmd.Flags().Int("interval", 0, "Check interval in minutes (required)")
	certsCreateCmd.Flags().Int("grace-period", 5, "Grace period in minutes")
	_ = certsCreateCmd.MarkFlagRequired("name")
	_ = certsCreateCmd.MarkFlagRequired("interval")

	// Add flags to update command
	certsUpdateCmd.Flags().String("name", "", "Cert name")
	certsUpdateCmd.Flags().Int("interval", 0, "Check interval in minutes")
	certsUpdateCmd.Flags().Int("grace-period", 0, "Grace period in minutes")
	certsUpdateCmd.Flags().String("status", "", "Cert status (active, inactive, paused)")
	certsUpdateCmd.Flags().String("webhook-url", "", "Webhook URL")
	certsUpdateCmd.Flags().String("webhook-secret", "", "Webhook secret")

	// Add flags to incidents command
	certsIncidentsCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to delete command
	certsDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	// Add subcommands
	certsCmd.AddCommand(certsListCmd)
	certsCmd.AddCommand(certsShowCmd)
	certsCmd.AddCommand(certsCreateCmd)
	certsCmd.AddCommand(certsUpdateCmd)
	certsCmd.AddCommand(certsPauseCmd)
	certsCmd.AddCommand(certsResumeCmd)
	certsCmd.AddCommand(certsIncidentsCmd)
	certsCmd.AddCommand(certsDeleteCmd)

	// Add certs command to root
	rootCmd.AddCommand(certsCmd)
}
