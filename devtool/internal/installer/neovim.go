package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lukeberry99/devtool/internal/ui"
)

// Replicate exact Neovim build process from runs/neovim
type NeovimBuilder struct {
	logger   *ui.Logger
	dryRun   bool
	version  string // "nightly", "stable", or specific version
	buildDir string
}

func NewNeovimBuilder(logger *ui.Logger, dryRun bool, version string) *NeovimBuilder {
	homeDir, _ := os.UserHomeDir()
	buildDir := filepath.Join(homeDir, ".devtool", "builds", "neovim")

	return &NeovimBuilder{
		logger:   logger,
		dryRun:   dryRun,
		version:  version,
		buildDir: buildDir,
	}
}

func (n *NeovimBuilder) Build() error {
	n.logger.Info(fmt.Sprintf("Installing Neovim version: %s", n.version))

	// 1. Install dependencies via Homebrew
	if err := n.installDependencies(); err != nil {
		return err
	}

	// 2. Clone or update repository
	if err := n.prepareRepository(); err != nil {
		return err
	}

	// 3. Build Neovim
	if err := n.build(); err != nil {
		return err
	}

	// 4. Install Neovim
	if err := n.install(); err != nil {
		return err
	}

	n.logger.Info("✅ Neovim installation completed")
	return nil
}

func (n *NeovimBuilder) installDependencies() error {
	n.logger.Info("Installing Neovim build dependencies...")

	dependencies := []string{"cmake", "gettext", "lua"}

	if n.dryRun {
		n.logger.Info(fmt.Sprintf("[DRY RUN] Would install dependencies: %v", dependencies))
		return nil
	}

	homebrew := NewHomebrewManager(n.logger, n.dryRun)
	if err := homebrew.EnsureInstalled(); err != nil {
		return fmt.Errorf("failed to ensure Homebrew is installed: %w", err)
	}

	return homebrew.InstallPackages(dependencies)
}

func (n *NeovimBuilder) prepareRepository() error {
	n.logger.Info("Preparing Neovim repository...")

	if n.dryRun {
		n.logger.Info("[DRY RUN] Would clone/update Neovim repository")
		return nil
	}

	// Ensure build directory exists
	if err := os.MkdirAll(n.buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	repoPath := filepath.Join(n.buildDir, "neovim")

	// Check if repository already exists
	if _, err := os.Stat(filepath.Join(repoPath, ".git")); err == nil {
		n.logger.Debug("Repository exists, updating...")
		return n.updateRepository(repoPath)
	}

	// Clone repository
	n.logger.Info("Cloning Neovim repository...")
	cmd := exec.Command("git", "clone", "https://github.com/neovim/neovim.git", repoPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone Neovim repository: %w", err)
	}

	// Checkout appropriate branch/tag
	return n.checkoutVersion(repoPath)
}

func (n *NeovimBuilder) updateRepository(repoPath string) error {
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

	return n.checkoutVersion(repoPath)
}

func (n *NeovimBuilder) checkoutVersion(repoPath string) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(repoPath); err != nil {
		return fmt.Errorf("failed to change to repository directory: %w", err)
	}

	var target string
	switch n.version {
	case "nightly", "":
		target = "master"
	case "stable":
		target = "stable"
	default:
		target = n.version
	}

	n.logger.Debug(fmt.Sprintf("Checking out %s", target))
	cmd := exec.Command("git", "checkout", target)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout %s: %w", target, err)
	}

	// Pull latest changes if on a branch
	if target == "master" || target == "stable" {
		cmd = exec.Command("git", "pull", "origin", target)
		if err := cmd.Run(); err != nil {
			n.logger.Warn(fmt.Sprintf("Failed to pull latest changes: %v", err))
		}
	}

	return nil
}

func (n *NeovimBuilder) build() error {
	n.logger.Info("Building Neovim...")

	if n.dryRun {
		n.logger.Info("[DRY RUN] Would build Neovim")
		return nil
	}

	repoPath := filepath.Join(n.buildDir, "neovim")

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(repoPath); err != nil {
		return fmt.Errorf("failed to change to repository directory: %w", err)
	}

	// Clean previous builds
	cmd := exec.Command("make", "distclean")
	cmd.Run() // Ignore errors for clean

	// Build with RelWithDebInfo configuration
	cmd = exec.Command("make", "CMAKE_BUILD_TYPE=RelWithDebInfo")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build Neovim: %w", err)
	}

	return nil
}

func (n *NeovimBuilder) install() error {
	n.logger.Info("Installing Neovim...")

	if n.dryRun {
		n.logger.Info("[DRY RUN] Would install Neovim")
		return nil
	}

	repoPath := filepath.Join(n.buildDir, "neovim")

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(repoPath); err != nil {
		return fmt.Errorf("failed to change to repository directory: %w", err)
	}

	// Install Neovim
	cmd := exec.Command("sudo", "make", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin // For sudo password prompt

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Neovim: %w", err)
	}

	return nil
}
