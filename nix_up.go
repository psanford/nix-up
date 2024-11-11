package main

import (
	"fmt"
	"os"

	"github.com/psanford/nix-up/run"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd.AddCommand(run.Command())

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "nix-up",
	Short: "Nix pull based updater",
}
