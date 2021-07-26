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

	"github.com/lmacrc/weather/pkg/event"
	"github.com/lmacrc/weather/pkg/weather"
	whttp "github.com/lmacrc/weather/pkg/weather/http"
	"github.com/lmacrc/weather/pkg/weather/reporting"
	"github.com/lmacrc/weather/pkg/weather/service/archive"
	"github.com/lmacrc/weather/pkg/weather/service/camera"
	ftp2 "github.com/lmacrc/weather/pkg/weather/service/ftp"
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
				return err
			}

			log.Info("Loaded config.", zap.String("path", viper.ConfigFileUsed()))

			db, err := weather.OpenDb(viper.GetString("database_path"))
			if err != nil {
				return err
			}

			s, err := store.New(db, store.WithBus(bus))
			if err != nil {
				return err
			}

			vp := viper.GetViper()
			whttp.InitViper(vp)
			archive.InitViper(vp)

			var ftpSvc *ftp2.Service
			if viper.GetBool("ftp.enabled") {
				ftpSvc, err = ftp2.New(log, db, vp)
				if err != nil {
					return err
				}

				go func() {
					ftpSvc.Run(ctx)
				}()
			} else {
				log.Info("FTP service disabled.")
			}

			cs := cron.New()

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

			if viper.GetBool("archive.enabled") {
				archiveSvc, err := archive.New(log, db, vp, s, ftpSvc)
				if err != nil {
					return err
				}

				// run at 00:30 each day
				sch, err := cron.ParseStandard("30 0 * * *")
				if err != nil {
					panic(fmt.Sprintf("Unable to parse cron spec: %s", err))
				}

				cs.Schedule(sch, archiveSvc)
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
