package cmd

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/scookdev/groovekit-cli/internal/api"
	"github.com/scookdev/groovekit-cli/internal/config"
	"github.com/scookdev/groovekit-cli/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  "Login, logout, and check authentication status",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to GrooveKit",
	Long:  "Authenticate with your GrooveKit account and save credentials locally",
	RunE: func(_ *cobra.Command, _ []string) error {
		// Prompt for email
		fmt.Print("Email: ")
		var email string
		_, _ = fmt.Scanln(&email)

		// Prompt for password (hidden)
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println() // New line after password input
		password := string(passwordBytes)

		// Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Create API client and login
		client := api.NewClient(cfg)

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Start()
		token, err := client.Login(email, password)
		s.Stop()

		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		// Save credentials
		cfg.AccessToken = token
		cfg.Email = email
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		output.SuccessMessage(fmt.Sprintf("Logged in successfully as %s", output.Bold(email)))
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from GrooveKit",
	Long:  "Remove locally stored credentials",
	RunE: func(_ *cobra.Command, _ []string) error {
		if err := config.Clear(); err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Not currently logged in")
				return nil
			}
			return fmt.Errorf("failed to logout: %w", err)
		}

		output.SuccessMessage("Logged out successfully")
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(authCmd)
}
