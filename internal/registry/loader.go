package registry

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader loads registry data from various sources.
type Loader struct {
	basePath string
}

// NewLoader creates a new Loader with the given base path.
func NewLoader(basePath string) *Loader {
	return &Loader{basePath: basePath}
}

// LoadRegistry loads the registry index from registry.yaml.
func (l *Loader) LoadRegistry() (*Registry, error) {
	path := filepath.Join(l.basePath, "registry.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry.yaml: %w", err)
	}

	var reg Registry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("failed to parse registry.yaml: %w", err)
	}

	return &reg, nil
}

// LoadPackage loads a package definition from pkgs/{name}.yaml.
func (l *Loader) LoadPackage(name string) (*Package, error) {
	path := filepath.Join(l.basePath, "pkgs", name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read package %s: %w", name, err)
	}

	var pkg Package
	if err := yaml.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse package %s: %w", name, err)
	}

	return &pkg, nil
}

// LoadAllPackages loads all packages listed in the registry.
func (l *Loader) LoadAllPackages() ([]*Package, error) {
	reg, err := l.LoadRegistry()
	if err != nil {
		return nil, err
	}

	packages := make([]*Package, 0, len(reg.Packages))
	for _, name := range reg.Packages {
		pkg, err := l.LoadPackage(name)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

// Search searches packages by name or tag.
func (l *Loader) Search(query string) ([]*Package, error) {
	packages, err := l.LoadAllPackages()
	if err != nil {
		return nil, err
	}

	var results []*Package
	for _, pkg := range packages {
		if matchesQuery(pkg, query) {
			results = append(results, pkg)
		}
	}

	return results, nil
}

func matchesQuery(pkg *Package, query string) bool {
	if contains(pkg.Name, query) || contains(pkg.DisplayName, query) {
		return true
	}
	for _, tag := range pkg.Tags {
		if contains(tag, query) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchString(s, substr)))
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
