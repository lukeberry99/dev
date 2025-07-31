package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LocalStateManager struct {
	statePath string
	state     *LocalState
}

func NewLocalStateManager() (*LocalStateManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	statePath := filepath.Join(homeDir, ".devtool", "state.json")

	manager := &LocalStateManager{
		statePath: statePath,
	}

	if err := manager.load(); err != nil {
		// If loading fails, create new state with machine info
		hostname, _ := os.Hostname()
		manager.state = &LocalState{
			MachineID:     generateMachineID(),
			Version:       "1.0",
			LastUpdated:   time.Now(),
			LastSync:      time.Time{},
			Hostname:      hostname,
			OS:            "darwin", // TODO: detect actual OS
			Arch:          "arm64",  // TODO: detect actual arch
			Tools:         make(map[string]ToolStatus),
			ActiveProfile: "default",
			Preferences: MachinePreferences{
				AutoUpdate:            true,
				BuildNeovimFromSource: false,
				BackupBeforeChanges:   true,
			},
			LastBackup: time.Time{},
		}
	}

	return manager, nil
}

func generateMachineID() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s-%d", hostname, time.Now().Unix())
}
func (m *LocalStateManager) load() error {
	data, err := os.ReadFile(m.statePath)
	if err != nil {
		return err
	}

	var state LocalState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	m.state = &state
	return nil
}

func (m *LocalStateManager) Save() error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(m.statePath), 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	m.state.LastUpdated = time.Now()

	data, err := json.MarshalIndent(m.state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(m.statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

func (m *LocalStateManager) IsToolCurrent(name, expectedVersion string) bool {
	status, exists := m.state.Tools[name]
	if !exists {
		return false
	}

	return status.Installed && status.Version == expectedVersion
}

func (m *LocalStateManager) UpdateToolStatus(name string, status ToolStatus) {
	status.LastChecked = time.Now()
	m.state.Tools[name] = status
}

func (m *LocalStateManager) GetToolStatus(name string) (ToolStatus, bool) {
	status, exists := m.state.Tools[name]
	return status, exists
}

func (m *LocalStateManager) GetAllTools() map[string]ToolStatus {
	return m.state.Tools
}
