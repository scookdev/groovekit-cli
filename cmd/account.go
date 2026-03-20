// Package cmd provides the command-line interface commands for GrooveKit CLI
package cmd

import (
	"fmt"
	"strings"
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
	RunE: func(cmd *cobra.Command, _ []string) error {
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
			jobUsage := fmt.Sprintf("%d / %d", account.JobMonitorCount, account.Subscription.MaxJobs)
			jobPercent := 0.0
			if account.Subscription.MaxJobs > 0 {
				jobPercent = float64(account.JobMonitorCount) / float64(account.Subscription.MaxJobs) * 100
			}
			fmt.Printf("Job Monitors:         %s %s\n", jobUsage, formatUsageBar(jobPercent))

			// Api Monitors
			apiMonitorUsage := fmt.Sprintf("%d / %d", account.ApiMonitorCount, account.Subscription.MaxApiMonitors)
			apiMonitorPercent := 0.0
			if account.Subscription.MaxApiMonitors > 0 {
				apiMonitorPercent = float64(account.ApiMonitorCount) / float64(account.Subscription.MaxApiMonitors) * 100
			}
			fmt.Printf("Api Monitors:         %s %s\n", apiMonitorUsage, formatUsageBar(apiMonitorPercent))

			// Ssl Monitors
			sslMonitorUsage := fmt.Sprintf("%d / %d", account.SslMonitorCount, account.Subscription.MaxSslMonitors)
			sslMonitorPercent := 0.0
			if account.Subscription.MaxSslMonitors > 0 {
				sslMonitorPercent = float64(account.SslMonitorCount) / float64(account.Subscription.MaxSslMonitors) * 100
			}
			fmt.Printf("SSL Monitors:         %s %s\n", sslMonitorUsage, formatUsageBar(sslMonitorPercent))

			// Domain Monitors
			domainMonitorUsage := fmt.Sprintf("%d / %d", account.DomainMonitorCount, account.Subscription.MaxDomainMonitors)
			domainMonitorPercent := 0.0
			if account.Subscription.MaxDomainMonitors > 0 {
				domainMonitorPercent = float64(account.DomainMonitorCount) / float64(account.Subscription.MaxDomainMonitors) * 100
			}
			fmt.Printf("Domain Monitors:      %s %s\n", domainMonitorUsage , formatUsageBar(domainMonitorPercent))
			
			// Dns Monitors
			dnsMonitorUsage := fmt.Sprintf("%d / %d", account.DnsMonitorCount, account.Subscription.MaxDnsMonitors)
			dnsMonitorPercent := 0.0
			if account.Subscription.MaxDnsMonitors > 0 {
				dnsMonitorPercent = float64(account.DnsMonitorCount) / float64(account.Subscription.MaxDnsMonitors) * 100
			}
			fmt.Printf("DNS Monitors:         %s %s\n", dnsMonitorUsage, formatUsageBar(dnsMonitorPercent))

			// SMS
			if account.Subscription.SMSLimit > 0 {
				smsUsage := fmt.Sprintf("%d / %d", account.SMSUsed, account.Subscription.SMSLimit)
				smsPercent := 0.0
				if account.Subscription.SMSLimit > 0 {
					smsPercent = float64(account.SMSUsed) / float64(account.Subscription.SMSLimit) * 100
				}
				fmt.Printf("SMS this month:       %s %s\n", smsUsage, formatUsageBar(smsPercent))
			} else {
				fmt.Printf("SMS this month:       %s\n", output.Yellow("Not available on this plan"))
			}

			// Check interval
			fmt.Printf("Min check interval:   %s\n",  output.FormatDuration(account.Subscription.MinCheckInterval))
		} else {
			fmt.Printf("\n%s\n",                output.Yellow("No active subscription"))
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

	var sb strings.Builder
	sb.WriteString("[")
	for i := range barLength {
		if i < filled {
			switch {
			case percent >= 90:
				sb.WriteString(output.Red("█"))
			case percent >= 75:
				sb.WriteString(output.Yellow("█"))
			default:
				sb.WriteString(output.Green("█"))
			}
		} else {
			sb.WriteString("░")
		}
	}
	fmt.Fprintf(&sb, "] %.0f%%", percent)
	bar := sb.String()

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
