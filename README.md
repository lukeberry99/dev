# DevTool

Development environment setup tool for macOS. Installs tools via Homebrew and custom builds, manages dotfiles.

- [My dotfiles](https://github.com/lukeberry99/dotfiles)

## Usage

```bash
# Build
go build -o devtool .

# Install tools
./devtool install

# Deploy dotfiles 
./devtool configure

# Check status
./devtool status
```

## Configuration

Uses `$HOME/.devtool.yaml` or specify with `--config`.

```yaml
tools:
  go:
    source: "homebrew"
    enabled: true

  aerospace:
    source: "homebrew"
    cask: true
    app_name: "Aerospace"
    enabled: true
    
  neovim:
    source: "build"
    build_config:
      repository: "https://github.com/neovim/neovim.git"
      build_steps:
        - "make CMAKE_BUILD_TYPE=RelWithDebInfo"
    enabled: true

dotfiles:
  mappings:
    "env/.config": "~/.config"
    "env/.zshrc": "~/.zshrc"
```

## Options

- `--dry-run`: Preview without executing
- `--verbose`: Detailed output
- `--force`: Reinstall existing tools
