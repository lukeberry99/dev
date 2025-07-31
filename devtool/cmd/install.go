package cmd

import (
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

	// Initialize logger
	logger := ui.NewLogger(verbose)

	logger.Info("Starting tool installation...")

	if dryRun {
		logger.Info("DRY RUN MODE: No actual installations will be performed")
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
		logger.Errorf("Failed to load configuration: %v", err)
		// Create a minimal default config for basic functionality
		cfg = &config.Config{
			Tools: getDefaultTools(),
		}
		logger.Info("Using default tool configuration")
	}

	// Initialize tool runner
	runner := installer.NewToolRunner(logger, dryRun, verbose, stateManager)

	// Install tools
	if err := runner.InstallTools(cfg.Tools); err != nil {
		logger.Errorf("Installation failed: %v", err)
		return
	}

	logger.Info("✅ Installation completed successfully")
}

func getDefaultTools() map[string]config.ToolConfig {
	return map[string]config.ToolConfig{
		"git": {
			Source:  "homebrew",
			Enabled: true,
		},
		"tmux": {
			Source:  "homebrew",
			Enabled: true,
		},
		"ripgrep": {
			Source:  "homebrew",
			Enabled: true,
		},
		"fzf": {
			Source:  "homebrew",
			Enabled: true,
		},
		"jq": {
			Source:  "homebrew",
			Enabled: true,
		},
	}
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringSlice("tools", []string{}, "Specific tools to install")
	installCmd.Flags().String("profile", "", "Install tools for specific profile")
}
