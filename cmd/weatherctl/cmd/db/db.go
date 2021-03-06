package db

import (
	"github.com/lmacrc/weather/pkg/event"
	"github.com/lmacrc/weather/pkg/weather"
	whttp "github.com/lmacrc/weather/pkg/weather/http"
	"github.com/lmacrc/weather/pkg/weather/service/archive"
	"github.com/lmacrc/weather/pkg/weather/service/camera"
	"github.com/lmacrc/weather/pkg/weather/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	// install camera drivers
	_ "github.com/lmacrc/weather/pkg/weather/service/camera/remote"
	_ "github.com/lmacrc/weather/pkg/weather/service/camera/rpi"
)

var dbFlags = struct {
	Config string
	Path   string
}{
	Path: "weather.db",
}

var (
	bus *event.Bus
	db  *gorm.DB
	st  *store.Store
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

			vp := viper.GetViper()
			whttp.InitViper(vp)
			archive.InitViper(vp)
			camera.InitViper(vp)

			db, err = weather.OpenDb(viper.GetString("database.url"))
			if err != nil {
				return err
			}

			bus = event.New()

			st, err = store.New(db, bus)
			return
		},
	}

	cmd.PersistentFlags().StringVar(&dbFlags.Config, "config", "", "Override config file for weather service")
	cmd.AddCommand(newGetLastCommand())
	cmd.AddCommand(newGetStatsCommand())
	cmd.AddCommand(newGetImageCommand())
	cmd.AddCommand(newArchiveCommand())
	cmd.AddCommand(newArchiveAllCommand())

	return cmd
}
