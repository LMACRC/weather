package store

import (
	"fmt"
	"time"

	"github.com/lmacrc/weather/pkg/weather"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type config struct {
	Path string
}

type OptionFn func(c *config)

func WithPath(path string) OptionFn {
	return func(c *config) {
		c.Path = path
	}
}

type Store struct {
	db *gorm.DB
}

func New(opts ...OptionFn) (*Store, error) {
	c := config{
		Path: "weather.db",
	}

	for _, opt := range opts {
		opt(&c)
	}

	db, err := gorm.Open(sqlite.Open(c.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db open error: %w", err)
	}

	err = db.AutoMigrate(weather.Observation{})
	if err != nil {
		return nil, fmt.Errorf("db migrate error: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Write(value interface{}) (int64, error) {
	tx := s.db.Create(value)
	return tx.RowsAffected, tx.Error
}

func (s *Store) LastObservation() *weather.Observation {
	var res weather.Observation
	tx := s.db.Order("timestamp DESC").First(&res)
	if tx.RowsAffected == 0 {
		return nil
	}
	return &res
}

func (s *Store) Statistics(now time.Time) *weather.Statistics {
	var stats weather.Statistics

	return &stats
}
