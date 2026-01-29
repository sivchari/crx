package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

// Config represents the user configuration.
type Config struct {
	Extensions []string `yaml:"extensions"`
	Settings   Settings `yaml:"settings"`
}

// Settings represents the application settings.
type Settings struct {
	PolicyPath string `yaml:"policy_path"`
	Mode       string `yaml:"mode"` // "force_install", "normal_install", "allowed"
}

// InstallMode constants.
const (
	ModeForceInstall  = "force_install"
	ModeNormalInstall = "normal_install"
	ModeAllowed       = "allowed"
)

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Extensions: []string{},
		Settings: Settings{
			PolicyPath: defaultPolicyPath(),
			Mode:       ModeForceInstall,
		},
	}
}

func defaultPolicyPath() string {
	switch os := getOS(); os {
	case "darwin":
		return "/Library/Google/Chrome/policies/managed"
	case "linux":
		return "/etc/opt/chrome/policies/managed"
	default:
		return ""
	}
}

func getOS() string {
	// Simplified OS detection
	if _, err := os.Stat("/Library"); err == nil {
		return "darwin"
	}
	return "linux"
}

// ConfigDir returns the configuration directory path.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "crx"), nil
}

// ConfigPath returns the configuration file path.
func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// Load loads the configuration from the default path.
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	return LoadFrom(path)
}

// LoadFrom loads the configuration from the specified path.
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to the default path.
func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	return c.SaveTo(path)
}

// SaveTo saves the configuration to the specified path.
func (c *Config) SaveTo(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// AddExtension adds an extension to the configuration.
func (c *Config) AddExtension(name string) bool {
	if slices.Contains(c.Extensions, name) {
		return false
	}
	c.Extensions = append(c.Extensions, name)
	return true
}

// RemoveExtension removes an extension from the configuration.
func (c *Config) RemoveExtension(name string) bool {
	for i, ext := range c.Extensions {
		if ext == name {
			c.Extensions = append(c.Extensions[:i], c.Extensions[i+1:]...)
			return true
		}
	}
	return false
}
