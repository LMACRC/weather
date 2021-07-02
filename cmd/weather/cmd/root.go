package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	cmd := &cobra.Command{
		Use: "weather",
	}

	cmd.AddCommand(newServer())

	if err := cmd.Execute(); err != nil {
		cmd.PrintErr(err)
		os.Exit(1)
	}
}
