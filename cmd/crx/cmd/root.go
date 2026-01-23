package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crx",
	Short: "Chrome Extension Manager",
	Long: `crx is a declarative Chrome extension manager.

It allows you to manage Chrome extensions using a YAML configuration file
and generates Chrome Enterprise Policy JSON for installation.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(browseCmd)
}

func exitWithError(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
	os.Exit(1)
}
