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
	Short: "Manage SSL certificate monitors",
	Long:  "List, create, show, update, and delete SSL certificate monitors",
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

		if len(result.SslMonitors) == 0 {
			output.InfoMessage("No SSL certificate monitors found")
			fmt.Println("\nCreate your first SSL certificate monitor:")
			fmt.Println("  groovekit certs create --name 'example.com SSL' --domain example.com")
			return nil
		}

		// Create table
		table := output.NewTable([]string{"ID", "NAME", "DOMAIN", "PORT", "DAYS LEFT", "STATUS"})
		table.Render()

		// Add rows
		for _, cert := range result.SslMonitors {
			status := cert.Status
			if cert.Status == "active" {
				status = output.Green(status)
			}

			// Truncate ID to first 8 chars (like Docker)
			shortID := cert.ID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}

			// Format days until expiration with color coding
			daysLeft := fmt.Sprintf("%d", cert.DaysUntilExpiration)
			if cert.DaysUntilExpiration <= cert.CriticalThreshold {
				daysLeft = output.Red(daysLeft)
			} else if cert.DaysUntilExpiration <= cert.UrgentThreshold {
				daysLeft = output.Yellow(daysLeft)
			} else if cert.DaysUntilExpiration <= cert.WarningThreshold {
				daysLeft = output.Yellow(daysLeft)
			} else {
				daysLeft = output.Green(daysLeft)
			}

			table.Append([]string{
				output.Cyan(shortID),
				cert.Name,
				cert.Domain,
				cert.Port,
				daysLeft,
				status,
			})
		}

		table.Flush()
		fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d SSL certificate monitor(s)", len(result.SslMonitors))))
		return nil
	},
}

// certs show <id>
var certsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show cert details",
	Long:  "Display detailed information about a specific SSL certificate monitor",
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
		fmt.Printf("ID:                       %s\n", output.Cyan(cert.ID))
		fmt.Printf("Name:                     %s\n", output.Bold(cert.Name))
		fmt.Printf("Domain:                   %s\n", cert.Domain)
		fmt.Printf("Port:                     %s\n", cert.Port)
		fmt.Printf("Status:                   %s\n", cert.Status)
		fmt.Printf("Check Interval:           %s\n", output.FormatDuration(cert.Interval))
		fmt.Printf("Grace Period:             %s\n", output.FormatDuration(cert.GracePeriod))
		fmt.Printf("Warning Threshold:        %d days\n", cert.WarningThreshold)
		fmt.Printf("Urgent Threshold:         %d days\n", cert.UrgentThreshold)
		fmt.Printf("Critical Threshold:       %d days\n", cert.CriticalThreshold)
		fmt.Printf("Days Until Expiration:    %d\n", cert.DaysUntilExpiration)
		fmt.Printf("Certificate Expires At:   %s\n", cert.CertificateExpiresAt)
		fmt.Printf("Certificate Issuer:       %s\n", cert.CertificateIssuer)
		fmt.Printf("Certificate Subject:      %s\n", cert.CertificateSubject)
		fmt.Printf("Last Check At:            %s\n", cert.LastCheckAt)
		fmt.Printf("Last Successful Check:    %s\n", cert.LastSuccessfulCheckAt)
		fmt.Printf("Consecutive Failures:     %d\n", cert.ConsecutiveFailures)
		fmt.Printf("Created At:               %s\n", cert.CreatedAt)
		fmt.Printf("Updated At:               %s\n", cert.UpdatedAt)

		return nil
	},
}

// certs create
var certsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new SSL certificate monitor",
	Long:  "Create a new SSL certificate monitor",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		domain, _ := cmd.Flags().GetString("domain")
		port, _ := cmd.Flags().GetString("port")
		interval, _ := cmd.Flags().GetInt("interval")

		if name == "" {
			return fmt.Errorf("--name is required")
		}
		if domain == "" {
			return fmt.Errorf("--domain is required")
		}

		req := &api.CreateSslMonitorRequest{
			Name:     name,
			Domain:   domain,
			Port:     port,
			Interval: interval,
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		cert, err := client.CreateCert(req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to create SSL monitor: %w", err)
		}

		output.SuccessMessage("SSL certificate monitor created successfully\n")
		fmt.Printf("ID:       %s\n", output.Cyan(cert.ID))
		fmt.Printf("Name:     %s\n", output.Bold(cert.Name))
		fmt.Printf("Domain:   %s\n", cert.Domain)
		fmt.Printf("Port:     %s\n", cert.Port)
		fmt.Printf("Interval: %s\n", output.FormatDuration(cert.Interval))

		return nil
	},
}

// certs update <id>
var certsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an SSL certificate monitor",
	Long:  "Update an existing SSL certificate monitor",
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
		req := &api.UpdateSslMonitorRequest{}
		hasUpdates := false

		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			req.Name = &name
			hasUpdates = true
		}

		if cmd.Flags().Changed("domain") {
			domain, _ := cmd.Flags().GetString("domain")
			req.Domain = &domain
			hasUpdates = true
		}

		if cmd.Flags().Changed("port") {
			port, _ := cmd.Flags().GetString("port")
			req.Port = &port
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

		if cmd.Flags().Changed("warning-threshold") {
			warning, _ := cmd.Flags().GetInt("warning-threshold")
			req.WarningThreshold = &warning
			hasUpdates = true
		}

		if cmd.Flags().Changed("urgent-threshold") {
			urgent, _ := cmd.Flags().GetInt("urgent-threshold")
			req.UrgentThreshold = &urgent
			hasUpdates = true
		}

		if cmd.Flags().Changed("critical-threshold") {
			critical, _ := cmd.Flags().GetInt("critical-threshold")
			req.CriticalThreshold = &critical
			hasUpdates = true
		}

		if cmd.Flags().Changed("status") {
			status, _ := cmd.Flags().GetString("status")
			req.Status = &status
			hasUpdates = true
		}

		if !hasUpdates {
			return fmt.Errorf("no fields to update. Use --name, --domain, --port, --interval, --grace-period, --warning-threshold, --urgent-threshold, --critical-threshold, or --status")
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		cert, err := client.UpdateCert(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to update SSL monitor: %w", err)
		}

		output.SuccessMessage("SSL certificate monitor updated successfully\n")
		fmt.Printf("ID:       %s\n", output.Cyan(cert.ID))
		fmt.Printf("Name:     %s\n", output.Bold(cert.Name))
		fmt.Printf("Domain:   %s\n", cert.Domain)
		fmt.Printf("Port:     %s\n", cert.Port)
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
		req := &api.UpdateSslMonitorRequest{Status: &status}

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
		req := &api.UpdateSslMonitorRequest{Status: &status}

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
		return "", fmt.Errorf("failed to list SSL monitors: %w", err)
	}

	var matches []string
	for _, cert := range result.SslMonitors {
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
	certsCreateCmd.Flags().String("name", "", "SSL monitor name (required)")
	certsCreateCmd.Flags().String("domain", "", "Domain to monitor (required)")
	certsCreateCmd.Flags().String("port", "443", "Port number")
	certsCreateCmd.Flags().Int("interval", 1440, "Check interval in minutes (default: daily)")
	_ = certsCreateCmd.MarkFlagRequired("name")
	_ = certsCreateCmd.MarkFlagRequired("domain")

	// Add flags to update command
	certsUpdateCmd.Flags().String("name", "", "SSL monitor name")
	certsUpdateCmd.Flags().String("domain", "", "Domain to monitor")
	certsUpdateCmd.Flags().String("port", "", "Port number")
	certsUpdateCmd.Flags().Int("interval", 0, "Check interval in minutes")
	certsUpdateCmd.Flags().Int("grace-period", 0, "Grace period in minutes")
	certsUpdateCmd.Flags().Int("warning-threshold", 0, "Warning threshold in days")
	certsUpdateCmd.Flags().Int("urgent-threshold", 0, "Urgent threshold in days")
	certsUpdateCmd.Flags().Int("critical-threshold", 0, "Critical threshold in days")
	certsUpdateCmd.Flags().String("status", "", "Monitor status (active, inactive, paused)")

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
