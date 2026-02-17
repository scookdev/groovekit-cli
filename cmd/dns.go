package cmd

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/scookdev/groovekit-cli/internal/api"
	"github.com/scookdev/groovekit-cli/internal/output"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage DNS record monitors",
	Long:  "List, create, show, update, and delete DNS record monitors",
}

// dns list
var dnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DNS monitors",
	Long:  "List all DNS record monitors for your account",
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

		result, err := client.ListDnsMonitors()

		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to list DNS monitors: %w", err)
		}
		if jsonOutput {
			return outputJSON(result)
		}

		if len(result.DnsMonitors) == 0 {
			output.InfoMessage("No DNS monitors found")
			fmt.Println("\nCreate your first DNS monitor:")
			fmt.Println("  groovekit dns create --name 'Example MX' --domain example.com --type MX --expected mail.example.com")
			return nil
		}

		// Create table
		table := output.NewTable([]string{"ID", "NAME", "DOMAIN", "TYPE", "MISMATCH", "STATUS"})
		table.Render()

		// Add rows
		for _, dns := range result.DnsMonitors {
			status := dns.Status
			if dns.Status == "active" {
				status = output.Green(status)
			}

			// Truncate ID to first 8 chars (like Docker)
			shortID := dns.ID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}

			// Color-code mismatch
			var mismatch string
			if dns.HasMismatch {
				mismatch = output.Red("Yes")
			} else {
				mismatch = output.Green("No")
			}

			table.Append([]string{
				output.Cyan(shortID),
				dns.Name,
				dns.Domain,
				dns.RecordType,
				mismatch,
				status,
			})
		}

		table.Flush()
		fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d DNS monitor(s)", len(result.DnsMonitors))))
		return nil
	},
}

// dns show <id>
var dnsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show DNS monitor details",
	Long:  "Display detailed information about a specific DNS monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDnsMonitorID(client, args[0])
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

		dns, err := client.GetDnsMonitor(fullID)

		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to get DNS monitor: %w", err)
		}
		if jsonOutput {
			return outputJSON(dns)
		}

		// Print DNS monitor details
		fmt.Printf("ID:                       %s\n", output.Cyan(dns.ID))
		fmt.Printf("Name:                     %s\n", output.Bold(dns.Name))
		fmt.Printf("Domain:                   %s\n", dns.Domain)
		fmt.Printf("Record Type:              %s\n", dns.RecordType)
		fmt.Printf("Status:                   %s\n", dns.Status)
		fmt.Printf("Check Interval:           %s\n", output.FormatDuration(dns.Interval))
		fmt.Printf("Grace Period:             %s\n", output.FormatDuration(dns.GracePeriod))

		// Show expected values
		fmt.Printf("\nExpected Values:\n")
		if len(dns.ExpectedValues) == 0 {
			fmt.Printf("  (none)\n")
		} else {
			for _, val := range dns.ExpectedValues {
				fmt.Printf("  - %s\n", output.Green(val))
			}
		}

		// Show current values
		fmt.Printf("\nCurrent Values:\n")
		if len(dns.CurrentValues) == 0 {
			fmt.Printf("  (none)\n")
		} else {
			for _, val := range dns.CurrentValues {
				// Highlight if this value is not in expected values
				if slices.Contains(dns.ExpectedValues, val) {
					fmt.Printf("  - %s (unexpected)\n", output.Red(val))
				} else {
					fmt.Printf("  - %s\n", val)
				}
			}
		}

		// Show mismatch status
		if dns.HasMismatch {
			fmt.Printf("\nMismatch:                 %s\n", output.Red("Yes - values don't match!"))
		} else {
			fmt.Printf("\nMismatch:                 %s\n", output.Green("No - values match"))
		}

		if dns.LastChanged != nil {
			fmt.Printf("Last Changed:             %s\n", *dns.LastChanged)
		}
		fmt.Printf("Last Check At:            %s\n", dns.LastCheckAt)
		fmt.Printf("Last Successful Check:    %s\n", dns.LastSuccessfulCheckAt)
		fmt.Printf("Consecutive Failures:     %d\n", dns.ConsecutiveFailures)
		fmt.Printf("Created At:               %s\n", dns.CreatedAt)
		fmt.Printf("Updated At:               %s\n", dns.UpdatedAt)

		return nil
	},
}

// dns create
var dnsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new DNS monitor",
	Long:  "Create a new DNS record monitor",
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		domain, _ := cmd.Flags().GetString("domain")
		recordType, _ := cmd.Flags().GetString("type")
		expectedValues, _ := cmd.Flags().GetStringSlice("expected")
		interval, _ := cmd.Flags().GetInt("interval")
		gracePeriod, _ := cmd.Flags().GetInt("grace-period")

		if name == "" {
			return fmt.Errorf("--name is required")
		}
		if domain == "" {
			return fmt.Errorf("--domain is required")
		}
		if recordType == "" {
			return fmt.Errorf("--type is required")
		}
		if len(expectedValues) == 0 {
			return fmt.Errorf("--expected is required (at least one value)")
		}

		// Validate record type
		validTypes := []string{"A", "AAAA", "MX", "CNAME", "TXT", "NS"}
		recordType = strings.ToUpper(recordType)
		if !slices.Contains(validTypes, recordType) {
			return fmt.Errorf("invalid record type '%s'. Must be one of: %s", recordType, strings.Join(validTypes, ", "))
		}

		req := &api.CreateDnsMonitorRequest{
			Name:           name,
			Domain:         domain,
			RecordType:     recordType,
			ExpectedValues: expectedValues,
			Interval:       interval,
			GracePeriod:    gracePeriod,
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		dnsMonitor, err := client.CreateDnsMonitor(req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to create DNS monitor: %w", err)
		}

		output.SuccessMessage("DNS monitor created successfully\n")
		fmt.Printf("ID:       %s\n", output.Cyan(dnsMonitor.ID))
		fmt.Printf("Name:     %s\n", output.Bold(dnsMonitor.Name))
		fmt.Printf("Domain:   %s\n", dnsMonitor.Domain)
		fmt.Printf("Type:     %s\n", dnsMonitor.RecordType)
		fmt.Printf("Interval: %s\n", output.FormatDuration(dnsMonitor.Interval))
		fmt.Printf("Expected: %s\n", strings.Join(dnsMonitor.ExpectedValues, ", "))

		return nil
	},
}

// dns update <id>
var dnsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a DNS monitor",
	Long:  "Update an existing DNS record monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDnsMonitorID(client, args[0])
		if err != nil {
			return err
		}

		// Build update request with only provided flags
		req := &api.UpdateDnsMonitorRequest{}
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

		if cmd.Flags().Changed("type") {
			recordType, _ := cmd.Flags().GetString("type")
			recordType = strings.ToUpper(recordType)
			// Validate record type
			validTypes := []string{"A", "AAAA", "MX", "CNAME", "TXT", "NS"}
			if !slices.Contains(validTypes, recordType) {
				return fmt.Errorf("invalid record type '%s'. Must be one of: %s", recordType, strings.Join(validTypes, ", "))
			}
			req.RecordType = &recordType
			hasUpdates = true
		}

		if cmd.Flags().Changed("expected") {
			expectedValues, _ := cmd.Flags().GetStringSlice("expected")
			req.ExpectedValues = &expectedValues
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

		if !hasUpdates {
			return fmt.Errorf("no fields to update. Use --name, --domain, --type, --expected, --interval, --grace-period, or --status")
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		dnsMonitor, err := client.UpdateDnsMonitor(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to update DNS monitor: %w", err)
		}

		output.SuccessMessage("DNS monitor updated successfully\n")
		fmt.Printf("ID:       %s\n", output.Cyan(dnsMonitor.ID))
		fmt.Printf("Name:     %s\n", output.Bold(dnsMonitor.Name))
		fmt.Printf("Domain:   %s\n", dnsMonitor.Domain)
		fmt.Printf("Type:     %s\n", dnsMonitor.RecordType)
		fmt.Printf("Interval: %s\n", output.FormatDuration(dnsMonitor.Interval))
		fmt.Printf("Status:   %s\n", dnsMonitor.Status)

		return nil
	},
}

// dns pause <id>
var dnsPauseCmd = &cobra.Command{
	Use:   "pause <id>",
	Short: "Pause a DNS monitor",
	Long:  "Pause a DNS record monitor (sets status to paused)",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDnsMonitorID(client, args[0])
		if err != nil {
			return err
		}

		// Update status to paused
		status := "paused"
		req := &api.UpdateDnsMonitorRequest{Status: &status}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		_, err = client.UpdateDnsMonitor(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to pause DNS monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("DNS monitor %s paused successfully", args[0]))
		return nil
	},
}

// dns resume <id>
var dnsResumeCmd = &cobra.Command{
	Use:   "resume <id>",
	Short: "Resume a DNS monitor",
	Long:  "Resume a paused DNS record monitor (sets status to active)",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDnsMonitorID(client, args[0])
		if err != nil {
			return err
		}

		// Update status to active
		status := "active"
		req := &api.UpdateDnsMonitorRequest{Status: &status}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		_, err = client.UpdateDnsMonitor(fullID, req)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to resume DNS monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("DNS monitor %s resumed successfully", args[0]))
		return nil
	},
}

// dns incidents <id>
var dnsIncidentsCmd = &cobra.Command{
	Use:   "incidents <id>",
	Short: "Show incident history",
	Long:  "Display incident history (downtime periods) for a DNS monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDnsMonitorID(client, args[0])
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

		incidents, err := client.ListDnsMonitorIncidents(fullID)

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
			output.InfoMessage("No incidents found - this DNS monitor has been running smoothly!")
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

// dns delete <id>
var dnsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a DNS monitor",
	Long:  "Delete a DNS record monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveDnsMonitorID(client, args[0])
		if err != nil {
			return err
		}

		// Confirm deletion
		confirm, _ := cmd.Flags().GetBool("force")
		if !confirm {
			fmt.Printf("Are you sure you want to delete DNS monitor %s? (y/N): ", args[0])
			var response string
			_, _ = fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cancelled")
				return nil
			}
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		err = client.DeleteDnsMonitor(fullID)
		s.Stop()

		if err != nil {
			return fmt.Errorf("failed to delete DNS monitor: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("DNS monitor %s deleted successfully", args[0]))
		return nil
	},
}

// Helper function to resolve a short DNS monitor ID to a full ID
func resolveDnsMonitorID(client *api.Client, shortID string) (string, error) {
	// If it looks like a full UUID, use it as-is
	if len(shortID) >= 32 {
		return shortID, nil
	}

	// Otherwise, fetch all DNS monitors and match by prefix
	result, err := client.ListDnsMonitors()
	if err != nil {
		return "", fmt.Errorf("failed to list DNS monitors: %w", err)
	}

	var matches []string
	for _, dns := range result.DnsMonitors {
		if len(dns.ID) >= len(shortID) && dns.ID[:len(shortID)] == shortID {
			matches = append(matches, dns.ID)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no DNS monitor found with ID prefix '%s'", shortID)
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("ambiguous ID prefix '%s' matches multiple DNS monitors", shortID)
	}

	return matches[0], nil
}

func init() {
	// Add flags to list command
	dnsListCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to show command
	dnsShowCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to create command
	dnsCreateCmd.Flags().String("name", "", "DNS monitor name (required)")
	dnsCreateCmd.Flags().String("domain", "", "Domain to monitor (required)")
	dnsCreateCmd.Flags().String("type", "", "DNS record type: A, AAAA, MX, CNAME, TXT, NS (required)")
	dnsCreateCmd.Flags().StringSlice("expected", []string{}, "Expected value(s) - can be specified multiple times or comma-separated (required)")
	dnsCreateCmd.Flags().Int("interval", 1440, "Check interval in minutes (default: daily)")
	dnsCreateCmd.Flags().Int("grace-period", 0, "Grace period in minutes")
	_ = dnsCreateCmd.MarkFlagRequired("name")
	_ = dnsCreateCmd.MarkFlagRequired("domain")
	_ = dnsCreateCmd.MarkFlagRequired("type")
	_ = dnsCreateCmd.MarkFlagRequired("expected")

	// Add flags to update command
	dnsUpdateCmd.Flags().String("name", "", "DNS monitor name")
	dnsUpdateCmd.Flags().String("domain", "", "Domain to monitor")
	dnsUpdateCmd.Flags().String("type", "", "DNS record type: A, AAAA, MX, CNAME, TXT, NS")
	dnsUpdateCmd.Flags().StringSlice("expected", []string{}, "Expected value(s) - can be specified multiple times or comma-separated")
	dnsUpdateCmd.Flags().Int("interval", 0, "Check interval in minutes")
	dnsUpdateCmd.Flags().Int("grace-period", 0, "Grace period in minutes")
	dnsUpdateCmd.Flags().String("status", "", "Monitor status (active, inactive, paused)")

	// Add flags to incidents command
	dnsIncidentsCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to delete command
	dnsDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	// Add subcommands
	dnsCmd.AddCommand(dnsListCmd)
	dnsCmd.AddCommand(dnsShowCmd)
	dnsCmd.AddCommand(dnsCreateCmd)
	dnsCmd.AddCommand(dnsUpdateCmd)
	dnsCmd.AddCommand(dnsPauseCmd)
	dnsCmd.AddCommand(dnsResumeCmd)
	dnsCmd.AddCommand(dnsIncidentsCmd)
	dnsCmd.AddCommand(dnsDeleteCmd)

	// Add dns command to root
	rootCmd.AddCommand(dnsCmd)
}
