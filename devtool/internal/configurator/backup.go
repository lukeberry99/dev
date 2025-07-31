package configurator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lukeberry99/devtool/internal/ui"
)

type BackupManager struct {
	backupDir string
	logger    *ui.Logger
	dryRun    bool
}

func NewBackupManager(backupDir string, logger *ui.Logger, dryRun bool) *BackupManager {
	return &BackupManager{
		backupDir: backupDir,
		logger:    logger,
		dryRun:    dryRun,
	}
}

func (b *BackupManager) CreateBackup(targetPath string) error {
	// Check if target exists
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		// Target doesn't exist, no backup needed
		return nil
	}

	if b.backupDir == "" {
		b.logger.Debug("No backup directory configured, skipping backup")
		return nil
	}

	backupDir := expandPath(b.backupDir)

	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(backupDir, 0755); err != nil && !b.dryRun {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	basename := filepath.Base(targetPath)
	backupName := fmt.Sprintf("%s.%s.backup", basename, timestamp)
	backupPath := filepath.Join(backupDir, backupName)

	b.logger.Debug(fmt.Sprintf("Creating backup: %s -> %s", targetPath, backupPath))

	if b.dryRun {
		b.logger.Info(fmt.Sprintf("[DRY RUN] Would create backup: %s", backupPath))
		return nil
	}

	// Create backup by copying
	return b.copyToBackup(targetPath, backupPath)
}

func (b *BackupManager) copyToBackup(source, backup string) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("failed to stat source for backup: %w", err)
	}

	if sourceInfo.IsDir() {
		return b.copyDirectoryToBackup(source, backup)
	} else {
		return b.copyFileToBackup(source, backup)
	}
}

func (b *BackupManager) copyFileToBackup(source, backup string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file for backup: %w", err)
	}
	defer sourceFile.Close()

	backupFile, err := os.Create(backup)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer backupFile.Close()

	_, err = sourceFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek source file: %w", err)
	}

	buffer := make([]byte, 1024*1024) // 1MB buffer
	for {
		n, readErr := sourceFile.Read(buffer)
		if n > 0 {
			if _, writeErr := backupFile.Write(buffer[:n]); writeErr != nil {
				return fmt.Errorf("failed to write to backup file: %w", writeErr)
			}
		}
		if readErr != nil {
			break
		}
	}

	return nil
}

func (b *BackupManager) copyDirectoryToBackup(source, backup string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		backupPath := filepath.Join(backup, relPath)

		if info.IsDir() {
			return os.MkdirAll(backupPath, info.Mode())
		} else {
			return b.copyFileToBackup(path, backupPath)
		}
	})
}
