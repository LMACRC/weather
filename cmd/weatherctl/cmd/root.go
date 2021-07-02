package cmd

import (
	"os"

	"github.com/lmacrc/weather/cmd/weatherctl/cmd/db"
	"github.com/spf13/cobra"
)

func Execute() {
	cmd := &cobra.Command{
		Use:   "weatherctl",
		Short: "weather control interface",
	}

	cmd.AddCommand(db.NewDbCommand())

	if err := cmd.Execute(); err != nil {
		cmd.PrintErr(err)
		os.Exit(1)
	}
}
