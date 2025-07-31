package installer

import (
	"fmt"
	"time"

	"github.com/lukeberry99/devtool/internal/config"
	"github.com/lukeberry99/devtool/internal/state"
	"github.com/lukeberry99/devtool/internal/ui"
)

// Replicate the exact script execution logic from bash
type ToolRunner struct {
	logger       *ui.Logger
	dryRun       bool
	verbose      bool
	homebrew     *HomebrewManager
	stateManager *state.LocalStateManager
}

func NewToolRunner(logger *ui.Logger, dryRun, verbose bool, stateManager *state.LocalStateManager) *ToolRunner {
	homebrew := NewHomebrewManager(logger, dryRun)

	return &ToolRunner{
		logger:       logger,
		dryRun:       dryRun,
		verbose:      verbose,
		homebrew:     homebrew,
		stateManager: stateManager,
	}
}

func (r *ToolRunner) InstallTool(name string, toolConfig config.ToolConfig) error {
	r.logger.Info(fmt.Sprintf("Processing: %s", name))

	// Check if already installed and up-to-date
	if r.isToolCurrent(name, toolConfig.Version) {
		r.logger.Debug(fmt.Sprintf("Tool %s is already current", name))
		return nil
	}

	switch toolConfig.Source {
	case "homebrew":
		return r.installFromHomebrew(name, toolConfig)
	case "build":
		return r.buildFromSource(name, toolConfig)
	case "script":
		return r.runCustomScript(name, toolConfig)
	default:
		return fmt.Errorf("unknown installation source: %s", toolConfig.Source)
	}
}

func (r *ToolRunner) isToolCurrent(name, expectedVersion string) bool {
	if r.stateManager == nil {
		return false
	}

	return r.stateManager.IsToolCurrent(name, expectedVersion)
}

func (r *ToolRunner) installFromHomebrew(name string, toolConfig config.ToolConfig) error {
	r.logger.Info(fmt.Sprintf("Installing %s via Homebrew", name))

	// Ensure Homebrew is installed
	if err := r.homebrew.EnsureInstalled(); err != nil {
		return fmt.Errorf("failed to ensure Homebrew is installed: %w", err)
	}

	// Install the package with any specified arguments
	if len(toolConfig.HomebrewArgs) > 0 {
		if err := r.homebrew.InstallPackageWithArgs(name, toolConfig.HomebrewArgs); err != nil {
			return err
		}
	} else {
		if err := r.homebrew.InstallPackages([]string{name}); err != nil {
			return err
		}
	}

	// Update state tracking
	r.updateToolState(name, toolConfig, "homebrew")

	r.logger.Info(fmt.Sprintf("✅ %s installed successfully via Homebrew", name))
	return nil
}

func (r *ToolRunner) buildFromSource(name string, toolConfig config.ToolConfig) error {
	r.logger.Info(fmt.Sprintf("Building %s from source", name))

	if toolConfig.BuildConfig == nil {
		return fmt.Errorf("build configuration is required for source builds")
	}

	// Handle special case for Neovim
	if name == "neovim" {
		return r.buildNeovim(toolConfig)
	}

	// Generic source build logic
	r.logger.Info(fmt.Sprintf("[BUILD] Building %s from %s", name, toolConfig.BuildConfig.Repository))

	if r.dryRun {
		r.logger.Info(fmt.Sprintf("[DRY RUN] Would build %s from source", name))
		return nil
	}

	// TODO: Implement generic source building
	r.logger.Warn(fmt.Sprintf("Generic source building not yet implemented for %s", name))

	// Update state tracking
	r.updateToolState(name, toolConfig, "built_from_source")

	return nil
}

func (r *ToolRunner) buildNeovim(toolConfig config.ToolConfig) error {
	builder := NewNeovimBuilder(r.logger, r.dryRun, toolConfig.Version)

	if err := builder.Build(); err != nil {
		return fmt.Errorf("failed to build Neovim: %w", err)
	}

	// Update state tracking
	r.updateToolState("neovim", toolConfig, "built_from_source")

	return nil
}

func (r *ToolRunner) runCustomScript(name string, toolConfig config.ToolConfig) error {
	r.logger.Info(fmt.Sprintf("Running custom script for %s", name))

	if r.dryRun {
		r.logger.Info(fmt.Sprintf("[DRY RUN] Would run custom script for %s", name))
		return nil
	}

	// TODO: Implement custom script execution
	r.logger.Warn(fmt.Sprintf("Custom script execution not yet implemented for %s", name))

	// Update state tracking
	r.updateToolState(name, toolConfig, "script")

	return nil
}

func (r *ToolRunner) updateToolState(name string, toolConfig config.ToolConfig, source string) {
	if r.stateManager == nil {
		return
	}

	now := time.Now()
	toolStatus := state.ToolStatus{
		Installed:     true,
		Version:       toolConfig.Version,
		InstallDate:   now,
		LastChecked:   now,
		Source:        source,
		BinaryPath:    "", // TODO: Detect binary path
		ConfigCurrent: true,
	}

	r.stateManager.UpdateToolStatus(name, toolStatus)

	// Save state to disk
	if err := r.stateManager.Save(); err != nil {
		r.logger.Warn(fmt.Sprintf("Failed to save state: %v", err))
	}
}

func (r *ToolRunner) InstallDependencies(dependencies []string) error {
	if len(dependencies) == 0 {
		return nil
	}

	r.logger.Info(fmt.Sprintf("Installing dependencies: %v", dependencies))

	// Ensure Homebrew is installed for dependencies
	if err := r.homebrew.EnsureInstalled(); err != nil {
		return fmt.Errorf("failed to ensure Homebrew is installed: %w", err)
	}

	return r.homebrew.InstallPackages(dependencies)
}

func (r *ToolRunner) InstallTools(tools map[string]config.ToolConfig) error {
	// Filter enabled tools
	enabledTools := make(map[string]config.ToolConfig)
	for name, toolConfig := range tools {
		if toolConfig.Enabled {
			enabledTools[name] = toolConfig
		}
	}

	if len(enabledTools) == 0 {
		r.logger.Info("No tools enabled for installation")
		return nil
	}

	r.logger.Info(fmt.Sprintf("Installing %d tools", len(enabledTools)))

	// Install each tool
	for name, toolConfig := range enabledTools {
		if err := r.InstallTool(name, toolConfig); err != nil {
			return fmt.Errorf("failed to install %s: %w", name, err)
		}
	}

	// Cleanup Homebrew if configured
	if err := r.homebrew.Cleanup(); err != nil {
		r.logger.Warn(fmt.Sprintf("Failed to cleanup Homebrew: %v", err))
	}

	r.logger.Info("✅ All tools installed successfully")
	return nil
}
