package configurator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lukeberry99/devtool/internal/config"
	"github.com/lukeberry99/devtool/internal/ui"
)

type DotfilesManager struct {
	config        *config.Config
	logger        *ui.Logger
	dryRun        bool
	backupManager *BackupManager
	filesystem    *FilesystemManager
}

func NewDotfilesManager(cfg *config.Config, logger *ui.Logger, dryRun bool) (*DotfilesManager, error) {
	if cfg.Dotfiles.SourceRoot == "" {
		return nil, fmt.Errorf("dotfiles.source_root is required in configuration")
	}

	// Validate source root exists
	sourceRoot := expandPath(cfg.Dotfiles.SourceRoot)
	if _, err := os.Stat(sourceRoot); err != nil {
		return nil, fmt.Errorf("source_root does not exist: %s", sourceRoot)
	}

	backupManager := NewBackupManager(cfg.Dotfiles.BackupDir, logger, dryRun)
	filesystem := NewFilesystemManager(logger, dryRun)

	return &DotfilesManager{
		config:        cfg,
		logger:        logger,
		dryRun:        dryRun,
		backupManager: backupManager,
		filesystem:    filesystem,
	}, nil
}

func (d *DotfilesManager) Deploy() error {
	d.logger.Info("Deploying configuration files...")

	if d.dryRun {
		d.logger.Info("DRY RUN MODE: No actual files will be modified")
	}

	// Validate all source paths exist before deployment
	if err := d.validateSourcePaths(); err != nil {
		return err
	}

	// Deploy each mapping
	for source, target := range d.config.Dotfiles.Mappings {
		if err := d.deployPath(source, target); err != nil {
			return fmt.Errorf("failed to deploy %s: %w", source, err)
		}
	}

	d.logger.Info("Configuration deployment completed successfully")
	return nil
}

func (d *DotfilesManager) validateSourcePaths() error {
	d.logger.Debug("Validating source paths...")

	for source := range d.config.Dotfiles.Mappings {
		sourcePath, err := d.getSourcePath(source)
		if err != nil {
			return err
		}

		if _, err := os.Stat(sourcePath); err != nil {
			return fmt.Errorf("source path does not exist: %s", sourcePath)
		}
	}

	return nil
}

func (d *DotfilesManager) deployPath(source, target string) error {
	sourcePath, err := d.getSourcePath(source)
	if err != nil {
		return err
	}

	targetPath := expandPath(target)

	d.logger.Info(fmt.Sprintf("Copying files from: %s to %s", sourcePath, targetPath))

	// Check if source is directory or file
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to stat source path: %w", err)
	}

	if sourceInfo.IsDir() {
		return d.deployDirectory(sourcePath, targetPath)
	} else {
		return d.deployFile(sourcePath, targetPath)
	}
}

func (d *DotfilesManager) deployDirectory(sourcePath, targetPath string) error {
	// Create backup if target exists
	if err := d.backupManager.CreateBackup(targetPath); err != nil {
		return err
	}

	// Remove existing target (destructive like bash script)
	if err := d.filesystem.RemoveExisting(targetPath); err != nil {
		return err
	}

	// Copy directory
	return d.filesystem.CopyDirectory(sourcePath, targetPath)
}

func (d *DotfilesManager) deployFile(sourcePath, targetPath string) error {
	// Create backup if target exists
	if err := d.backupManager.CreateBackup(targetPath); err != nil {
		return err
	}

	// Remove existing target
	if err := d.filesystem.RemoveExisting(targetPath); err != nil {
		return err
	}

	// Copy file
	return d.filesystem.CopyFile(sourcePath, targetPath)
}

func (d *DotfilesManager) getSourcePath(mapping string) (string, error) {
	sourceRoot := expandPath(d.config.Dotfiles.SourceRoot)
	return filepath.Join(sourceRoot, mapping), nil
}
