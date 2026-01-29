package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sivchari/crx/internal/config"
	"github.com/sivchari/crx/internal/logger"
	"github.com/sivchari/crx/internal/registry"
	"github.com/sivchari/crx/internal/tui"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse and select extensions interactively",
	Long:  `Opens an interactive TUI to browse and select extensions from the registry.`,
	Run:   runBrowse,
}

func runBrowse(cmd *cobra.Command, args []string) {
	logger.Debug("fetching packages from registry")
	packages, err := fetchAllPackages()
	if err != nil {
		exitWithError("Failed to fetch packages", err)
	}
	logger.Debug("packages fetched", "count", len(packages))

	if len(packages) == 0 {
		fmt.Println("No packages found in registry.")
		return
	}

	selected, err := tui.Run(packages)
	if err != nil {
		exitWithError("TUI error", err)
	}
	logger.Debug("extensions selected", "count", len(selected))

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
			logger.Debug("extension added", "name", name)
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

// fetchAllPackages fetches all packages from the registry.
func fetchAllPackages() ([]*registry.Package, error) {
	fetcher := registry.NewDefaultFetcher()
	return fetcher.FetchAllPackages()
}
