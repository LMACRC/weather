package ftp

import (
	"testing"
	"time"

	sqlite2 "github.com/lmacrc/weather/pkg/sql/driver/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func mustOpenDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to open database")
	}
	err = db.AutoMigrate(&QueueEntry{})
	if err != nil {
		panic("failed to migrate database")
	}
	return db
}

func TestService_nextDue(t *testing.T) {
	t.Run("with data", func(t *testing.T) {
		rows := []QueueEntry{
			{
				Due: sqlite2.FromTime(time.Date(2004, 4, 9, 12, 13, 14, 0, time.Local)),
			},
			{
				Due: sqlite2.FromTime(time.Date(2006, 9, 1, 12, 13, 14, 0, time.Local)),
			},
		}

		db := mustOpenDb()
		tx := db.Create(&rows)
		require.EqualValues(t, len(rows), tx.RowsAffected)
		require.NoError(t, tx.Error)

		s := Service{
			log: zaptest.NewLogger(t),
			db:  db,
		}
		d, err := s.nextDue()
		assert.NoError(t, err)
		assert.NotNil(t, d)
		assert.True(t, d.Equal(time.Date(2004, 4, 9, 12, 13, 14, 0, time.Local)))
	})

	t.Run("without data", func(t *testing.T) {
		db := mustOpenDb()

		s := Service{
			log: zaptest.NewLogger(t),
			db:  db,
		}
		d, err := s.nextDue()
		assert.NoError(t, err)
		assert.Nil(t, d)
	})
}

func TestService_findEntries(t *testing.T) {
	t.Run("with data", func(t *testing.T) {
		rows := []QueueEntry{
			{
				Due: sqlite2.FromTime(time.Date(2004, 4, 9, 12, 13, 14, 0, time.Local)),
			},
			{
				Due: sqlite2.FromTime(time.Date(2006, 9, 1, 12, 13, 14, 0, time.Local)),
			},
		}

		db := mustOpenDb()
		tx := db.Create(&rows)
		require.EqualValues(t, len(rows), tx.RowsAffected)
		require.NoError(t, tx.Error)

		s := Service{
			log: zaptest.NewLogger(t),
			db:  db,
		}

		t.Run("none due yet", func(t *testing.T) {
			entries, err := s.findEntries(time.Date(2000, 1, 1, 0, 0, 0, 1, time.Local))
			assert.NoError(t, err)
			assert.Empty(t, entries)
		})

		t.Run("one due", func(t *testing.T) {
			entries, err := s.findEntries(time.Date(2005, 1, 1, 10, 0, 0, 0, time.Local))
			assert.NoError(t, err)
			assert.Len(t, entries, 1)
		})

		t.Run("one due exact time", func(t *testing.T) {
			entries, err := s.findEntries(time.Date(2004, 4, 9, 12, 13, 14, 0, time.Local))
			assert.NoError(t, err)
			assert.Len(t, entries, 1)
		})

		t.Run("all due", func(t *testing.T) {
			entries, err := s.findEntries(time.Now())
			assert.NoError(t, err)
			assert.Len(t, entries, 2)
		})
	})

	t.Run("without data", func(t *testing.T) {
		db := mustOpenDb()

		s := Service{
			log: zaptest.NewLogger(t),
			db:  db,
		}
		entries, err := s.findEntries(time.Now())
		assert.NoError(t, err)
		assert.Empty(t, entries)
	})
}
