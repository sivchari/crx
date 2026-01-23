package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/crx/internal/config"
	"github.com/user/crx/internal/registry"
	"github.com/user/crx/internal/tui"
)

var browseRegistryPath string

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse and select extensions interactively",
	Long:  `Opens an interactive TUI to browse and select extensions from the registry.`,
	Run:   runBrowse,
}

func init() {
	browseCmd.Flags().StringVar(&browseRegistryPath, "registry", "", "Local registry path (for testing)")
}

func runBrowse(cmd *cobra.Command, args []string) {
	packages, err := fetchAllPackages(browseRegistryPath)
	if err != nil {
		exitWithError("Failed to fetch packages", err)
	}

	if len(packages) == 0 {
		fmt.Println("No packages found in registry.")
		return
	}

	selected, err := tui.Run(packages)
	if err != nil {
		exitWithError("TUI error", err)
	}

	if len(selected) == 0 {
		fmt.Println("No extensions selected.")
		return
	}

	// Add selected extensions to config
	cfg, err := config.Load()
	if err != nil {
		exitWithError("Failed to load configuration", err)
	}

	added := 0
	for _, name := range selected {
		if cfg.AddExtension(name) {
			added++
		}
	}

	if added > 0 {
		if err := cfg.Save(); err != nil {
			exitWithError("Failed to save configuration", err)
		}
		fmt.Printf("Added %d extension(s) to configuration.\n", added)
		fmt.Println("Run 'crx apply' to generate policy.")
	} else {
		fmt.Println("All selected extensions were already in configuration.")
	}
}

// fetchAllPackages fetches all packages from local or remote registry.
func fetchAllPackages(localPath string) ([]*registry.Package, error) {
	if localPath != "" {
		loader := registry.NewLoader(localPath)
		return loader.LoadAllPackages()
	}

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	if len(cfg.Registries) == 0 {
		return nil, fmt.Errorf("no registries configured")
	}

	reg := cfg.Registries[0]
	if reg.Type != "github" {
		return nil, fmt.Errorf("unsupported registry type: %s", reg.Type)
	}

	fetcher := registry.NewGitHubFetcher(reg.Repo, reg.Ref)
	return fetcher.FetchAllPackages()
}
