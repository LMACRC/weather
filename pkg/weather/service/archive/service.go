package archive

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/jinzhu/now"
	"github.com/lmacrc/weather/pkg/compress/brotli"
	"github.com/lmacrc/weather/pkg/filepath/template"
	"github.com/lmacrc/weather/pkg/weather/service"
	"github.com/lmacrc/weather/pkg/weather/store"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	log         *zap.Logger
	db          *gorm.DB
	store       *store.Store
	compression Compression
	ftp         service.Ftp

	localDir  string
	remoteDir string
	filename  *template.Template
}

func InitViper(v *viper.Viper) {
	v.SetDefault("archive.compression", CompressionGzip)
}

func New(log *zap.Logger, db *gorm.DB, v *viper.Viper, s *store.Store, ftp service.Ftp) (*Service, error) {
	var cfg Config
	if err := v.UnmarshalKey("archive", &cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc())); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	a := &Service{
		log:         log,
		db:          db,
		store:       s,
		compression: cfg.Compression,
		ftp:         ftp,
		localDir:    cfg.LocalDir,
		remoteDir:   cfg.RemoteDir,
		filename:    template.Must(template.New("file").Parse(cfg.Filename)),
	}

	return a, nil
}

func (s *Service) Run() {
	s.log.Info("Starting archive process.")
	err := s.ArchiveAll()
	if err != nil {
		s.log.Error("Failed to archive data.", zap.Error(err))
	} else {
		s.log.Info("Completed archive process.")
	}
}

// ArchiveAll archives all days prior to now
func (s *Service) ArchiveAll() error {
	last := now.BeginningOfDay()
	s.log.Info("Archiving all data prior to today.", zap.Time("date", last))

	dates, err := s.findAllDates(last)
	if err != nil {
		s.log.Error("Unable to locate dates to archive.", zap.Error(err))
		return err
	}
	if len(dates) == 0 {
		s.log.Info("No data to archive.")
		return nil
	}

	for _, dt := range dates {
		log := s.log.With(zap.String("date", dt.Format("20060102")))

		log.Info("Archiving data for date.")
		path, err := s.Archive(dt)
		if err != nil {
			// TODO(sgc): This could be
			log.Error("Failed to archive data.", zap.Error(err))
			continue
		}
		log.Info("Data archived to file.", zap.String("path", path))

		if s.ftp != nil {
			log.Info("Queueing file for FTP.")
			s.ftp.Enqueue(service.FtpRequest{
				LocalPath:      path,
				RemoteDir:      s.remoteDir,
				RemoteFilename: filepath.Base(path),
				RemoveLocal:    true,
			})
		}
	}

	s.log.Info("Executing SQLite VACUUM.")
	s.db.Exec("VACUUM")
	s.log.Info("SQLite VACUUM complete.")

	return nil
}

func (s *Service) findAllDates(t time.Time) ([]time.Time, error) {
	db := s.store.DB()
	var dates []string
	tx := db.Table("observations").
		Where("timestamp < ?", t.UTC()).
		Select("strftime('%Y%m%d', datetime(timestamp, 'localtime')) as date").
		Group("date").
		Find(&dates)

	var ts []time.Time
	for _, ds := range dates {
		dt, err := time.ParseInLocation("20060102", ds, time.Local)
		if err != nil {
			continue
		}
		ts = append(ts, dt)
	}

	return ts, tx.Error
}

var (
	ErrNoData = errors.New("no data")
)

// Archive will archive the data for the day specified by t and return a path to the archived
// file.
func (s *Service) Archive(t time.Time) (path string, err error) {
	tt := now.With(t)
	start := tt.BeginningOfDay()
	end := start.AddDate(0, 0, 1)

	var rows []*store.Observation
	rows, err = s.findRows(s.db, start, end)
	if err != nil {
		return "", err
	}

	if len(rows) == 0 {
		return "", ErrNoData
	}

	path = fmt.Sprintf("observations_%s.csv", start.Format("20060102"))
	var useBrotli = brotli.IsAvailable() && s.compression == CompressionBrotli

	if useBrotli {
		path = path + ".br"
	} else {
		path = path + ".gz"
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	var wr io.WriteCloser
	if useBrotli {
		// Ignore err as brotli is available, given earlier test
		wr, _ = brotli.NewWriter(f)
	} else {
		wr, _ = gzip.NewWriterLevel(f, gzip.BestCompression)
	}
	defer func() { _ = wr.Close() }()

	err = gocsv.Marshal(rows, wr)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (s *Service) findRows(tx *gorm.DB, start, end time.Time) ([]*store.Observation, error) {
	var rows []*store.Observation
	tx.Where("timestamp >= ? AND timestamp < ?", start.UTC(), end.UTC()).
		Order("timestamp").
		Find(&rows)

	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}

	var deleted []*store.Observation
	tx.Delete(&deleted, ids)

	return rows, tx.Error
}
