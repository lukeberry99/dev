package state

import (
	"os/exec"
	"strings"
)

type ToolDetector struct{}

func NewToolDetector() *ToolDetector {
	return &ToolDetector{}
}

func (d *ToolDetector) IsInstalled(toolName string) bool {
	_, err := exec.LookPath(toolName)
	return err == nil
}

func (d *ToolDetector) GetVersion(toolName string) (string, error) {
	switch toolName {
	case "go":
		return d.getGoVersion()
	case "node":
		return d.getNodeVersion()
	case "nvim":
		return d.getNeovimVersion()
	case "git":
		return d.getGitVersion()
	default:
		return d.getGenericVersion(toolName)
	}
}

func (d *ToolDetector) getGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse "go version go1.21.0 darwin/arm64"
	parts := strings.Fields(string(output))
	if len(parts) >= 3 {
		return strings.TrimPrefix(parts[2], "go"), nil
	}

	return string(output), nil
}

func (d *ToolDetector) getNodeVersion() (string, error) {
	cmd := exec.Command("node", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimPrefix(string(output), "v")), nil
}

func (d *ToolDetector) getNeovimVersion() (string, error) {
	cmd := exec.Command("nvim", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		// Parse "NVIM v0.9.1"
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			return strings.TrimPrefix(parts[1], "v"), nil
		}
	}

	return string(output), nil
}

func (d *ToolDetector) getGitVersion() (string, error) {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse "git version 2.39.2"
	parts := strings.Fields(string(output))
	if len(parts) >= 3 {
		return parts[2], nil
	}

	return string(output), nil
}

func (d *ToolDetector) getGenericVersion(toolName string) (string, error) {
	// Try common version flags
	versionFlags := []string{"--version", "-v", "version"}

	for _, flag := range versionFlags {
		cmd := exec.Command(toolName, flag)
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output)), nil
		}
	}

	return "unknown", nil
}
