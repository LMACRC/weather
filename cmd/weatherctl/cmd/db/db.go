package db

import (
	"github.com/lmacrc/weather/pkg/weather"
	"github.com/lmacrc/weather/pkg/weather/store"
	"github.com/spf13/cobra"
)

var dbFlags = struct {
	Config string
	Path   string
}{
	Path: "weather.db",
}

var (
	config weather.Config
	db     *store.Store
)

func NewDbCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Commands to query and manage the weather database",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			config, err = weather.ReadConfig(dbFlags.Config)
			if err != nil {
				return
			}

			db, err = store.New(store.WithPath(config.DbPath))
			return
		},
	}

	cmd.PersistentFlags().StringVar(&dbFlags.Config, "config", "", "Path to config file")
	_ = cmd.MarkPersistentFlagRequired("config")
	cmd.AddCommand(newGetLastCommand())
	cmd.AddCommand(newGetStatsCommand())

	return cmd
}
