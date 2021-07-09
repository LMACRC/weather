package cmd

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/lmacrc/weather/pkg/weather"
	"github.com/lmacrc/weather/pkg/weather/ftp"
	"github.com/lmacrc/weather/pkg/weather/realtime"
	"github.com/lmacrc/weather/pkg/weather/reporting"
	"github.com/lmacrc/weather/pkg/weather/service"
	"github.com/lmacrc/weather/pkg/weather/store"
	"github.com/pelletier/go-toml/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
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

			var cfg weather.Config

			if data, err := os.ReadFile(flags.Config); err != nil {
				log.Error("Unable to read config path.", zap.String("path", flags.Config), zap.Error(err))
				return err
			} else {
				dec := toml.NewDecoder(bytes.NewBuffer(data))
				if err := dec.Decode(&cfg); err != nil {
					log.Error("Unable to decode config file.", zap.String("path", flags.Config), zap.Error(err))
					return err
				}
			}

			s, err := store.New()
			if err != nil {
				return err
			}

			ftpSvc := ftp.New(cfg.Ftp)
			reportSvc := reporting.New(log, s, cfg.Latitude, cfg.Longitude)
			realtimeSvc, err := realtime.New(log, cfg.Realtime, reportSvc, ftpSvc)
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

	cmd.Flags().IntVarP(&flags.Port, "port", "p", 9876, "Port to listen for requests")
	cmd.Flags().StringVar(&flags.Config, "config", "weather.toml", "Config file for weather service")
	_ = cmd.MarkFlagRequired("config")

	return cmd
}
