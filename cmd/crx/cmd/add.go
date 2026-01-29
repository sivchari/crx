package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sivchari/crx/internal/config"
	"github.com/sivchari/crx/internal/registry"
)

var addCmd = &cobra.Command{
	Use:   "add <extension>",
	Short: "Add an extension to the configuration",
	Long:  `Adds the specified extension to your configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run:   runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	name := args[0]

	cfg, err := config.Load()
	if err != nil {
		exitWithError("Failed to load configuration", err)
	}

	// Verify extension exists in registry
	pkg, err := verifyPackage(name)
	if err != nil {
		exitWithError("Extension not found in registry", err)
	}

	if cfg.AddExtension(name) {
		if err := cfg.Save(); err != nil {
			exitWithError("Failed to save configuration", err)
		}
		fmt.Printf("Added extension: %s (%s)\n", pkg.DisplayName, name)
	} else {
		fmt.Printf("Extension already exists: %s\n", name)
	}
}

// verifyPackage checks if a package exists in the registry.
func verifyPackage(name string) (*registry.Package, error) {
	fetcher := registry.NewDefaultFetcher()
	return fetcher.FetchPackage(name)
}
