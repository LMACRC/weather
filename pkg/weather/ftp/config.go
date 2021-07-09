package ftp

import (
	"time"
)

type Duration struct{ time.Duration }

func (d *Duration) UnmarshalText(text []byte) error {
	res, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	(*d).Duration = res
	return nil
}

type Config struct {
	Interval   Duration
	Address    string
	Username   string
	Password   string
	UploadPath string `toml:"upload_path"`
}
