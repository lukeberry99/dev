package types

type InstallationMethod int

const (
	HomebrewInstall InstallationMethod = iota
	BuildFromSource
	CustomScript
	ManualInstall
)

func (i InstallationMethod) String() string {
	switch i {
	case HomebrewInstall:
		return "homebrew"
	case BuildFromSource:
		return "build"
	case CustomScript:
		return "script"
	case ManualInstall:
		return "manual"
	default:
		return "unknown"
	}
}

type Tool struct {
	Name            string             `json:"name"`
	DisplayName     string             `json:"display_name"`
	Description     string             `json:"description"`
	Method          InstallationMethod `json:"installation_method"`
	Version         string             `json:"version"`
	Dependencies    []string           `json:"dependencies"`
	HomebrewPackage string             `json:"homebrew_package,omitempty"`
	BuildConfig     *BuildInfo         `json:"build_config,omitempty"`
	CheckCommand    string             `json:"check_command"`
	VersionCommand  string             `json:"version_command"`
}

type BuildInfo struct {
	Repository   string   `json:"repository"`
	Branch       string   `json:"branch"`
	BuildDir     string   `json:"build_dir"`
	BuildSteps   []string `json:"build_steps"`
	InstallSteps []string `json:"install_steps"`
	CleanupSteps []string `json:"cleanup_steps"`
}
