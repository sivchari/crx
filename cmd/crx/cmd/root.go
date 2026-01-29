package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/sivchari/crx/internal/logger"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "crx",
	Short: "Chrome Extension Manager",
	Long: `crx is a declarative Chrome extension manager.

It allows you to manage Chrome extensions using a YAML configuration file
and generates Chrome Enterprise Policy JSON for installation.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Init(verbose)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(browseCmd)
}

func exitWithError(msg string, err error) {
	if err != nil {
		logger.Error(msg, "error", err)
	} else {
		logger.Error(msg)
	}
	os.Exit(1)
}
