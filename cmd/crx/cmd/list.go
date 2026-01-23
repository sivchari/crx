package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/crx/internal/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured extensions",
	Long:  `Lists all extensions in your configuration file.`,
	Run:   runList,
}

func runList(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		exitWithError("Failed to load configuration", err)
	}

	if len(cfg.Extensions) == 0 {
		fmt.Println("No extensions configured.")
		fmt.Println("Use 'crx add <extension>' to add extensions.")
		return
	}

	fmt.Println("Configured extensions:")
	for _, ext := range cfg.Extensions {
		fmt.Printf("  - %s\n", ext)
	}
}
