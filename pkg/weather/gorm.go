package weather

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDb(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	return db, nil
}
