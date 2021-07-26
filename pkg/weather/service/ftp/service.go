package ftp

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lmacrc/weather/pkg/sql/driver/sqlite"
	"github.com/lmacrc/weather/pkg/weather/service"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type optionFn func(s *Service)

func WithClient(client Client) optionFn {
	if client == nil {
		panic("client == nil")
	}

	return func(s *Service) {
		s.client = client
	}
}

type Service struct {
	log *zap.Logger
	db  *gorm.DB
	ch  chan struct{}

	client   Client
	address  string
	username string
	password string
}

func New(log *zap.Logger, db *gorm.DB, v *viper.Viper, opts ...optionFn) (*Service, error) {
	var cfg Config
	if err := v.UnmarshalKey("ftp", &cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	err := db.AutoMigrate(QueueEntry{})
	if err != nil {
		return nil, fmt.Errorf("db migrate: %w", err)
	}

	s := &Service{
		log:      log.With(zap.String("service", "ftp")),
		db:       db,
		ch:       make(chan struct{}),
		client:   &ftpClient{},
		address:  cfg.Address,
		username: cfg.Username,
		password: cfg.Password,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

func (s *Service) Enqueue(req service.FtpRequest) (err error) {
	defer func() {
		if err == nil {
			select {
			case s.ch <- struct{}{}:
				// wakeup the message pump
			default:
			}
		}
	}()

	rec := NewFromFtpRequest(req)
	rec.Due = sqlite.Timestamp{Time: time.Now().UTC()}
	return s.db.Create(&rec).Error
}

func (s *Service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	for {
		select {
		case <-s.ch:
			s.RunQueue(ctx)
		case <-ctx.Done():
			s.log.Info("Shutting down.")
			return
		}
	}
}

func (s *Service) RunQueue(ctx context.Context) {
	entries, err := s.findEntries(time.Now())
	if err != nil {
		// TODO(sgc): Should we fatal here and let the process restart?
		s.log.Fatal("Unable to read queue entries. Terminating.", zap.Error(err))
		return
	}

	for i := range entries {
		e := entries[i]
		log := s.log.With(zap.String("local_path", e.LocalPath), zap.String("remote_dir", e.RemoteDir), zap.String("remote_filename", e.RemoteFilename))
		log.Info("Sending file via FTP.")

		if e.ExpiresAt != nil && e.ExpiresAt.Before(time.Now()) {
			log.Info("File upload has expired.", zap.Time("expired_at", (*e.ExpiresAt).Time))
			s.completeEntry(e)
			continue
		}

		s.processEntry(ctx, log, e)
	}
}

func (s *Service) processEntry(ctx context.Context, log *zap.Logger, e *QueueEntry) {
	f, err := os.Open(e.LocalPath)
	if err != nil {
		log.Error("Unable to read local file. Removing from queue.", zap.Error(err))
		s.completeEntry(e)
		return
	}
	defer func() { _ = f.Close() }()

	err = s.client.Upload(ctx, s.address, s.username, s.password, f, e.RemoteDir, e.RemoteFilename)
	if err != nil {
		// TODO(sgc): retry?
		return
	}

	log.Info("Upload complete.")

	s.completeEntry(e)
}

func (s *Service) completeEntry(e *QueueEntry) {
	if e.RemoveLocal {
		_ = os.Remove(e.LocalPath)
	}
	s.db.Delete(e)
}

func (s *Service) findEntries(ts time.Time) ([]*QueueEntry, error) {
	var rows []*QueueEntry
	return rows, s.db.Model(&QueueEntry{}).
		Where("due <= ?", ts.UTC()).
		Find(&rows).Error
}

func (s *Service) nextDue() (*time.Time, error) {
	var res struct {
		MinDue sqlite.Timestamp
	}
	tx := s.db.Model(&QueueEntry{}).
		Select("MIN(due) as min_due").
		Find(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 || res.MinDue.Time.IsZero() {
		return nil, nil
	}

	return &res.MinDue.Time, nil
}
