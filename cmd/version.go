package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of the CLI
	Version = "dev"
	// Commit is the git commit hash
	Commit = "none"
	// Date is the build date
	Date = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("groovekit version %s\n", Version)
		fmt.Printf("commit: %s\n", Commit)
		fmt.Printf("built: %s\n", Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
