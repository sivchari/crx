package registry

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// GitHubFetcher fetches registry data from GitHub.
type GitHubFetcher struct {
	repo   string // e.g., "user/crx-registry"
	ref    string // e.g., "main"
	client *http.Client
	cache  *Cache
}

// NewGitHubFetcher creates a new GitHubFetcher.
func NewGitHubFetcher(repo, ref string) *GitHubFetcher {
	if ref == "" {
		ref = "main"
	}
	return &GitHubFetcher{
		repo: repo,
		ref:  ref,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache: NewCache(),
	}
}

// FetchRegistry fetches the registry index.
func (f *GitHubFetcher) FetchRegistry() (*Registry, error) {
	url := f.rawURL("registry.yaml")

	data, err := f.fetch(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry.yaml: %w", err)
	}

	var reg Registry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("failed to parse registry.yaml: %w", err)
	}

	return &reg, nil
}

// FetchPackage fetches a package definition.
func (f *GitHubFetcher) FetchPackage(name string) (*Package, error) {
	url := f.rawURL(fmt.Sprintf("pkgs/%s.yaml", name))

	data, err := f.fetch(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch package %s: %w", name, err)
	}

	var pkg Package
	if err := yaml.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse package %s: %w", name, err)
	}

	return &pkg, nil
}

// FetchAllPackages fetches all packages listed in the registry.
func (f *GitHubFetcher) FetchAllPackages() ([]*Package, error) {
	reg, err := f.FetchRegistry()
	if err != nil {
		return nil, err
	}

	packages := make([]*Package, 0, len(reg.Packages))
	for _, name := range reg.Packages {
		pkg, err := f.FetchPackage(name)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

// Search searches packages by name or tag.
func (f *GitHubFetcher) Search(query string) ([]*Package, error) {
	packages, err := f.FetchAllPackages()
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

func (f *GitHubFetcher) rawURL(path string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", f.repo, f.ref, path)
}

func (f *GitHubFetcher) fetch(url string) ([]byte, error) {
	// Check cache first
	if data, ok := f.cache.Get(url); ok {
		return data, nil
	}

	resp, err := f.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Store in cache
	f.cache.Set(url, data)

	return data, nil
}

// Cache provides simple in-memory caching.
type Cache struct {
	data map[string]cacheEntry
}

type cacheEntry struct {
	data      []byte
	timestamp time.Time
}

// NewCache creates a new Cache.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]cacheEntry),
	}
}

// Get retrieves data from cache.
func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.data[key]
	if !ok {
		return nil, false
	}

	// Cache expires after 5 minutes
	if time.Since(entry.timestamp) > 5*time.Minute {
		delete(c.data, key)
		return nil, false
	}

	return entry.data, true
}

// Set stores data in cache.
func (c *Cache) Set(key string, data []byte) {
	c.data[key] = cacheEntry{
		data:      data,
		timestamp: time.Now(),
	}
}

// CacheDir returns the cache directory path.
func CacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache", "crx"), nil
}
