package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of installed tools and configuration",
	Long: `Display information about which tools are installed,
their versions, and the status of configuration files.`,
	Run: runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	fmt.Println("Status command called")

	// TODO: Implement actual status checking logic
	fmt.Println("Status checking logic not yet implemented")
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
