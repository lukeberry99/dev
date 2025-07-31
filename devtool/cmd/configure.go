package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lukeberry99/devtool/internal/config"
	"github.com/lukeberry99/devtool/internal/configurator"
	"github.com/lukeberry99/devtool/internal/ui"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Deploy configuration files and dotfiles",
	Long: `Deploy configuration files from env/ directory to system locations.
	
Equivalent to running ./configure from the bash version.
This includes .config directories, .local scripts, and dotfiles.`,
	Run: runConfigure,
}

func runConfigure(cmd *cobra.Command, args []string) {
	dryRun := viper.GetBool("dry-run")
	verbose := viper.GetBool("verbose")

	// Initialize logger
	logger := ui.NewLogger(verbose)

	logger.Info("Starting configuration deployment...")

	if dryRun {
		logger.Info("DRY RUN MODE: No actual file operations will be performed")
	}

	// Load configuration
	configFile := viper.GetString("config")
	cfg, err := config.Load(configFile)
	if err != nil {
		logger.Errorf("Failed to load configuration: %v", err)
		return
	}

	// Initialize dotfiles manager
	dotfilesManager, err := configurator.NewDotfilesManager(cfg, logger, dryRun)
	if err != nil {
		logger.Errorf("Failed to initialize dotfiles manager: %v", err)
		return
	}

	// Deploy configuration files
	if err := dotfilesManager.Deploy(); err != nil {
		logger.Errorf("Configuration deployment failed: %v", err)
		return
	}

	logger.Info("✅ Configuration deployment completed successfully")
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
