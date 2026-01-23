package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sivchari/crx/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize crx configuration",
	Long:  `Creates a new configuration file with default settings.`,
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	path, err := config.ConfigPath()
	if err != nil {
		exitWithError("Failed to get config path", err)
	}

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Configuration already exists at %s\n", path)
		return
	}

	cfg := config.DefaultConfig()
	if err := cfg.Save(); err != nil {
		exitWithError("Failed to save configuration", err)
	}

	fmt.Printf("Configuration created at %s\n", path)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit the configuration to add your extensions")
	fmt.Println("  2. Run 'crx add <extension>' to add extensions")
	fmt.Println("  3. Run 'crx apply' to generate the policy JSON")
}
