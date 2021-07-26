package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "weather",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErr(err)
		os.Exit(1)
	}
}
