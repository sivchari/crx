package registry

// Package represents a Chrome extension package in the registry.
type Package struct {
	Name        string   `yaml:"name"`
	ID          string   `yaml:"id"`
	DisplayName string   `yaml:"display_name"`
	Description string   `yaml:"description,omitempty"`
	Homepage    string   `yaml:"homepage,omitempty"`
	Repository  string   `yaml:"repository,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
}

// Registry represents the registry index.
type Registry struct {
	Version  int      `yaml:"version"`
	Packages []string `yaml:"packages"`
}

// CRXUpdateURL is the Chrome Web Store update URL.
const CRXUpdateURL = "https://clients2.google.com/service/update2/crx"
