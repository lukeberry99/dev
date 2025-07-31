package configurator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lukeberry99/devtool/internal/ui"
)

type FilesystemManager struct {
	logger *ui.Logger
	dryRun bool
}

func NewFilesystemManager(logger *ui.Logger, dryRun bool) *FilesystemManager {
	return &FilesystemManager{
		logger: logger,
		dryRun: dryRun,
	}
}

func (f *FilesystemManager) RemoveExisting(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Path doesn't exist, nothing to remove
		return nil
	}

	f.logger.Debug(fmt.Sprintf("Removing: %s", path))

	if f.dryRun {
		f.logger.Info(fmt.Sprintf("[DRY RUN] Would remove: %s", path))
		return nil
	}

	return os.RemoveAll(path)
}

func (f *FilesystemManager) CopyFile(source, target string) error {
	f.logger.Info(fmt.Sprintf("Copying: %s to %s", source, target))

	if f.dryRun {
		f.logger.Info(fmt.Sprintf("[DRY RUN] Would copy: %s to %s", source, target))
		return nil
	}

	// Create target directory if it doesn't exist
	targetDir := filepath.Dir(target)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Copy file permissions
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	return os.Chmod(target, sourceInfo.Mode())
}

func (f *FilesystemManager) CopyDirectory(source, target string) error {
	f.logger.Info(fmt.Sprintf("Copying directory: %s to %s", source, target))

	if f.dryRun {
		f.logger.Info(fmt.Sprintf("[DRY RUN] Would copy directory: %s to %s", source, target))
		return nil
	}

	// Create target directory
	if err := os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Walk source directory and copy each file/directory
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source directory itself
		if path == source {
			return nil
		}

		// Calculate relative path from source
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		targetPath := filepath.Join(target, relPath)

		if info.IsDir() {
			f.logger.Debug(fmt.Sprintf("Creating directory: %s", targetPath))
			return os.MkdirAll(targetPath, info.Mode())
		} else {
			f.logger.Debug(fmt.Sprintf("Copying file: %s -> %s", path, targetPath))
			return f.copyFileContent(path, targetPath, info.Mode())
		}
	})
}

func (f *FilesystemManager) copyFileContent(source, target string, mode os.FileMode) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create target directory if needed
	targetDir := filepath.Dir(target)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	targetFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return os.Chmod(target, mode)
}
