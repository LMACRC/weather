package sqlite

import (
	"database/sql/driver"
	"errors"
	"time"
)

const CurrentTimeStamp = "2006-01-02 15:04:05"

type Timestamp struct{ time.Time }

func FromTime(t time.Time) Timestamp {
	return Timestamp{t}
}

func (dt Timestamp) Value() (driver.Value, error) {
	return dt.Time.UTC().Format(CurrentTimeStamp), nil
}

func (dt *Timestamp) Scan(value interface{}) error {
	var err error
	switch v := value.(type) {
	case string:
		dt.Time, err = time.Parse(CurrentTimeStamp, v)
	case []byte:
		dt.Time, err = time.Parse(CurrentTimeStamp, string(v))
	case time.Time:
		dt.Time = v
	default:
		err = errors.New("invalid type for current_timestamp")
	}
	return err
}

func (dt *Timestamp) MarshalCSV() (string, error) {
	return dt.Time.UTC().Format(CurrentTimeStamp), nil
}

func (dt *Timestamp) UnmarshalCSV(csv string) (err error) {
	dt.Time, err = time.Parse(CurrentTimeStamp, csv)
	return err
}
