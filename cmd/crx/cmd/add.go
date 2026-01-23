package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/crx/internal/config"
	"github.com/user/crx/internal/registry"
)

var addRegistryPath string

var addCmd = &cobra.Command{
	Use:   "add <extension>",
	Short: "Add an extension to the configuration",
	Long:  `Adds the specified extension to your configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run:   runAdd,
}

func init() {
	addCmd.Flags().StringVar(&addRegistryPath, "registry", "", "Local registry path (for testing)")
}

func runAdd(cmd *cobra.Command, args []string) {
	name := args[0]

	cfg, err := config.Load()
	if err != nil {
		exitWithError("Failed to load configuration", err)
	}

	// Verify extension exists in registry
	pkg, err := verifyPackage(cfg, name, addRegistryPath)
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
func verifyPackage(cfg *config.Config, name, localPath string) (*registry.Package, error) {
	if localPath != "" {
		loader := registry.NewLoader(localPath)
		return loader.LoadPackage(name)
	}

	if len(cfg.Registries) == 0 {
		return nil, fmt.Errorf("no registries configured")
	}

	reg := cfg.Registries[0]
	if reg.Type != "github" {
		return nil, fmt.Errorf("unsupported registry type: %s", reg.Type)
	}

	fetcher := registry.NewGitHubFetcher(reg.Repo, reg.Ref)
	return fetcher.FetchPackage(name)
}
