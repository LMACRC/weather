// Package realtime is responsible for generating weather statistics
package realtime

import (
	"gorm.io/gorm"
)

// https://cumuluswiki.org/a/Realtime.txt#List_of_fields_in_the_file

type Service struct {
	DB *gorm.DB
}
