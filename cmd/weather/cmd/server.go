package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/lmacrc/weather/pkg/cronzap"
	"github.com/lmacrc/weather/pkg/event"
	"github.com/lmacrc/weather/pkg/weather"
	whttp "github.com/lmacrc/weather/pkg/weather/http"
	"github.com/lmacrc/weather/pkg/weather/reporting"
	"github.com/lmacrc/weather/pkg/weather/service/archive"
	"github.com/lmacrc/weather/pkg/weather/service/camera"
	"github.com/lmacrc/weather/pkg/weather/service/ftp"
	"github.com/lmacrc/weather/pkg/weather/service/influxdb"
	"github.com/lmacrc/weather/pkg/weather/service/realtime"
	"github.com/lmacrc/weather/pkg/weather/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	// install camera drivers
	_ "github.com/lmacrc/weather/pkg/weather/service/camera/remote"
	_ "github.com/lmacrc/weather/pkg/weather/service/camera/rpi"
)

func init() {
	rootCmd.AddCommand(newServer())
}

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
				log.Error("Failed to read configuration.", zap.Error(err))
				return err
			}

			log.Info("Loaded config.", zap.String("path", viper.ConfigFileUsed()))

			db, err := weather.OpenDb(viper.GetString("database.url"))
			if err != nil {
				log.Error("Failed to open database.", zap.Error(err))
				return err
			}

			s, err := store.New(db, bus)
			if err != nil {
				return err
			}

			vp := viper.GetViper()
			whttp.InitViper(vp)
			archive.InitViper(vp)
			influxdb.InitViper(vp)

			var ftpSvc *ftp.Service
			if viper.GetBool("ftp.enabled") {
				ftpSvc, err = ftp.New(log, db, vp)
				if err != nil {
					log.Error("Failed to initialise FTP service.", zap.Error(err))
					return err
				}

				go func() {
					ftpSvc.Run(ctx)
				}()
			} else {
				log.Info("FTP service disabled.")
			}

			if viper.GetBool("influxdb.enabled") {
				influxDb, err := influxdb.New(log, vp, bus)
				if err != nil {
					log.Error("Failed to initialise influxdb service.", zap.Error(err))
					return err
				}

				go func() {
					influxDb.Run(ctx)
				}()
			} else {
				log.Info("InfluxDB service disabled.")
			}

			cs := cron.New(cron.WithLogger(&cronzap.Adapter{Log: log.With(zap.String("service", "cron"))}))

			if viper.GetBool("realtime.enabled") {
				reportSvc, err := reporting.New(log, vp, s)
				if err != nil {
					log.Error("Failed to initialise reporting service.", zap.Error(err))
					return err
				}

				realtimeSvc, err := realtime.New(log, vp, reportSvc, ftpSvc, bus)
				if err != nil {
					log.Error("Failed to initialise realtime service.", zap.Error(err))
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
					log.Error("Failed to initialise camera service.", zap.Error(err))
					return err
				}

				go func() {
					cameraSvc.Run(ctx)
				}()
			} else {
				log.Info("Camera service disabled.")
			}

			if viper.GetBool("archive.enabled") {
				log.Info("Archive service enabled.")

				archiveSvc, err := archive.New(log, db, vp, s, ftpSvc)
				if err != nil {
					log.Error("Failed to initialise archive service.", zap.Error(err))
					return err
				}

				// run at 00:30 each day
				sch, err := cron.ParseStandard("30 0 * * *")
				if err != nil {
					// Should never happen and represents a programming error
					panic(fmt.Sprintf("Unable to parse cron spec: %s", err))
				}

				cs.Schedule(sch, archiveSvc)
			} else {
				log.Info("Archive service disabled.")
			}

			cs.Start()

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
