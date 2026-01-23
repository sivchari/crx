package main

import (
	"os"

	"github.com/user/crx/cmd/crx/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
