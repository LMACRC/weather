// Package realtime is responsible for generating weather statistics
package realtime

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/lmacrc/weather/pkg/weather/reporting"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// https://cumuluswiki.org/a/Realtime.txt#List_of_fields_in_the_file

type Service struct {
	log        *zap.Logger
	reporter   Reporter
	ftp        Ftp
	schedule   cron.Schedule
	uploadPath string
}

type Reporter interface {
	Generate(ts time.Time) *reporting.Statistics
}

type Ftp interface {
	Upload(ctx context.Context, dir, filename string, r io.Reader) error
}

func New(log *zap.Logger, cfg Config, reporter Reporter, ftp Ftp) (*Service, error) {
	schedule, err := cron.ParseStandard(cfg.Cron)
	if err != nil {
		return nil, fmt.Errorf("error parsing cron: %w", err)
	}

	return &Service{
		log:        log.With(zap.String("service", "realtime")),
		reporter:   reporter,
		ftp:        ftp,
		schedule:   schedule,
		uploadPath: cfg.UploadPath,
	}, nil
}

func (s Service) Run(ctx context.Context) {
	for {
		ts := time.Now()
		next := s.schedule.Next(ts)
		sleep := next.Sub(ts)
		s.log.Info("Next upload scheduled", zap.Time("time", next), zap.Duration("wait_time", sleep))

		select {
		case <-ctx.Done():
			s.log.Info("Shutting down")
			return

		case <-time.After(sleep):
			s.log.Info("Generating realtime.txt data.")

			stats := s.reporter.Generate(ts)
			data, err := Statistics(*stats).MarshalText()
			if err != nil {
				s.log.Error("Unable to generate statistics", zap.Error(err))
				continue
			}

			s.log.Info("Uploading realtime.txt.")
			err = s.ftp.Upload(ctx, s.uploadPath, "realtime.txt", bytes.NewReader(data))
			if err != nil {
				s.log.Error("Failed to upload realtime.txt.", zap.Error(err))
			} else {
				s.log.Info("Completed uploading realtime.txt.")
			}
		}
	}
}
