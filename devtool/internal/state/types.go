package state

import "time"

// Enhanced LocalState to match rewrite.md vision
type LocalState struct {
	MachineID     string                `json:"machine_id"`
	Version       string                `json:"version"`
	LastUpdated   time.Time             `json:"last_updated"`
	LastSync      time.Time             `json:"last_sync"`
	Hostname      string                `json:"hostname"`
	OS            string                `json:"os"`
	Arch          string                `json:"arch"`
	Tools         map[string]ToolStatus `json:"installed_tools"`
	ActiveProfile string                `json:"active_profile"`
	Preferences   MachinePreferences    `json:"preferences"`
	LastBackup    time.Time             `json:"last_backup"`
}

// Enhanced ToolStatus to match rewrite.md structure
type ToolStatus struct {
	Installed     bool      `json:"installed"`
	Version       string    `json:"version"`
	InstallDate   time.Time `json:"installed_at"`
	LastChecked   time.Time `json:"last_checked"`
	Source        string    `json:"source"` // "homebrew", "built_from_source", "manual"
	BinaryPath    string    `json:"binary_path"`
	ConfigCurrent bool      `json:"config_current"`
}

// Machine preferences as outlined in rewrite.md
type MachinePreferences struct {
	AutoUpdate            bool `json:"auto_update"`
	BuildNeovimFromSource bool `json:"build_neovim_from_source"`
	BackupBeforeChanges   bool `json:"backup_before_changes"`
}

type InstallationState int

const (
	NotInstalled InstallationState = iota
	Installed
	OutOfDate
	Failed
)
