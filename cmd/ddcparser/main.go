package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:   "ddcparser",
		Short: "dcc parser",
	}
	rootCmd.AddCommand(StartCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed executing ddcparser command: %s, exiting...\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
