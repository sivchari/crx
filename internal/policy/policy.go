package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/sivchari/crx/internal/config"
	"github.com/sivchari/crx/internal/registry"
)

// Policy represents Chrome Enterprise Policy for extensions.
type Policy struct {
	ExtensionInstallForcelist []string `json:"ExtensionInstallForcelist,omitempty"`
	ExtensionInstallAllowlist []string `json:"ExtensionInstallAllowlist,omitempty"`
	ExtensionInstallBlocklist []string `json:"ExtensionInstallBlocklist,omitempty"`
	ExtensionSettings         any      `json:"ExtensionSettings,omitempty"`
}

// Generator generates Chrome Enterprise Policy JSON.
type Generator struct {
	cfg      *config.Config
	packages []*registry.Package
}

// NewGenerator creates a new Generator.
func NewGenerator(cfg *config.Config, packages []*registry.Package) *Generator {
	return &Generator{
		cfg:      cfg,
		packages: packages,
	}
}

// Generate generates the policy based on the configuration.
func (g *Generator) Generate() (*Policy, error) {
	policy := &Policy{}

	for _, pkg := range g.packages {
		entry := formatEntry(pkg.ID)

		switch g.cfg.Settings.Mode {
		case config.ModeForceInstall:
			policy.ExtensionInstallForcelist = append(policy.ExtensionInstallForcelist, entry)
		case config.ModeNormalInstall:
			// normal_install uses ExtensionSettings
			if policy.ExtensionSettings == nil {
				policy.ExtensionSettings = make(map[string]any)
			}
			settings := policy.ExtensionSettings.(map[string]any)
			settings[pkg.ID] = map[string]any{
				"installation_mode": "normal_installed",
				"update_url":        registry.CRXUpdateURL,
			}
		case config.ModeAllowed:
			policy.ExtensionInstallAllowlist = append(policy.ExtensionInstallAllowlist, pkg.ID)
		}
	}

	return policy, nil
}

// formatEntry formats the extension ID with the update URL.
func formatEntry(id string) string {
	return fmt.Sprintf("%s;%s", id, registry.CRXUpdateURL)
}

// Apply applies the policy to the system using the appropriate method for the OS.
func (g *Generator) Apply(policy *Policy) error {
	switch runtime.GOOS {
	case "darwin":
		return g.applyDarwin(policy)
	case "linux":
		return g.applyLinux(policy)
	case "windows":
		return g.applyWindows(policy)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// applyDarwin applies the policy on macOS using a configuration profile.
// Since macOS Big Sur, profiles cannot be installed silently.
// This generates the profile and opens System Settings for manual installation.
func (g *Generator) applyDarwin(policy *Policy) error {
	// Generate mobileconfig content
	mobileconfig := buildMobileconfig(policy)

	// Determine output path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	profilePath := filepath.Join(homeDir, "Desktop", "crx-chrome-policy.mobileconfig")

	// Write the profile file
	if err := os.WriteFile(profilePath, []byte(mobileconfig), 0644); err != nil {
		return fmt.Errorf("failed to write mobileconfig: %w", err)
	}

	// Open System Settings with the profile
	cmd := exec.Command("open", profilePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open profile: %w", err)
	}

	return nil
}

// buildMobileconfig generates a macOS configuration profile for Chrome policies.
func buildMobileconfig(policy *Policy) string {
	var mcxSettings string

	if len(policy.ExtensionInstallForcelist) > 0 {
		mcxSettings += "\t\t\t\t\t\t\t\t<key>ExtensionInstallForcelist</key>\n\t\t\t\t\t\t\t\t<array>\n"
		for _, entry := range policy.ExtensionInstallForcelist {
			mcxSettings += fmt.Sprintf("\t\t\t\t\t\t\t\t\t<string>%s</string>\n", entry)
		}
		mcxSettings += "\t\t\t\t\t\t\t\t</array>\n"
	}

	if len(policy.ExtensionInstallAllowlist) > 0 {
		mcxSettings += "\t\t\t\t\t\t\t\t<key>ExtensionInstallAllowlist</key>\n\t\t\t\t\t\t\t\t<array>\n"
		for _, entry := range policy.ExtensionInstallAllowlist {
			mcxSettings += fmt.Sprintf("\t\t\t\t\t\t\t\t\t<string>%s</string>\n", entry)
		}
		mcxSettings += "\t\t\t\t\t\t\t\t</array>\n"
	}

	if len(policy.ExtensionInstallBlocklist) > 0 {
		mcxSettings += "\t\t\t\t\t\t\t\t<key>ExtensionInstallBlocklist</key>\n\t\t\t\t\t\t\t\t<array>\n"
		for _, entry := range policy.ExtensionInstallBlocklist {
			mcxSettings += fmt.Sprintf("\t\t\t\t\t\t\t\t\t<string>%s</string>\n", entry)
		}
		mcxSettings += "\t\t\t\t\t\t\t\t</array>\n"
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>PayloadContent</key>
	<array>
		<dict>
			<key>PayloadContent</key>
			<dict>
				<key>com.google.Chrome</key>
				<dict>
					<key>Forced</key>
					<array>
						<dict>
							<key>mcx_preference_settings</key>
							<dict>
%s							</dict>
						</dict>
					</array>
				</dict>
			</dict>
			<key>PayloadEnabled</key>
			<true/>
			<key>PayloadIdentifier</key>
			<string>com.crx.chrome.extensions.inner</string>
			<key>PayloadType</key>
			<string>com.apple.ManagedClient.preferences</string>
			<key>PayloadUUID</key>
			<string>A8B8E6D0-1234-5678-9ABC-DEF012345678</string>
			<key>PayloadVersion</key>
			<integer>1</integer>
		</dict>
	</array>
	<key>PayloadDescription</key>
	<string>Chrome Extension Policy managed by crx</string>
	<key>PayloadDisplayName</key>
	<string>crx Chrome Extensions</string>
	<key>PayloadIdentifier</key>
	<string>com.crx.chrome.extensions</string>
	<key>PayloadOrganization</key>
	<string>crx</string>
	<key>PayloadRemovalDisallowed</key>
	<false/>
	<key>PayloadScope</key>
	<string>System</string>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadUUID</key>
	<string>B9C9F7E1-2345-6789-ABCD-EF0123456789</string>
	<key>PayloadVersion</key>
	<integer>1</integer>
</dict>
</plist>
`, mcxSettings)
}

// applyLinux applies the policy on Linux using JSON file.
func (g *Generator) applyLinux(policy *Policy) error {
	dir := "/etc/opt/chrome/policies/managed"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create policy directory: %w", err)
	}

	path := filepath.Join(dir, "crx-extensions.json")

	data, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write policy file: %w", err)
	}

	return nil
}

// WriteToFile writes the policy to a JSON file (for dry-run or backup).
func (g *Generator) WriteToFile(policy *Policy, filename string) error {
	dir := g.cfg.Settings.PolicyPath
	if dir == "" {
		return fmt.Errorf("policy_path is not set")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create policy directory: %w", err)
	}

	path := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write policy file: %w", err)
	}

	return nil
}

// ToJSON returns the policy as a JSON string.
func (g *Generator) ToJSON(policy *Policy) (string, error) {
	data, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal policy: %w", err)
	}
	return string(data), nil
}

// Diff compares the generated policy with the existing one.
func (g *Generator) Diff(policy *Policy, filename string) (bool, string, error) {
	dir := g.cfg.Settings.PolicyPath
	if dir == "" {
		return false, "", fmt.Errorf("policy_path is not set")
	}

	path := filepath.Join(dir, filename)

	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			newJSON, _ := g.ToJSON(policy)
			return true, fmt.Sprintf("New file will be created:\n%s", newJSON), nil
		}
		return false, "", fmt.Errorf("failed to read existing policy: %w", err)
	}

	newJSON, err := g.ToJSON(policy)
	if err != nil {
		return false, "", err
	}

	if string(existing) == newJSON {
		return false, "No changes", nil
	}

	return true, fmt.Sprintf("Changes detected:\n--- existing\n%s\n+++ new\n%s", string(existing), newJSON), nil
}

// GetPolicyPath returns the appropriate policy path for the current OS.
func GetPolicyPath() string {
	switch runtime.GOOS {
	case "darwin":
		return "Configuration Profile (com.crx.chrome.extensions)"
	case "linux":
		return "/etc/opt/chrome/policies/managed"
	case "windows":
		return "HKLM\\SOFTWARE\\Policies\\Google\\Chrome"
	default:
		return ""
	}
}

// RemoveDarwinProfile removes the crx Chrome policy profile on macOS.
func RemoveDarwinProfile() error {
	cmd := exec.Command("profiles", "-R", "-p", "com.crx.chrome.extensions")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove profile (try with sudo): %w\nOutput: %s", err, string(output))
	}
	return nil
}
