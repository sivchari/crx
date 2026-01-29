//go:build !windows

package policy

import "fmt"

// applyWindows is a stub for non-Windows platforms.
func (g *Generator) applyWindows(policy *Policy) error {
	return fmt.Errorf("windows support is only available on Windows")
}

// RemoveWindowsPolicy is a stub for non-Windows platforms.
func RemoveWindowsPolicy() error {
	return fmt.Errorf("windows support is only available on Windows")
}
