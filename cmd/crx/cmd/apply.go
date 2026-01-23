package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/user/crx/internal/config"
	"github.com/user/crx/internal/policy"
	"github.com/user/crx/internal/registry"
)

var (
	dryRun       bool
	registryPath string
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the configuration and generate policy JSON",
	Long:  `Generates Chrome Enterprise Policy JSON from your configuration.`,
	Run:   runApply,
}

func init() {
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show changes without applying")
	applyCmd.Flags().StringVar(&registryPath, "registry", "", "Local registry path (for testing)")
}

func runApply(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		exitWithError("Failed to load configuration", err)
	}

	if len(cfg.Extensions) == 0 {
		fmt.Println("No extensions configured.")
		return
	}

	// Load packages from registry
	packages, err := loadPackages(cfg, registryPath)
	if err != nil {
		exitWithError("Failed to load packages", err)
	}

	// Generate policy
	gen := policy.NewGenerator(cfg, packages)
	pol, err := gen.Generate()
	if err != nil {
		exitWithError("Failed to generate policy", err)
	}

	if dryRun {
		json, err := gen.ToJSON(pol)
		if err != nil {
			exitWithError("Failed to generate JSON", err)
		}
		fmt.Println("Generated policy (dry-run):")
		fmt.Println(json)
		return
	}

	// Apply policy using OS-specific method
	if err := gen.Apply(pol); err != nil {
		exitWithError("Failed to apply policy", err)
	}

	fmt.Printf("Policy applied to: %s\n", policy.GetPolicyPath())

	switch runtime.GOOS {
	case "darwin":
		fmt.Println("\nmacOS: Profile opened in System Settings.")
		fmt.Println("Please click 'Install' to apply the policy.")
		fmt.Println("After installation, reload policies at chrome://policy or restart Chrome.")
		fmt.Println("\nTo remove: System Settings > Privacy & Security > Profiles > crx Chrome Extensions > Remove")
	default:
		fmt.Println("Reload policies at chrome://policy or restart Chrome to apply changes.")
	}
}

// loadPackages loads packages from local or remote registry.
func loadPackages(cfg *config.Config, localPath string) ([]*registry.Package, error) {
	packages := make([]*registry.Package, 0, len(cfg.Extensions))

	if localPath != "" {
		// Use local registry
		loader := registry.NewLoader(localPath)
		for _, name := range cfg.Extensions {
			pkg, err := loader.LoadPackage(name)
			if err != nil {
				return nil, fmt.Errorf("failed to load package %s: %w", name, err)
			}
			packages = append(packages, pkg)
		}
		return packages, nil
	}

	// Use remote registry from config
	if len(cfg.Registries) == 0 {
		return nil, fmt.Errorf("no registries configured")
	}

	reg := cfg.Registries[0]
	if reg.Type != "github" {
		return nil, fmt.Errorf("unsupported registry type: %s", reg.Type)
	}

	fetcher := registry.NewGitHubFetcher(reg.Repo, reg.Ref)
	for _, name := range cfg.Extensions {
		pkg, err := fetcher.FetchPackage(name)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch package %s: %w", name, err)
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}
