package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/scookdev/groovekit-cli/internal/api"
	"github.com/scookdev/groovekit-cli/internal/output"
	"github.com/spf13/cobra"
)

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Manage domain expiration monitors",
	Long:  "List, create, show, update, and delete domain expiration monitors",
}

// domains list
var domainsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all domains",
	Long:  "List all domain expiration monitors for your account",
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

		result, err := client.ListDomains()

		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to list domains: %w", err)
		}
		if jsonOutput {
			return outputJSON(result)
		}

		if len(result.DomainMonitors) == 0 {
			output.InfoMessage("No domain monitors found")
			fmt.Println("\nCreate your first domain monitor:")
			fmt.Println("  groovekit domains create --name 'example.com' --domain example.com")
			return nil
		}

		// Create table
		table := output.NewTable([]string{"ID", "NAME", "DOMAIN", "DAYS LEFT", "REGISTRAR", "STATUS"})
		table.Render()

		// Add rows
		for _, domain := range result.DomainMonitors {
			status := domain.Status
			if domain.Status == "active" {
				status = output.Green(status)
			}

			// Truncate ID to first 8 chars (like Docker)
			shortID := domain.ID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}

			// Format days until expiration with color coding
			daysLeft := fmt.Sprintf("%d", domain.DaysUntilExpiration)
			if domain.DaysUntilExpiration <= domain.CriticalThreshold {
				daysLeft = output.Red(daysLeft)
			} else if domain.DaysUntilExpiration <= domain.UrgentThreshold {
				daysLeft = output.Yellow(daysLeft)
			} else if domain.DaysUntilExpiration <= domain.WarningThreshold {
				daysLeft = output.Yellow(daysLeft)
			} else {
				daysLeft = output.Green(daysLeft)
			}

			table.Append([]string{
				output.Cyan(shortID),
				domain.Name,
				domain.Domain,
				daysLeft,
				truncate(domain.Registrar, 20),
				status,
			})
		}

		table.Flush()
		fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d domain monitor(s)", len(result.DomainMonitors))))
		return nil
	},
}

// domains show <id>
var domainsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show domain details",
	Long:  "Display detailed information about a specific domain monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDomainID(client, args[0])
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

		domain, err := client.GetDomain(fullID)

		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to get domain: %w", err)
		}
		if jsonOutput {
			return outputJSON(domain)
		}

		// Print domain details
		fmt.Printf("ID:                       %s\n", output.Cyan(domain.ID))
		fmt.Printf("Name:                     %s\n", output.Bold(domain.Name))
		fmt.Printf("Domain:                   %s\n", domain.Domain)
		fmt.Printf("Status:                   %s\n", domain.Status)
		fmt.Printf("Check Interval:           %s\n", output.FormatDuration(domain.Interval))
		fmt.Printf("Grace Period:             %s\n", output.FormatDuration(domain.GracePeriod))
		fmt.Printf("Warning Threshold:        %d days\n", domain.WarningThreshold)
		fmt.Printf("Urgent Threshold:         %d days\n", domain.UrgentThreshold)
		fmt.Printf("Critical Threshold:       %d days\n", domain.CriticalThreshold)
		fmt.Printf("Days Until Expiration:    %d\n", domain.DaysUntilExpiration)
		fmt.Printf("Expires At:               %s\n", domain.ExpiresAt)
		fmt.Printf("Registrar:                %s\n", domain.Registrar)
		if domain.RegistrarURL != nil {
			fmt.Printf("Registrar URL:            %s\n", *domain.RegistrarURL)
		}
		fmt.Printf("Last Check At:            %s\n", domain.LastCheckAt)
		fmt.Printf("Last Successful Check:    %s\n", domain.LastSuccessfulCheckAt)
		fmt.Printf("Consecutive Failures:     %d\n", domain.ConsecutiveFailures)
		fmt.Printf("Created At:               %s\n", domain.CreatedAt)
		fmt.Printf("Updated At:               %s\n", domain.UpdatedAt)

		return nil
	},
}

// domains create
var domainsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new domain monitor",
	Long:  "Create a new domain expiration monitor",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		domain, _ := cmd.Flags().GetString("domain")
		interval, _ := cmd.Flags().GetInt("interval")
		gracePeriod, _ := cmd.Flags().GetInt("grace-period")
		warningThreshold, _ := cmd.Flags().GetInt("warning-threshold")
		urgentThreshold, _ := cmd.Flags().GetInt("urgent-threshold")
		criticalThreshold, _ := cmd.Flags().GetInt("critical-threshold")

		if name == "" {
			return fmt.Errorf("--name is required")
		}
		if domain == "" {
			return fmt.Errorf("--domain is required")
		}

		req := &api.CreateDomainMonitorRequest{
			Name:              name,
			Domain:            domain,
			Interval:          interval,
			GracePeriod:       gracePeriod,
			WarningThreshold:  warningThreshold,
			UrgentThreshold:   urgentThreshold,
			CriticalThreshold: criticalThreshold,
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		domainMonitor, err := client.CreateDomain(req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to create domain monitor: %w", err)
		}

		output.SuccessMessage("Domain monitor created successfully\n")
		fmt.Printf("ID:       %s\n", output.Cyan(domainMonitor.ID))
		fmt.Printf("Name:     %s\n", output.Bold(domainMonitor.Name))
		fmt.Printf("Domain:   %s\n", domainMonitor.Domain)
		fmt.Printf("Interval: %s\n", output.FormatDuration(domainMonitor.Interval))

		return nil
	},
}

// domains update <id>
var domainsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a domain monitor",
	Long:  "Update an existing domain expiration monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDomainID(client, args[0])
		if err != nil {
			return err
		}

		// Build update request with only provided flags
		req := &api.UpdateDomainMonitorRequest{}
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
			return fmt.Errorf("no fields to update. Use --name, --domain, --interval, --grace-period, --warning-threshold, --urgent-threshold, --critical-threshold, or --status")
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		domainMonitor, err := client.UpdateDomain(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to update domain monitor: %w", err)
		}

		output.SuccessMessage("Domain monitor updated successfully\n")
		fmt.Printf("ID:       %s\n", output.Cyan(domainMonitor.ID))
		fmt.Printf("Name:     %s\n", output.Bold(domainMonitor.Name))
		fmt.Printf("Domain:   %s\n", domainMonitor.Domain)
		fmt.Printf("Interval: %s\n", output.FormatDuration(domainMonitor.Interval))
		fmt.Printf("Status:   %s\n", domainMonitor.Status)

		return nil
	},
}

// domains pause <id>
var domainsPauseCmd = &cobra.Command{
	Use:   "pause <id>",
	Short: "Pause a domain monitor",
	Long:  "Pause a domain expiration monitor (sets status to paused)",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDomainID(client, args[0])
		if err != nil {
			return err
		}

		// Update status to paused
		status := "paused"
		req := &api.UpdateDomainMonitorRequest{Status: &status}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		_, err = client.UpdateDomain(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to pause domain monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Domain monitor %s paused successfully", args[0]))
		return nil
	},
}

// domains resume <id>
var domainsResumeCmd = &cobra.Command{
	Use:   "resume <id>",
	Short: "Resume a domain monitor",
	Long:  "Resume a paused domain expiration monitor (sets status to active)",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDomainID(client, args[0])
		if err != nil {
			return err
		}

		// Update status to active
		status := "active"
		req := &api.UpdateDomainMonitorRequest{Status: &status}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		_, err = client.UpdateDomain(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to resume domain monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Domain monitor %s resumed successfully", args[0]))
		return nil
	},
}

// domains incidents <id>
var domainsIncidentsCmd = &cobra.Command{
	Use:   "incidents <id>",
	Short: "Show incident history",
	Long:  "Display incident history (downtime periods) for a domain monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDomainID(client, args[0])
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

		incidents, err := client.ListDomainIncidents(fullID)

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
			output.InfoMessage("No incidents found - this domain monitor has been running smoothly!")
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

// domains delete <id>
var domainsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a domain monitor",
	Long:  "Delete a domain expiration monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDomainID(client, args[0])
		if err != nil {
			return err
		}

		// Confirm deletion
		confirm, _ := cmd.Flags().GetBool("force")
		if !confirm {
			fmt.Printf("Are you sure you want to delete domain monitor %s? (y/N): ", args[0])
			var response string
			_, _ = fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cancelled")
				return nil
			}
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		err = client.DeleteDomain(fullID)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to delete domain monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Domain monitor %s deleted successfully", args[0]))
		return nil
	},
}

// Helper function to resolve a short domain ID to a full ID
func resolveDomainID(client *api.Client, shortID string) (string, error) {
	// If it looks like a full UUID, use it as-is
	if len(shortID) >= 32 {
		return shortID, nil
	}

	// Otherwise, fetch all domains and match by prefix
	result, err := client.ListDomains()
	if err != nil {
		return "", fmt.Errorf("failed to list domain monitors: %w", err)
	}

	var matches []string
	for _, domain := range result.DomainMonitors {
		if len(domain.ID) >= len(shortID) && domain.ID[:len(shortID)] == shortID {
			matches = append(matches, domain.ID)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no domain monitor found with ID prefix '%s'", shortID)
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("ambiguous ID prefix '%s' matches multiple domain monitors", shortID)
	}

	return matches[0], nil
}

func init() {
	// Add flags to list command
	domainsListCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to show command
	domainsShowCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to create command
	domainsCreateCmd.Flags().String("name", "", "Domain monitor name (required)")
	domainsCreateCmd.Flags().String("domain", "", "Domain to monitor (required)")
	domainsCreateCmd.Flags().Int("interval", 1440, "Check interval in minutes (default: daily)")
	domainsCreateCmd.Flags().Int("grace-period", 0, "Grace period in minutes")
	domainsCreateCmd.Flags().Int("warning-threshold", 30, "Warning threshold in days")
	domainsCreateCmd.Flags().Int("urgent-threshold", 14, "Urgent threshold in days")
	domainsCreateCmd.Flags().Int("critical-threshold", 7, "Critical threshold in days")
	_ = domainsCreateCmd.MarkFlagRequired("name")
	_ = domainsCreateCmd.MarkFlagRequired("domain")

	// Add flags to update command
	domainsUpdateCmd.Flags().String("name", "", "Domain monitor name")
	domainsUpdateCmd.Flags().String("domain", "", "Domain to monitor")
	domainsUpdateCmd.Flags().Int("interval", 0, "Check interval in minutes")
	domainsUpdateCmd.Flags().Int("grace-period", 0, "Grace period in minutes")
	domainsUpdateCmd.Flags().Int("warning-threshold", 0, "Warning threshold in days")
	domainsUpdateCmd.Flags().Int("urgent-threshold", 0, "Urgent threshold in days")
	domainsUpdateCmd.Flags().Int("critical-threshold", 0, "Critical threshold in days")
	domainsUpdateCmd.Flags().String("status", "", "Monitor status (active, inactive, paused)")

	// Add flags to incidents command
	domainsIncidentsCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to delete command
	domainsDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	// Add subcommands
	domainsCmd.AddCommand(domainsListCmd)
	domainsCmd.AddCommand(domainsShowCmd)
	domainsCmd.AddCommand(domainsCreateCmd)
	domainsCmd.AddCommand(domainsUpdateCmd)
	domainsCmd.AddCommand(domainsPauseCmd)
	domainsCmd.AddCommand(domainsResumeCmd)
	domainsCmd.AddCommand(domainsIncidentsCmd)
	domainsCmd.AddCommand(domainsDeleteCmd)

	// Add domains command to root
	rootCmd.AddCommand(domainsCmd)
}
