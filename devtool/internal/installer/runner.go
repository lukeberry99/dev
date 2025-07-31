package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	force        bool
	homebrew     *HomebrewManager
	stateManager *state.LocalStateManager
	detector     *state.ToolDetector
}

func NewToolRunner(logger *ui.Logger, dryRun, verbose, force bool, stateManager *state.LocalStateManager) *ToolRunner {
	homebrew := NewHomebrewManager(logger, dryRun)

	return &ToolRunner{
		logger:       logger,
		dryRun:       dryRun,
		verbose:      verbose,
		force:        force,
		homebrew:     homebrew,
		stateManager: stateManager,
		detector:     state.NewToolDetector(),
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

func (r *ToolRunner) validateToolExists(name string) bool {
	return r.detector.IsInstalled(name)
}

func (r *ToolRunner) isToolCurrent(name, expectedVersion string) bool {
	if r.force {
		r.logger.Debug(fmt.Sprintf("Force mode enabled - will reinstall %s", name))
		return false // --force flag bypasses all checks
	}

	if r.stateManager == nil {
		return false
	}

	status, exists := r.stateManager.GetToolStatus(name)
	if !exists || !status.Installed {
		r.logger.Debug(fmt.Sprintf("Tool %s not found in state or marked as not installed", name))
		return false
	}

	// Validate tool actually exists on system
	if !r.validateToolExists(name) {
		r.logger.Debug(fmt.Sprintf("Tool %s marked as installed but not found on system", name))
		return false
	}

	// If config doesn't specify version, any installed version is OK
	if expectedVersion == "" {
		r.logger.Debug(fmt.Sprintf("Tool %s is current (no version requirement)", name))
		return true
	}

	// Strict version matching when config specifies version
	if status.Version == expectedVersion {
		r.logger.Debug(fmt.Sprintf("Tool %s is current (version %s matches)", name, expectedVersion))
		return true
	}

	r.logger.Debug(fmt.Sprintf("Tool %s version mismatch: have %s, want %s", name, status.Version, expectedVersion))
	return false
}

func (r *ToolRunner) installFromHomebrew(name string, toolConfig config.ToolConfig) error {
	r.logger.Info(fmt.Sprintf("Installing %s via Homebrew", name))

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

	buildConfig := toolConfig.BuildConfig
	r.logger.Info(fmt.Sprintf("[BUILD] Building %s from %s", name, buildConfig.Repository))

	if r.dryRun {
		r.logger.Info(fmt.Sprintf("[DRY RUN] Would build %s from source", name))
		return nil
	}

	// 1. Install dependencies if specified
	if len(buildConfig.Dependencies) > 0 {
		if err := r.InstallDependencies(buildConfig.Dependencies); err != nil {
			return fmt.Errorf("failed to install dependencies for %s: %w", name, err)
		}
	}

	// 2. Prepare repository
	repoDir, err := r.prepareRepository(name, buildConfig.Repository, toolConfig.Version)
	if err != nil {
		return fmt.Errorf("failed to prepare repository for %s: %w", name, err)
	}

	// 3. Execute build steps
	if err := r.executeBuildSteps(name, repoDir, buildConfig.BuildSteps); err != nil {
		return fmt.Errorf("failed to build %s: %w", name, err)
	}

	// 4. Execute install steps
	if err := r.executeInstallSteps(name, repoDir, buildConfig.InstallSteps); err != nil {
		return fmt.Errorf("failed to install %s: %w", name, err)
	}

	// Update state tracking
	r.updateToolState(name, toolConfig, "built_from_source")

	r.logger.Info(fmt.Sprintf("✅ %s built and installed successfully", name))
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
	if r.stateManager == nil || r.dryRun {
		return
	}

	now := time.Now()

	// Detect actual installed version
	actualVersion := toolConfig.Version // Default to config version
	if source == "homebrew" {
		// Use Homebrew to get version for Homebrew-installed tools
		if brewVersion, err := r.homebrew.GetInstalledVersion(name); err == nil {
			actualVersion = brewVersion
			r.logger.Debug(fmt.Sprintf("Detected %s version via Homebrew: %s", name, actualVersion))
		} else {
			r.logger.Debug(fmt.Sprintf("Could not detect %s version via Homebrew: %v", name, err))
		}
	} else {
		// Fall back to tool-specific detection for non-Homebrew tools
		if detectedVersion, err := r.detector.GetVersion(name); err == nil {
			actualVersion = detectedVersion
			r.logger.Debug(fmt.Sprintf("Detected %s version: %s", name, actualVersion))
		} else {
			r.logger.Debug(fmt.Sprintf("Could not detect %s version: %v", name, err))
		}
	}

	// Detect binary path
	binaryPath := ""
	if path, err := exec.LookPath(name); err == nil {
		binaryPath = path
		r.logger.Debug(fmt.Sprintf("Found %s at: %s", name, binaryPath))
	}

	toolStatus := state.ToolStatus{
		Installed:     true,
		Version:       actualVersion,
		InstallDate:   now,
		LastChecked:   now,
		Source:        source,
		BinaryPath:    binaryPath,
		ConfigCurrent: true, // TODO: Compare actualVersion with toolConfig.Version
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

	return r.homebrew.InstallPackages(dependencies)
}

func (r *ToolRunner) InstallTools(tools map[string]config.ToolConfig) error {
	// Filter enabled tools and check if any need Homebrew
	enabledTools := make(map[string]config.ToolConfig)
	needsHomebrew := false

	for name, toolConfig := range tools {
		if toolConfig.Enabled {
			enabledTools[name] = toolConfig
			if toolConfig.Source == "homebrew" {
				needsHomebrew = true
			}
		}
	}

	if len(enabledTools) == 0 {
		r.logger.Info("No tools enabled for installation")
		return nil
	}

	// Ensure Homebrew once upfront if any tools need it
	if needsHomebrew {
		if err := r.homebrew.EnsureInstalled(); err != nil {
			return fmt.Errorf("failed to ensure Homebrew is installed: %w", err)
		}
	}

	r.logger.Info(fmt.Sprintf("Installing %d tools", len(enabledTools)))

	// Install each tool
	for name, toolConfig := range enabledTools {
		if err := r.InstallTool(name, toolConfig); err != nil {
			return fmt.Errorf("failed to install %s: %w", name, err)
		}
	}

	// Cleanup Homebrew if configured
	if needsHomebrew {
		if err := r.homebrew.Cleanup(); err != nil {
			r.logger.Warn(fmt.Sprintf("Failed to cleanup Homebrew: %v", err))
		}
	}

	r.logger.Info("✅ All tools installed successfully")
	return nil
}

func (r *ToolRunner) prepareRepository(name, repository, version string) (string, error) {
	r.logger.Info(fmt.Sprintf("Preparing %s repository...", name))

	// Create build directory
	homeDir, _ := os.UserHomeDir()
	buildDir := filepath.Join(homeDir, ".devtool", "builds", name)
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create build directory: %w", err)
	}

	repoPath := filepath.Join(buildDir, name)

	// Check if repository already exists
	if _, err := os.Stat(filepath.Join(repoPath, ".git")); err == nil {
		r.logger.Debug("Repository exists, updating...")
		if err := r.updateRepository(repoPath, version); err != nil {
			return "", err
		}
	} else {
		// Clone repository
		r.logger.Info(fmt.Sprintf("Cloning %s repository...", name))
		cmd := exec.Command("git", "clone", repository, repoPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to clone repository: %w", err)
		}

		// Checkout appropriate version
		if err := r.checkoutVersion(repoPath, version); err != nil {
			return "", err
		}
	}

	return repoPath, nil
}

func (r *ToolRunner) updateRepository(repoPath, version string) error {
	// Change to repository directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(repoPath); err != nil {
		return fmt.Errorf("failed to change to repository directory: %w", err)
	}

	// Fetch latest changes
	cmd := exec.Command("git", "fetch", "origin")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch updates: %w", err)
	}

	return r.checkoutVersion(repoPath, version)
}

func (r *ToolRunner) checkoutVersion(repoPath, version string) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(repoPath); err != nil {
		return fmt.Errorf("failed to change to repository directory: %w", err)
	}

	var target string
	switch version {
	case "nightly", "":
		target = "master"
	case "stable":
		target = "stable"
	default:
		target = version
	}

	r.logger.Debug(fmt.Sprintf("Checking out %s", target))
	cmd := exec.Command("git", "checkout", target)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout %s: %w", target, err)
	}

	// Pull latest changes if on a branch
	if target == "master" || target == "stable" {
		cmd = exec.Command("git", "pull", "origin", target)
		if err := cmd.Run(); err != nil {
			r.logger.Warn(fmt.Sprintf("Failed to pull latest changes: %v", err))
		}
	}

	return nil
}

func (r *ToolRunner) executeBuildSteps(name, repoDir string, buildSteps []string) error {
	if len(buildSteps) == 0 {
		r.logger.Debug(fmt.Sprintf("No build steps defined for %s", name))
		return nil
	}

	r.logger.Info(fmt.Sprintf("Building %s...", name))

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(repoDir); err != nil {
		return fmt.Errorf("failed to change to repository directory: %w", err)
	}

	// Clean previous builds if make is involved
	if len(buildSteps) > 0 {
		cmd := exec.Command("make", "distclean")
		cmd.Run() // Ignore errors for clean
	}

	// Execute each build step
	for i, step := range buildSteps {
		r.logger.Info(fmt.Sprintf("Executing build step %d/%d: %s", i+1, len(buildSteps), step))

		cmd := exec.Command("sh", "-c", step)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = repoDir

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("build step %d failed: %w", i+1, err)
		}
	}

	return nil
}

func (r *ToolRunner) executeInstallSteps(name, repoDir string, installSteps []string) error {
	if len(installSteps) == 0 {
		r.logger.Debug(fmt.Sprintf("No install steps defined for %s", name))
		return nil
	}

	r.logger.Info(fmt.Sprintf("Installing %s...", name))

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(repoDir); err != nil {
		return fmt.Errorf("failed to change to repository directory: %w", err)
	}

	// Execute each install step
	for i, step := range installSteps {
		r.logger.Info(fmt.Sprintf("Executing install step %d/%d: %s", i+1, len(installSteps), step))

		cmd := exec.Command("sh", "-c", step)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin // For sudo password prompts
		cmd.Dir = repoDir

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("install step %d failed: %w", i+1, err)
		}
	}

	return nil
}
