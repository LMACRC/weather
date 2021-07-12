package db

import (
	"github.com/lmacrc/weather/pkg/weather"
	_ "github.com/lmacrc/weather/pkg/weather/camera/remote"
	_ "github.com/lmacrc/weather/pkg/weather/camera/rpi"
	"github.com/lmacrc/weather/pkg/weather/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dbFlags = struct {
	Config string
	Path   string
}{
	Path: "weather.db",
}

var (
	db *store.Store
)

func NewDbCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Commands to query and manage the weather database",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			err = weather.ReadConfig(dbFlags.Config)
			if err != nil {
				return
			}

			db, err = store.New(store.WithPath(viper.GetString("database_path")))
			return
		},
	}

	cmd.PersistentFlags().StringVar(&dbFlags.Config, "config", "", "Override  config file for weather service")
	cmd.AddCommand(newGetLastCommand())
	cmd.AddCommand(newGetStatsCommand())
	cmd.AddCommand(newGetImageCommand())

	return cmd
}
