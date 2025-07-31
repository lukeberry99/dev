package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of devtool",
	Long:  `Print the version and build information for devtool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("devtool v1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
