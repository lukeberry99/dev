package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lukeberry99/devtool/internal/config"
	"github.com/lukeberry99/devtool/internal/installer"
	"github.com/lukeberry99/devtool/internal/state"
	"github.com/lukeberry99/devtool/internal/ui"
)

var installCmd = &cobra.Command{
	Use:   "install [pattern]",
	Short: "Install development tools and dependencies",
	Long: `Install development tools via Homebrew and custom build processes.
	
Equivalent to running ./install from the bash version.
Optionally filter which tools to install with a pattern.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runInstall,
}

func runInstall(cmd *cobra.Command, args []string) {
	dryRun := viper.GetBool("dry-run")
	verbose := viper.GetBool("verbose")
	force, _ := cmd.Flags().GetBool("force")

	// Initialize logger
	logger := ui.NewLogger(verbose)

	logger.Section("🚀 Starting Installation")

	if dryRun {
		logger.Step("DRY RUN MODE: No actual installations will be performed")
	}

	if force {
		logger.Step("FORCE MODE: Will reinstall tools even if they appear current")
	}

	// Initialize state manager
	stateManager, err := state.NewLocalStateManager()
	if err != nil {
		logger.Errorf("Failed to initialize state manager: %v", err)
		return
	}

	// Load configuration
	configFile := viper.GetString("config")
	cfg, err := config.Load(configFile)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		return
	}

	// Initialize tool runner
	runner := installer.NewToolRunner(logger, dryRun, verbose, force, stateManager)

	// Install tools
	if err := runner.InstallTools(cfg.Tools); err != nil {
		logger.Error(fmt.Sprintf("Installation failed: %v", err))
		return
	}

	logger.Success("Installation completed successfully!")
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringSlice("tools", []string{}, "Specific tools to install")
	installCmd.Flags().String("profile", "", "Install tools for specific profile")
	installCmd.Flags().Bool("force", false, "Force reinstall even if tools appear current")
}
