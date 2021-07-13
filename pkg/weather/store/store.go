package store

import (
	"fmt"
	"time"

	"github.com/lmacrc/weather/pkg/weather/model"
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
		return nil, fmt.Errorf("db open: %w", err)
	}

	err = db.AutoMigrate(Observation{})
	if err != nil {
		return nil, fmt.Errorf("db migrate: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) DB() *gorm.DB { return s.db }

func (s *Store) WriteObservation(o model.Observation) (*model.Observation, error) {
	var mo Observation
	mo.FromObservation(o)
	tx := s.db.Create(&mo)
	return mo.ToObservation(), tx.Error
}

func (s *Store) LastObservation(now time.Time) *model.Observation {
	now = now.UTC()

	var res Observation
	tx := s.db.Where("timestamp <= ?", now).Order("timestamp DESC").Limit(1).Find(&res)
	if tx.RowsAffected == 0 {
		return nil
	}
	return res.ToObservation()
}
