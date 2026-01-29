package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/sivchari/crx/internal/config"
	"github.com/sivchari/crx/internal/logger"
	"github.com/sivchari/crx/internal/policy"
	"github.com/sivchari/crx/internal/registry"
)

var dryRun bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the configuration and generate policy JSON",
	Long:  `Generates Chrome Enterprise Policy JSON from your configuration.`,
	Run:   runApply,
}

func init() {
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show changes without applying")
}

func runApply(cmd *cobra.Command, args []string) {
	logger.Debug("applying configuration", "dry_run", dryRun)

	cfg, err := config.Load()
	if err != nil {
		exitWithError("Failed to load configuration", err)
	}

	if len(cfg.Extensions) == 0 {
		fmt.Println("No extensions configured.")
		return
	}
	logger.Debug("extensions loaded", "count", len(cfg.Extensions))

	// Load packages from registry
	packages, err := loadPackages(cfg.Extensions)
	if err != nil {
		exitWithError("Failed to load packages", err)
	}
	logger.Debug("packages fetched from registry", "count", len(packages))

	// Generate policy
	gen := policy.NewGenerator(cfg, packages)
	pol, err := gen.Generate()
	if err != nil {
		exitWithError("Failed to generate policy", err)
	}
	logger.Debug("policy generated", "mode", cfg.Settings.Mode)

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
	logger.Debug("policy applied", "path", policy.GetPolicyPath())

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

// loadPackages loads packages from the registry.
func loadPackages(extensions []string) ([]*registry.Package, error) {
	packages := make([]*registry.Package, 0, len(extensions))
	fetcher := registry.NewDefaultFetcher()

	for _, name := range extensions {
		pkg, err := fetcher.FetchPackage(name)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch package %s: %w", name, err)
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}
