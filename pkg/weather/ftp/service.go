package ftp

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/jlaffaye/ftp"
	"go.uber.org/zap"
)

type Datasource interface {
	fmt.Stringer
	Data(ts time.Time) (path string, r io.ReadCloser, err error)
}

type Service struct {
	Log         *zap.Logger
	Config      Config
	DataSources []Datasource
}

func (s *Service) Run(ctx context.Context) {
	d := s.Config.Interval.Duration

	s.Log.Info("Starting FTP service.", zap.Duration("interval", d))

	for {
		nextWakeup := time.Now().Truncate(d).Add(d)
		delay := nextWakeup.Sub(time.Now())

		s.Log.Info("Waiting to upload data.", zap.Time("wakeup", nextWakeup), zap.Duration("wait_time", delay))

		select {
		case <-ctx.Done():
			s.Log.Info("Received shutdown signal.")
			return
		case <-time.After(delay):
			s.Log.Info("Uploading data.", zap.String("addr", s.Config.Address), zap.Time("time", nextWakeup))
			s.uploadData(ctx, nextWakeup)
		}
	}
}

func (s *Service) uploadData(ctx context.Context, ts time.Time) {
	conn, err := ftp.Dial(s.Config.Address, ftp.DialWithContext(ctx))
	if err != nil {
		s.Log.Error("Unable to connect to server.", zap.Error(err))
		return
	}
	defer func() { _ = conn.Quit() }()

	if err := conn.Login(s.Config.Username, s.Config.Password); err != nil {
		s.Log.Error("Unable to log in to server.", zap.String("user", s.Config.Username), zap.Error(err))
		return
	}

	if err := conn.ChangeDir(s.Config.UploadPath); err != nil {
		s.Log.Error("Unable to change dir.", zap.String("path", s.Config.UploadPath), zap.Error(err))
		return
	}

	for _, ds := range s.DataSources {
		log := s.Log.With(zap.Stringer("source", ds))

		path, r, err := ds.Data(ts)
		if err != nil {
			log.Error("Unable to fetch data for data source. Skipping.", zap.Error(err))
			continue
		}

		log = log.With(zap.String("path", path))

		func() {
			defer func() { _ = r.Close() }()
			if err := conn.Stor(path, r); err != nil {
				log.Error("Unable to store data for data source.", zap.Error(err))
				return
			}
			log.Info("Uploaded data to server.")
		}()
	}
}
