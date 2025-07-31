package configurator

import (
	"os"
	"path/filepath"
	"strings"
)

func expandPath(path string) string {
	// Expand environment variables
	path = os.ExpandEnv(path)

	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, path[2:])
	}

	return path
}
