// Package realtime is responsible for generating weather statistics
package realtime

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lmacrc/weather/pkg/event"
	"github.com/lmacrc/weather/pkg/weather/reporting"
	"github.com/lmacrc/weather/pkg/weather/service"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// https://cumuluswiki.org/a/Realtime.txt#List_of_fields_in_the_file

var (
	// NewStatistics is a topic for publishing new statistics calculations.
	NewStatistics = event.T("realtime:new_stats")
)

type Service struct {
	log        *zap.Logger
	reporter   Reporter
	ftp        service.Ftp
	schedule   cron.Schedule
	bus        *event.Bus
	remotePath string
}

type Reporter interface {
	Generate(ts time.Time) *reporting.Statistics
}

func New(log *zap.Logger, v *viper.Viper, reporter Reporter, ftp service.Ftp, bus *event.Bus) (*Service, error) {
	var cfg Config
	if err := v.UnmarshalKey("realtime", &cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	schedule, err := cron.ParseStandard(cfg.Cron)
	if err != nil {
		return nil, fmt.Errorf("parsing cron: %w", err)
	}

	return &Service{
		log:        log.With(zap.String("service", "realtime")),
		reporter:   reporter,
		ftp:        ftp,
		schedule:   schedule,
		bus:        bus,
		remotePath: cfg.RemoteDir,
	}, nil
}

func (s Service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	for {
		ts := time.Now()
		next := s.schedule.Next(ts)
		sleep := next.Sub(ts)
		s.log.Info("Next upload scheduled.", zap.Time("time", next), zap.Duration("wait_time", sleep))

		select {
		case <-ctx.Done():
			s.log.Info("Shutting down.")
			return

		case <-time.After(sleep):
			s.log.Info("Generating realtime.txt data.")

			stats := s.reporter.Generate(next)
			s.bus.Publish(NewStatistics, stats)

			data, err := Statistics(*stats).MarshalText()
			if err != nil {
				s.log.Error("Unable to marshal realtime.txt statistics data.", zap.Error(err))
				continue
			}

			err = os.WriteFile("realtime.txt", data, 0777)
			if err != nil {
				s.log.Error("Unable to write realtime.txt file.", zap.Error(err))
				continue
			}

			expiresAt := s.schedule.Next(time.Now())
			s.log.Info("Enqueue realtime.txt for upload.", zap.Time("expires_at", expiresAt))
			err = s.ftp.Enqueue(service.FtpRequest{
				LocalPath:      "realtime.txt",
				RemoteDir:      s.remotePath,
				RemoteFilename: "realtime.txt",
				ExpiresAt:      &expiresAt,
			})
			if err != nil {
				s.log.Error("Failed to enqueue realtime.txt for upload.", zap.Error(err))
			}
		}
	}
}
