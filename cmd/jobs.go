package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/scookdev/groovekit-cli/internal/api"
	"github.com/scookdev/groovekit-cli/internal/config"
	"github.com/scookdev/groovekit-cli/internal/output"
	"github.com/spf13/cobra"
)

var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Manage cron job monitors",
	Long:  "List, create, show, and delete cron job heartbeat monitors",
}

// jobs list
var jobsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all jobs",
	Long:  "List all cron job monitors for your account",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		result, err := client.ListJobs()
		if err != nil {
			return fmt.Errorf("failed to list jobs: %w", err)
		}

		// Check for --json flag
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			return outputJSON(result)
		}

		if len(result.Jobs) == 0 {
			output.InfoMessage("No jobs found")
			fmt.Println("\nCreate your first job:")
			fmt.Println("  groovekit jobs create --name 'Daily Backup' --interval 1440")
			return nil
		}

		// Create table
		table := output.NewTable([]string{"ID", "NAME", "INTERVAL", "STATUS", "HEALTH"})
		table.Render()

		// Add rows
		for _, job := range result.Jobs {
			status := job.Status
			if job.Status == "active" {
				status = output.Green(status)
			}

			health := output.Green("✓ Up")
			if job.Down {
				health = output.Red("✗ Down")
			}

			// Truncate ID to first 8 chars (like Docker)
			shortID := job.ID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}

			table.Append([]string{
				output.Cyan(shortID),
				job.Name,
				output.FormatDuration(job.Interval),
				status,
				health,
			})
		}

		table.Flush()
		fmt.Printf("\n%s\n", output.Bold(fmt.Sprintf("Total: %d job(s)", result.TotalCount)))
		return nil
	},
}

// jobs show <id>
var jobsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show job details",
	Long:  "Display detailed information about a specific job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveJobID(client, args[0])
		if err != nil {
			return err
		}

		job, err := client.GetJob(fullID)
		if err != nil {
			return fmt.Errorf("failed to get job: %w", err)
		}

		// Check for --json flag
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			return outputJSON(job)
		}

		// Print job details
		fmt.Printf("ID:            %s\n", job.ID)
		fmt.Printf("Name:          %s\n", job.Name)
		fmt.Printf("Status:        %s\n", job.Status)
		fmt.Printf("Interval:      %s\n", output.FormatDuration(job.Interval))
		fmt.Printf("Grace Period:  %s\n", output.FormatDuration(job.GracePeriod))
		fmt.Printf("Down:          %t\n", job.Down)

		if job.LastPingAt != nil {
			fmt.Printf("Last Ping:     %s\n", *job.LastPingAt)
		} else {
			fmt.Printf("Last Ping:     Never\n")
		}

		if job.LastRunAt != nil {
			fmt.Printf("Last Run:      %s\n", *job.LastRunAt)
		}

		fmt.Printf("\nPing URL:\n")
		fmt.Printf("  curl https://api.groovekit.io/pings/%s\n", job.PingToken)

		if len(job.AllowedIPs) > 0 {
			fmt.Printf("\nAllowed IPs:   %v\n", job.AllowedIPs)
		}

		if job.WebhookURL != "" {
			fmt.Printf("\nWebhook URL:   %s\n", job.WebhookURL)
		}

		return nil
	},
}

// jobs create
var jobsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new job",
	Long:  "Create a new cron job heartbeat monitor",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		req := &api.CreateJobRequest{
			Name:        name,
			Interval:    interval,
			GracePeriod: gracePeriod,
		}

		job, err := client.CreateJob(req)
		if err != nil {
			return fmt.Errorf("failed to create job: %w", err)
		}

		output.SuccessMessage("Job created successfully\n")
		fmt.Printf("ID:           %s\n", output.Cyan(job.ID))
		fmt.Printf("Name:         %s\n", output.Bold(job.Name))
		fmt.Printf("Interval:     %s\n", fmt.Sprintf("%d minutes", job.Interval))
		fmt.Printf("Grace Period: %s\n", fmt.Sprintf("%d minutes", job.GracePeriod))
		fmt.Printf("\n%s\n", output.Bold("Ping URL:"))
		fmt.Printf("  %s\n", output.Cyan(fmt.Sprintf("curl https://api.groovekit.io/pings/%s", job.PingToken)))

		return nil
	},
}

// jobs delete <id>
var jobsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a job",
	Long:  "Delete a cron job monitor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Resolve short ID to full ID
		fullID, err := resolveJobID(client, args[0])
		if err != nil {
			return err
		}

		// Confirm deletion
		confirm, _ := cmd.Flags().GetBool("force")
		if !confirm {
			fmt.Printf("Are you sure you want to delete job %s? (y/N): ", args[0])
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cancelled")
				return nil
			}
		}

		if err := client.DeleteJob(fullID); err != nil {
			return fmt.Errorf("failed to delete job: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Job %s deleted successfully", args[0]))
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
func resolveJobID(client *api.Client, shortID string) (string, error) {
	// If it looks like a full UUID, use it as-is
	if len(shortID) >= 32 {
		return shortID, nil
	}

	// Otherwise, fetch all jobs and match by prefix
	result, err := client.ListJobs()
	if err != nil {
		return "", fmt.Errorf("failed to list jobs: %w", err)
	}

	var matches []string
	for _, job := range result.Jobs {
		if len(job.ID) >= len(shortID) && job.ID[:len(shortID)] == shortID {
			matches = append(matches, job.ID)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no job found with ID prefix '%s'", shortID)
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("ambiguous ID prefix '%s' matches multiple jobs", shortID)
	}

	return matches[0], nil
}

func init() {
	// Add flags to list command
	jobsListCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to show command
	jobsShowCmd.Flags().Bool("json", false, "Output as JSON")

	// Add flags to create command
	jobsCreateCmd.Flags().String("name", "", "Job name (required)")
	jobsCreateCmd.Flags().Int("interval", 0, "Check interval in minutes (required)")
	jobsCreateCmd.Flags().Int("grace-period", 5, "Grace period in minutes")
	jobsCreateCmd.MarkFlagRequired("name")
	jobsCreateCmd.MarkFlagRequired("interval")

	// Add flags to delete command
	jobsDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	// Add subcommands
	jobsCmd.AddCommand(jobsListCmd)
	jobsCmd.AddCommand(jobsShowCmd)
	jobsCmd.AddCommand(jobsCreateCmd)
	jobsCmd.AddCommand(jobsDeleteCmd)

	// Add jobs command to root
	rootCmd.AddCommand(jobsCmd)
}
