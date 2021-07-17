package store

import (
	"fmt"
	"time"

	"github.com/lmacrc/weather/pkg/event"
	"github.com/lmacrc/weather/pkg/weather/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// NewObservation is a topic for publishing new observations.
	NewObservation = event.T("store:new_observation")
)

type config struct {
	Path string
	Bus  *event.Bus
}

type OptionFn func(c *config)

func WithPath(path string) OptionFn {
	return func(c *config) {
		c.Path = path
	}
}

func WithBus(eb *event.Bus) OptionFn {
	return func(c *config) {
		c.Bus = eb
	}
}

type Store struct {
	db  *gorm.DB
	bus *event.Bus
}

func New(opts ...OptionFn) (*Store, error) {
	c := config{
		Path: "weather.db",
	}

	for _, opt := range opts {
		opt(&c)
	}

	bus := c.Bus
	if bus == nil {
		bus = event.New()
	}

	db, err := gorm.Open(sqlite.Open(c.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	err = db.AutoMigrate(Observation{})
	if err != nil {
		return nil, fmt.Errorf("db migrate: %w", err)
	}

	return &Store{db: db, bus: bus}, nil
}

func (s *Store) DB() *gorm.DB { return s.db }

func (s *Store) WriteObservation(o model.Observation) (*model.Observation, error) {
	var mo Observation
	mo.FromObservation(o)
	tx := s.db.Create(&mo)
	res := mo.ToObservation()

	if tx.Error == nil {
		s.bus.Publish(NewObservation, res)
	}

	return res, tx.Error
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
