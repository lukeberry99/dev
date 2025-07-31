package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  string                `yaml:"version"`
	Tools    map[string]ToolConfig `yaml:"tools"`
	Profiles map[string]Profile    `yaml:"profiles"`
	Dotfiles DotfilesConfig        `yaml:"dotfiles"`
	Homebrew HomebrewConfig        `yaml:"homebrew"`
	Logging  LoggingConfig         `yaml:"logging"`
	Sync     SyncConfig            `yaml:"sync"`
}

type ToolConfig struct {
	Version      string       `yaml:"version"`
	Source       string       `yaml:"source"` // "homebrew", "build", "script"
	Dependencies []string     `yaml:"dependencies"`
	BuildConfig  *BuildConfig `yaml:"build_config,omitempty"`
	HomebrewArgs []string     `yaml:"homebrew_args,omitempty"`
	Profile      []string     `yaml:"profile"`
	Enabled      bool         `yaml:"enabled"`
}

type BuildConfig struct {
	Repository   string   `yaml:"repository"`
	BuildSteps   []string `yaml:"build_steps"`
	InstallSteps []string `yaml:"install_steps"`
	Dependencies []string `yaml:"dependencies"`
}

type DotfilesConfig struct {
	BackupDir string            `yaml:"backup_dir"`
	Strategy  string            `yaml:"strategy"` // "copy", "symlink"
	Mappings  map[string]string `yaml:"mappings"`
}

type HomebrewConfig struct {
	AutoUpdate   bool `yaml:"auto_update"`
	CleanupAfter bool `yaml:"cleanup_after"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Profile struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Include     []string `yaml:"include"`
	Exclude     []string `yaml:"exclude"`
}

type SyncConfig struct {
	Strategy string      `yaml:"strategy"` // "git", "cloud", "hybrid"
	Git      GitConfig   `yaml:"git"`
	Cloud    CloudConfig `yaml:"cloud"`
}

type GitConfig struct {
	Repository string `yaml:"repository"`
	Branch     string `yaml:"branch"`
	AuthType   string `yaml:"auth_type"` // "ssh", "token"
}

type CloudConfig struct {
	Provider string `yaml:"provider"` // "s3", "gcs"
	Bucket   string `yaml:"bucket"`
	Region   string `yaml:"region"`
	Prefix   string `yaml:"prefix"`
}

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func getDefaultConfigPath() string {
	if repoRoot := findRepoRoot(); repoRoot != "" {
		return filepath.Join(repoRoot, "devtool", "configs", "devtool.yml")
	}

	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".devtool.yaml")
}

func findRepoRoot() string {
	dir, _ := os.Getwd()
	for dir != "/" {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	return ""
}
