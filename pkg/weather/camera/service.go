package camera

import (
	"context"
	"io"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Ftp interface {
	Upload(ctx context.Context, dir, filename string, r io.Reader) error
}

type Service struct {
	log        *zap.Logger
	ftp        Ftp
	uploadPath string
	filename   string
	schedule   cron.Schedule
}

func New(log *zap.Logger) (*Service, error) {
	return &Service{
		log: log.With(zap.String("service", "camera")),
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
			s.log.Info("Generating webcam image.")

			// TODO(sgc): Generate image
			var r io.Reader

			s.log.Info("Starting image upload.")
			err := s.ftp.Upload(ctx, s.uploadPath, ".txt", r)
			if err != nil {
				s.log.Error("Finish image upload. Failure.", zap.Error(err))
			} else {
				s.log.Info("Finish image upload. Success.")
			}
		}
	}
}
