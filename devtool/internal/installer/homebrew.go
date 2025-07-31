package installer

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/lukeberry99/devtool/internal/ui"
)

type HomebrewManager struct {
	logger *ui.Logger
	dryRun bool
}

func NewHomebrewManager(logger *ui.Logger, dryRun bool) *HomebrewManager {
	return &HomebrewManager{
		logger: logger,
		dryRun: dryRun,
	}
}

func (h *HomebrewManager) EnsureInstalled() error {
	// Replicate exact Homebrew installation logic from bash
	if h.isInstalled() {
		h.logger.Info("Homebrew is already installed")
		return h.update()
	}

	h.logger.Info("Installing Homebrew...")
	if h.dryRun {
		h.logger.Info("[DRY RUN] Would install Homebrew")
		return nil
	}

	// Execute Homebrew installation script
	return h.install()
}

func (h *HomebrewManager) isInstalled() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func (h *HomebrewManager) install() error {
	// Check if we're on macOS (required for Homebrew)
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("homebrew installation is only supported on macOS")
	}

	h.logger.Info("Downloading and installing Homebrew...")

	// Use the official Homebrew installation script
	cmd := exec.Command("/bin/bash", "-c",
		`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Homebrew: %w", err)
	}

	h.logger.Info("✅ Homebrew installation completed")
	return nil
}

func (h *HomebrewManager) update() error {
	if h.dryRun {
		h.logger.Info("[DRY RUN] Would update Homebrew")
		return nil
	}

	h.logger.Debug("Updating Homebrew...")
	cmd := exec.Command("brew", "update")

	if err := cmd.Run(); err != nil {
		h.logger.Warn("Failed to update Homebrew, continuing anyway")
		return nil // Don't fail on update errors
	}

	return nil
}

func (h *HomebrewManager) InstallPackages(packages []string) error {
	for _, pkg := range packages {
		if h.isPackageInstalled(pkg) {
			h.logger.Debug(fmt.Sprintf("Package %s already installed", pkg))
			continue
		}

		h.logger.Info(fmt.Sprintf("Installing %s...", pkg))
		if err := h.installPackage(pkg); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}
	return nil
}

func (h *HomebrewManager) InstallPackageWithArgs(pkg string, args []string) error {
	if h.isPackageInstalled(pkg) {
		h.logger.Debug(fmt.Sprintf("Package %s already installed", pkg))
		return nil
	}

	h.logger.Info(fmt.Sprintf("Installing %s with args %v...", pkg, args))
	if err := h.installPackageWithArgs(pkg, args); err != nil {
		return fmt.Errorf("failed to install %s: %w", pkg, err)
	}
	return nil
}

func (h *HomebrewManager) isPackageInstalled(pkg string) bool {
	if h.dryRun {
		return false // In dry-run mode, assume nothing is installed
	}

	// Check if package is installed by trying to find it in brew list
	cmd := exec.Command("brew", "list", "--formula", pkg)
	if err := cmd.Run(); err == nil {
		return true
	}

	// Also check casks
	cmd = exec.Command("brew", "list", "--cask", pkg)
	return cmd.Run() == nil
}

func (h *HomebrewManager) installPackage(pkg string) error {
	if h.dryRun {
		h.logger.Info(fmt.Sprintf("[DRY RUN] Would install package: %s", pkg))
		return nil
	}

	cmd := exec.Command("brew", "install", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("brew install %s failed: %w", pkg, err)
	}

	return nil
}

func (h *HomebrewManager) installPackageWithArgs(pkg string, args []string) error {
	if h.dryRun {
		h.logger.Info(fmt.Sprintf("[DRY RUN] Would install package: %s with args %v", pkg, args))
		return nil
	}

	// Build command: brew install [args...] [package]
	cmdArgs := append([]string{"install"}, args...)
	cmdArgs = append(cmdArgs, pkg)

	cmd := exec.Command("brew", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("brew install %s %s failed: %w", strings.Join(args, " "), pkg, err)
	}

	return nil
}

func (h *HomebrewManager) GetInstalledVersion(pkg string) (string, error) {
	if h.dryRun {
		return "dry-run-version", nil
	}

	// Try to get version from brew info
	cmd := exec.Command("brew", "info", pkg, "--json")
	_, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get info for %s: %w", pkg, err)
	}

	// For now, just return that it's installed
	// TODO: Parse JSON to get actual version
	return "installed", nil
}

func (h *HomebrewManager) Cleanup() error {
	if h.dryRun {
		h.logger.Info("[DRY RUN] Would cleanup Homebrew")
		return nil
	}

	h.logger.Debug("Cleaning up Homebrew...")
	cmd := exec.Command("brew", "cleanup")

	if err := cmd.Run(); err != nil {
		h.logger.Warn("Failed to cleanup Homebrew, continuing anyway")
		return nil // Don't fail on cleanup errors
	}

	return nil
}
