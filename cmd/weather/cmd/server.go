package cmd

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/lmacrc/weather/pkg/weather"
	_ "github.com/lmacrc/weather/pkg/weather/camera/remote"
	_ "github.com/lmacrc/weather/pkg/weather/camera/rpi"
	"github.com/lmacrc/weather/pkg/weather/ftp"
	"github.com/lmacrc/weather/pkg/weather/realtime"
	"github.com/lmacrc/weather/pkg/weather/reporting"
	"github.com/lmacrc/weather/pkg/weather/service"
	"github.com/lmacrc/weather/pkg/weather/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func newServer() *cobra.Command {
	flags := struct {
		Port   int
		Config string
	}{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the LMACRC weather server",
		RunE: func(cmd *cobra.Command, args []string) error {
			logCfg := zap.NewDevelopmentConfig()
			logCfg.DisableCaller = true
			logCfg.Development = false
			log, _ := logCfg.Build()

			err := weather.ReadConfig(flags.Config)
			if err != nil {
				return err
			}

			log.Info("Loaded config.", zap.String("path", viper.ConfigFileUsed()))

			s, err := store.New(store.WithPath(viper.GetString("database_path")))
			if err != nil {
				return err
			}

			ftpSvc, _ := ftp.New(viper.GetViper())
			reportSvc, err := reporting.New(log, viper.GetViper(), s)
			if err != nil {
				return err
			}

			realtimeSvc, err := realtime.New(log, viper.GetViper(), reportSvc, ftpSvc)
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			go func() {
				realtimeSvc.Run(ctx)
			}()

			mux := http.NewServeMux()

			mux.Handle("/metrics", promhttp.Handler())

			wh := service.Handler{Path: "/weather", Store: s}
			wh.Handle(mux)

			v := viper.GetInt("server.port")
			_ = v

			var lc net.ListenConfig
			addr := ":" + strconv.Itoa(flags.Port)
			ln, err := lc.Listen(ctx, "tcp", addr)
			if err != nil {
				return err
			}

			go func() {
				_ = http.Serve(ln, mux)
			}()

			<-ctx.Done()
			log.Info("Shutting down.")

			_ = ln.Close()

			return nil
		},
	}

	fs := cmd.Flags()
	fs.IntVarP(&flags.Port, "port", "p", 9876, "Port to listen for requests")
	fs.StringVar(&flags.Config, "config", "", "Override config file for weather service")

	_ = viper.BindPFlag("server.port", fs.Lookup("port"))

	return cmd
}
