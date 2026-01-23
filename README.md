# crx - Chrome Extension Manager

Declarative Chrome extension manager using Chrome Enterprise Policy.

## Overview

crx allows you to manage Chrome extensions across all profiles using a YAML configuration file. Extensions are installed via Chrome Enterprise Policy, ensuring they are applied to all browser profiles automatically.

## Features

- Declarative configuration (YAML)
- Registry-based extension management (inspired by [aqua](https://github.com/aquaproj/aqua))
- Cross-platform support (macOS, Linux, Windows)
- Interactive TUI for browsing extensions
- Force-install extensions across all Chrome profiles

## Installation

```bash
go install github.com/user/crx/cmd/crx@latest
```

Or build from source:

```bash
git clone https://github.com/user/crx.git
cd crx
go build ./cmd/crx
```

## Quick Start

```bash
# Initialize configuration
crx init

# Add an extension
crx add vimium

# Apply the policy
crx apply
```

## Commands

| Command | Description |
|---------|-------------|
| `crx init` | Initialize configuration file |
| `crx add <name>` | Add an extension to configuration |
| `crx list` | List configured extensions |
| `crx browse` | Interactive TUI to browse and select extensions |
| `crx apply` | Apply the policy to the system |
| `crx apply --dry-run` | Show policy without applying |

## Configuration

Configuration file location: `~/.config/crx/config.yaml`

```yaml
registries:
  - name: standard
    type: github
    repo: user/crx-registry
    ref: main

extensions:
  - vimium
  - dark-reader

settings:
  mode: force_install  # force_install, normal_install, or allowed
```

### Installation Modes

| Mode | Description |
|------|-------------|
| `force_install` | Extensions are installed automatically and cannot be removed by users |
| `normal_install` | Extensions are installed but users can disable/remove them |
| `allowed` | Extensions are allowed but not automatically installed |

## OS-Specific Notes

### macOS

**Requirements:**
- macOS Big Sur (11.0) or later
- Manual profile installation required (Apple security restriction)

**How it works:**
1. `crx apply` generates a `.mobileconfig` profile on your Desktop
2. System Settings opens automatically
3. Click "Install" to apply the profile
4. Reload policies at `chrome://policy` or restart Chrome

**Apply policy:**
```bash
crx apply
# Follow the System Settings prompt to install the profile
```

**Remove policy:**
```
System Settings > Privacy & Security > Profiles > crx Chrome Extensions > Remove
```

Or via terminal:
```bash
sudo profiles -R -p com.crx.chrome.extensions
```

**Limitation:** Every time you add/remove extensions, you need to manually approve the profile in System Settings. This is an Apple security requirement and cannot be bypassed without MDM.

### Linux

**Requirements:**
- Root access (sudo)

**How it works:**
- Writes JSON policy file to `/etc/opt/chrome/policies/managed/`
- Chrome reads this directory on startup

**Apply policy:**
```bash
sudo crx apply
```

**Remove policy:**
```bash
sudo rm /etc/opt/chrome/policies/managed/crx-extensions.json
```

**Note:** Fully automatic, no user interaction required after sudo.

### Windows

**Status:** Not yet implemented

**Manual workaround:**
Add registry entries to:
```
HKLM\SOFTWARE\Policies\Google\Chrome\ExtensionInstallForcelist
```

Each extension should be added as a string value:
- Name: `1`, `2`, `3`, etc.
- Value: `<extension_id>;https://clients2.google.com/service/update2/crx`

## Registry

crx uses a registry to map extension names to Chrome Web Store IDs. This allows you to use human-readable names instead of 32-character extension IDs.

### Registry Structure

```
crx-registry/
├── registry.yaml      # Index of all packages
├── schema.md          # Schema documentation
└── pkgs/
    ├── vimium.yaml
    ├── dark-reader.yaml
    └── ...
```

### Package Format

```yaml
name: vimium
id: dbepggeogbaibhgnhhndojpepiihcmeb
display_name: Vimium
description: The Hacker's Browser
homepage: https://vimium.github.io/
repository: https://github.com/philc/vimium
tags:
  - productivity
  - keyboard
  - vim
```

### Using a Local Registry

For testing or private extensions:

```bash
crx add my-extension --registry /path/to/local/registry
crx apply --registry /path/to/local/registry
```

## Extension Updates

Chrome automatically updates force-installed extensions using the `update_url` specified in the policy (Chrome Web Store). No action required from crx.

## Troubleshooting

### Check Applied Policies

Open `chrome://policy` in Chrome to see all applied policies.

### Verify Extension Installation

Open `chrome://extensions` to see installed extensions. Force-installed extensions will not have a "Remove" button.

### macOS: Profile Not Working

1. Check `chrome://policy` - ensure the policy shows as "Mandatory" not "Recommended"
2. Verify the profile is installed: System Settings > Privacy & Security > Profiles
3. Reload policies at `chrome://policy` or restart Chrome completely

### Linux: Permission Denied

Ensure you're running with sudo:
```bash
sudo crx apply
```

### Extension Not Installing

1. Verify the extension ID is correct
2. Check if the extension is available in Chrome Web Store
3. Some extensions (like uBlock Origin) have been removed due to Manifest V3 - use alternatives like uBlock Origin Lite

## License

MIT
