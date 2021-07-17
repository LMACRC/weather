package cmd

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/lmacrc/weather/pkg/event"
	"github.com/lmacrc/weather/pkg/ftp"
	"github.com/lmacrc/weather/pkg/weather"
	"github.com/lmacrc/weather/pkg/weather/camera"
	_ "github.com/lmacrc/weather/pkg/weather/camera/remote"
	_ "github.com/lmacrc/weather/pkg/weather/camera/rpi"
	whttp "github.com/lmacrc/weather/pkg/weather/http"
	"github.com/lmacrc/weather/pkg/weather/realtime"
	"github.com/lmacrc/weather/pkg/weather/reporting"
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
			bus := event.New()

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			logCfg := zap.NewDevelopmentConfig()
			logCfg.DisableCaller = true
			logCfg.Development = false
			log, _ := logCfg.Build()

			err := weather.ReadConfig(flags.Config)
			if err != nil {
				return err
			}

			log.Info("Loaded config.", zap.String("path", viper.ConfigFileUsed()))

			s, err := store.New(store.WithPath(viper.GetString("database_path")), store.WithBus(bus))
			if err != nil {
				return err
			}

			vp := viper.GetViper()
			whttp.InitViper(vp)

			ftpSvc, err := ftp.New(vp)
			if err != nil {
				return err
			}

			if viper.GetBool("realtime.enabled") {
				reportSvc, err := reporting.New(log, vp, s)
				if err != nil {
					return err
				}

				realtimeSvc, err := realtime.New(log, vp, reportSvc, ftpSvc)
				if err != nil {
					return err
				}

				go func() {
					realtimeSvc.Run(ctx)
				}()
			} else {
				log.Info("Realtime service disabled.")
			}

			if viper.GetBool("camera.enabled") {
				cameraSvc, err := camera.New(log, vp, ftpSvc)
				if err != nil {
					return err
				}

				go func() {
					cameraSvc.Run(ctx)
				}()
			} else {
				log.Info("Camera service disabled.")
			}

			mux := http.NewServeMux()

			mux.Handle("/metrics", promhttp.Handler())

			wh, err := whttp.New(log, vp, s)
			if err != nil {
				return err
			}

			wh.Handle(mux)

			var lc net.ListenConfig
			addr := ":" + strconv.Itoa(flags.Port)
			ln, err := lc.Listen(ctx, "tcp", addr)
			if err != nil {
				return err
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()

				_ = http.Serve(ln, mux)
				log.Info("HTTP server shut down.")
			}()

			<-ctx.Done()
			log.Info("Shutting down.")

			_ = ln.Close()

			wg.Wait()

			log.Info("Shutdown complete.")

			return nil
		},
	}

	fs := cmd.Flags()
	fs.IntVarP(&flags.Port, "port", "p", 9876, "Port to listen for requests")
	fs.StringVar(&flags.Config, "config", "", "Override config file for weather service")

	_ = viper.BindPFlag("server.port", fs.Lookup("port"))

	return cmd
}
