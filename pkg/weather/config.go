package weather

import (
	"bytes"
	"fmt"
	"os"

	"github.com/lmacrc/weather/pkg/weather/ftp"
	"github.com/lmacrc/weather/pkg/weather/realtime"
	"github.com/lmacrc/weather/pkg/weather/reporting"
	"github.com/pelletier/go-toml/v2"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type Config struct {
	DbPath   string `toml:"database_path"`
	Location Location

	Ftp       ftp.Config
	Realtime  realtime.Config
	Reporting reporting.Config
}

func ReadConfig(path string) (Config, error) {
	var cfg Config
	if data, err := os.ReadFile(path); err != nil {
		return Config{}, fmt.Errorf("read error: %w", err)
	} else {
		dec := toml.NewDecoder(bytes.NewBuffer(data))
		if err := dec.Decode(&cfg); err != nil {
			return Config{}, fmt.Errorf("decode error: %w", err)
		}
	}

	return cfg, nil
}
