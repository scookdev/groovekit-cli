package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/scookdev/groovekit-cli/internal/api"
	"github.com/scookdev/groovekit-cli/internal/output"
	"github.com/spf13/cobra"
)

var certsCmd = &cobra.Command{
	Use:   "certs",
	Short: "Manage SSL certificate certs",
	Long:  "List, create, show, and delete SSL certificate certs",
}

// certs list
var certsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all certs",
	Long:  "List all API endpoint certs for your account",
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

		result, err := client.ListCerts()

		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to list certs: %w", err)
		}
		if jsonOutput {
			return outputJSON(result)
		}

		if len(result.APICerts) == 0 {
			output.InfoMessage("No certs found")
			fmt.Println("\nCreate your first cert:")
			fmt.Println("  groovekit certs create --name 'Production API' --url https://api.example.com/health --interval 60")
			return nil
		}

		// Create table
		table := output.NewTable([]string{"ID", "NAME", "URL", "INTERVAL", "STATUS", "HEALTH"})
		table.Render()

		// Add rows
		for _, cert := range result.APICerts {
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
				truncate(cert.URL, 40),
				output.FormatDuration(cert.Interval),
				status,
				health,
			})
		}

		table.Flush()
		fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d cert(s)", len(result.APICerts))))
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
		fmt.Printf("ID:               %s\n", output.Cyan(cert.ID))
		fmt.Printf("Name:             %s\n", output.Bold(cert.Name))
		fmt.Printf("URL:              %s\n", cert.URL)
		fmt.Printf("HTTP Method:      %s\n", cert.HTTPMethod)
		fmt.Printf("Status:           %s\n", cert.Status)
		fmt.Printf("Interval:         %s\n", output.FormatDuration(cert.Interval))
		fmt.Printf("Timeout:          %d seconds\n", cert.Timeout)
		fmt.Printf("Grace Period:     %s\n", output.FormatDuration(cert.GracePeriod))
		fmt.Printf("Down:             %t\n", cert.Down)

		if len(cert.ExpectedStatusCodes) > 0 {
			fmt.Printf("Expected Status:  %v\n", cert.ExpectedStatusCodes)
		}

		if cert.LastCheckAt != nil {
			fmt.Printf("Last Check:       %s\n", *cert.LastCheckAt)
		} else {
			fmt.Printf("Last Check:       Never\n")
		}

		if cert.UptimePercentage != nil {
			fmt.Printf("Uptime (30d):     %.2f%%\n", *cert.UptimePercentage)
		}

		if cert.AverageResponseTime != nil {
			fmt.Printf("Avg Response:     %.0fms\n", *cert.AverageResponseTime)
		}

		if len(cert.ValidateResponsePaths) > 0 {
			fmt.Printf("\nJSON Path Validation:\n")
			for _, path := range cert.ValidateResponsePaths {
				fmt.Printf("  - %s\n", path)
			}
		}

		return nil
	},
}

// certs create
var certsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cert",
	Long:  "Create a new API endpoint cert",
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

		req := &api.CreateCertRequest{
			Name:       name,
			URL:        url,
			Interval:   interval,
			HTTPMethod: method,
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		cert, err := client.CreateCert(req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to create cert: %w", err)
		}

		output.SuccessMessage("Cert created successfully\n")
		fmt.Printf("ID:          %s\n", output.Cyan(cert.ID))
		fmt.Printf("Name:        %s\n", output.Bold(cert.Name))
		fmt.Printf("URL:         %s\n", cert.URL)
		fmt.Printf("Interval:    %s\n", fmt.Sprintf("%d minutes", cert.Interval))

		return nil
	},
}

// certs update <id>
var certsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a cert",
	Long:  "Update an existing API endpoint cert",
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
		cert, err := client.UpdateCert(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to update cert: %w", err)
		}

		output.SuccessMessage("Cert updated successfully\n")
		fmt.Printf("ID:       %s\n", output.Cyan(cert.ID))
		fmt.Printf("Name:     %s\n", output.Bold(cert.Name))
		fmt.Printf("URL:      %s\n", cert.URL)
		fmt.Printf("Method:   %s\n", cert.HTTPMethod)
		fmt.Printf("Interval: %s\n", output.FormatDuration(cert.Interval))
		fmt.Printf("Status:   %s\n", cert.Status)

		return nil
	},
}

// certs pause <id>
var certsPauseCmd = &cobra.Command{
	Use:   "pause <id>",
	Short: "Pause a cert",
	Long:  "Pause an API endpoint cert (sets status to paused)",
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
	Long:  "Resume a paused API endpoint cert (sets status to active)",
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

// certs delete <id>
var certsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a cert",
	Long:  "Delete an API endpoint cert",
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

// Helper function to resolve a short cert ID to a full ID
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
	for _, cert := range result.APICerts {
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

func init() {
	// Add flags to list command
	certsListCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to show command
	certsShowCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to create command
	certsCreateCmd.Flags().String("name", "", "cert name (required)")
	certsCreateCmd.Flags().String("url", "", "URL to cert (required)")
	certsCreateCmd.Flags().Int("interval", 60, "Check interval in minutes")
	certsCreateCmd.Flags().String("method", "GET", "HTTP method")
	_ = certsCreateCmd.MarkFlagRequired("name")
	_ = certsCreateCmd.MarkFlagRequired("url")

	// Add flags to update command
	certsUpdateCmd.Flags().String("name", "", "Cert name")
	certsUpdateCmd.Flags().String("url", "", "URL to cert")
	certsUpdateCmd.Flags().String("http-method", "", "HTTP method (GET, POST, etc)")
	certsUpdateCmd.Flags().Int("interval", 0, "Check interval in minutes")
	certsUpdateCmd.Flags().Int("timeout", 0, "Request timeout in seconds")
	certsUpdateCmd.Flags().Int("grace-period", 0, "Grace period in minutes")
	certsUpdateCmd.Flags().String("status", "", "Cert status (active, inactive, paused)")
	certsUpdateCmd.Flags().IntSlice("expected-status-codes", nil, "Expected HTTP status codes (comma-separated)")

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
