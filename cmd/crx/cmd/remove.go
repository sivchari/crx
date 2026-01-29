package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sivchari/crx/internal/config"
	"github.com/sivchari/crx/internal/logger"
)

var removeCmd = &cobra.Command{
	Use:   "remove <extension>",
	Short: "Remove an extension from the configuration",
	Long:  `Removes the specified extension from your configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run:   runRemove,
}

func runRemove(cmd *cobra.Command, args []string) {
	name := args[0]
	logger.Debug("removing extension", "name", name)

	cfg, err := config.Load()
	if err != nil {
		exitWithError("Failed to load configuration", err)
	}

	if cfg.RemoveExtension(name) {
		if err := cfg.Save(); err != nil {
			exitWithError("Failed to save configuration", err)
		}
		logger.Debug("extension removed from configuration")
		fmt.Printf("Removed extension: %s\n", name)
		fmt.Println("Run 'crx apply' to update the policy.")
	} else {
		logger.Debug("extension not found in configuration")
		fmt.Printf("Extension not found: %s\n", name)
	}
}
