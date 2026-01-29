//go:build windows

package policy

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

const (
	chromePolicyKey = `SOFTWARE\Policies\Google\Chrome`
	forcelistKey    = `ExtensionInstallForcelist`
	allowlistKey    = `ExtensionInstallAllowlist`
	blocklistKey    = `ExtensionInstallBlocklist`
)

// applyWindows applies the policy on Windows using registry.
func (g *Generator) applyWindows(policy *Policy) error {
	// Open or create the Chrome policy key
	chromeKey, _, err := registry.CreateKey(registry.LOCAL_MACHINE, chromePolicyKey, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open Chrome policy key (try running as Administrator): %w", err)
	}
	defer chromeKey.Close()

	// Apply ExtensionInstallForcelist
	if len(policy.ExtensionInstallForcelist) > 0 {
		if err := writeStringList(chromeKey, forcelistKey, policy.ExtensionInstallForcelist); err != nil {
			return fmt.Errorf("failed to write ExtensionInstallForcelist: %w", err)
		}
	}

	// Apply ExtensionInstallAllowlist
	if len(policy.ExtensionInstallAllowlist) > 0 {
		if err := writeStringList(chromeKey, allowlistKey, policy.ExtensionInstallAllowlist); err != nil {
			return fmt.Errorf("failed to write ExtensionInstallAllowlist: %w", err)
		}
	}

	// Apply ExtensionInstallBlocklist
	if len(policy.ExtensionInstallBlocklist) > 0 {
		if err := writeStringList(chromeKey, blocklistKey, policy.ExtensionInstallBlocklist); err != nil {
			return fmt.Errorf("failed to write ExtensionInstallBlocklist: %w", err)
		}
	}

	return nil
}

// writeStringList writes a list of strings to a registry subkey.
// Each entry is stored as a numbered value (1, 2, 3, ...).
func writeStringList(parentKey registry.Key, subkeyName string, values []string) error {
	// Delete existing subkey if it exists
	_ = registry.DeleteKey(parentKey, subkeyName)

	// Create new subkey
	subkey, _, err := registry.CreateKey(parentKey, subkeyName, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer subkey.Close()

	// Write each value with a numbered name
	for i, value := range values {
		name := fmt.Sprintf("%d", i+1)
		if err := subkey.SetStringValue(name, value); err != nil {
			return err
		}
	}

	return nil
}

// RemoveWindowsPolicy removes the crx Chrome policy from Windows registry.
func RemoveWindowsPolicy() error {
	chromeKey, err := registry.OpenKey(registry.LOCAL_MACHINE, chromePolicyKey, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open Chrome policy key: %w", err)
	}
	defer chromeKey.Close()

	// Delete policy subkeys
	_ = registry.DeleteKey(chromeKey, forcelistKey)
	_ = registry.DeleteKey(chromeKey, allowlistKey)
	_ = registry.DeleteKey(chromeKey, blocklistKey)

	return nil
}
