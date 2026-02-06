package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/scookdev/groovekit-cli/internal/output"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "View account information",
	Long:  "View your account details, plan limits, and usage",
}

// account show
var accountShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show account details",
	Long:  "Display your account information, plan limits, and current usage",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")

		// Start spinner
		var s *spinner.Spinner
		if !jsonOutput {
			s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Start()
		}

		account, err := client.GetAccount()

		// Stop spinner
		if s != nil {
			s.Stop()
		}

		if err != nil {
			return fmt.Errorf("failed to get account: %w", err)
		}

		if jsonOutput {
			return outputJSON(account)
		}

		// Print account details
		fmt.Printf("%s\n\n", output.Bold("Account Information"))
		fmt.Printf("Email:            %s\n", account.Email)
		fmt.Printf("Name:             %s\n", account.FullName)

		if account.Subscription != nil {
			fmt.Printf("\n%s\n\n", output.Bold("Subscription"))
			fmt.Printf("Plan:             %s\n", output.Cyan(account.Subscription.PlanName))
			fmt.Printf("Status:           %s\n", formatStatus(account.Subscription.Status))

			if account.Subscription.CurrentPeriodEnd != nil {
				fmt.Printf("Renews:           %s\n", *account.Subscription.CurrentPeriodEnd)
			}

			// Usage and Limits
			fmt.Printf("\n%s\n\n", output.Bold("Usage & Limits"))

			// Jobs
			jobUsage := fmt.Sprintf("%d / %d", account.JobCount, account.Subscription.MaxJobs)
			jobPercent := 0.0
			if account.Subscription.MaxJobs > 0 {
				jobPercent = float64(account.JobCount) / float64(account.Subscription.MaxJobs) * 100
			}
			fmt.Printf("Jobs:             %s %s\n", jobUsage, formatUsageBar(jobPercent))

			// Monitors
			monitorUsage := fmt.Sprintf("%d / %d", account.MonitorCount, account.Subscription.MaxMonitors)
			monitorPercent := 0.0
			if account.Subscription.MaxMonitors > 0 {
				monitorPercent = float64(account.MonitorCount) / float64(account.Subscription.MaxMonitors) * 100
			}
			fmt.Printf("Monitors:         %s %s\n", monitorUsage, formatUsageBar(monitorPercent))

			// SMS
			if account.Subscription.SMSLimit > 0 {
				smsUsage := fmt.Sprintf("%d / %d", account.SMSUsed, account.Subscription.SMSLimit)
				smsPercent := 0.0
				if account.Subscription.SMSLimit > 0 {
					smsPercent = float64(account.SMSUsed) / float64(account.Subscription.SMSLimit) * 100
				}
				fmt.Printf("SMS this month:   %s %s\n", smsUsage, formatUsageBar(smsPercent))
			} else {
				fmt.Printf("SMS this month:   %s\n", output.Yellow("Not available on this plan"))
			}

			// Check interval
			fmt.Printf("Min check interval: %s\n", output.FormatDuration(account.Subscription.MinCheckInterval))
		} else {
			fmt.Printf("\n%s\n", output.Yellow("No active subscription"))
		}

		return nil
	},
}

// Helper function to format status with color
func formatStatus(status string) string {
	switch status {
	case "active":
		return output.Green(status)
	case "canceled", "past_due":
		return output.Red(status)
	case "trialing":
		return output.Cyan(status)
	default:
		return status
	}
}

// Helper function to format usage bar
func formatUsageBar(percent float64) string {
	barLength := 20
	filled := int(percent / 100 * float64(barLength))

	bar := "["
	for i := 0; i < barLength; i++ {
		if i < filled {
			if percent >= 90 {
				bar += output.Red("█")
			} else if percent >= 75 {
				bar += output.Yellow("█")
			} else {
				bar += output.Green("█")
			}
		} else {
			bar += "░"
		}
	}
	bar += fmt.Sprintf("] %.0f%%", percent)

	return bar
}

func init() {
	// Add flags to show command
	accountShowCmd.Flags().Bool("json", false, "Output as JSON")

	// Add subcommands
	accountCmd.AddCommand(accountShowCmd)

	// Add account command to root
	rootCmd.AddCommand(accountCmd)
}
