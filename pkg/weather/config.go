package weather

import (
	"fmt"
	"os"

	"github.com/lmacrc/weather/pkg/weather/camera"
	"github.com/lmacrc/weather/pkg/weather/camera/remote"
	"github.com/lmacrc/weather/pkg/weather/camera/rpi"
	"github.com/lmacrc/weather/pkg/weather/ftp"
	"github.com/lmacrc/weather/pkg/weather/realtime"
	"github.com/lmacrc/weather/pkg/weather/reporting"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type Config struct {
	DbPath   string `toml:"database_path" mapstructure:"database_path"`
	Location Location

	Ftp       ftp.Config
	Realtime  realtime.Config
	Reporting reporting.Config

	Camera camera.Config

	CameraDriver struct {
		Remote remote.Config
		Rpi    rpi.Config
	} `toml:"camera_driver" mapstructure:"camera_driver"`
}

func NewConfig() Config {
	return Config{
		DbPath:    "weather.db",
		Location:  Location{},
		Ftp:       ftp.NewConfig(),
		Realtime:  realtime.NewConfig(),
		Reporting: reporting.NewConfig(),
		Camera:    camera.NewConfig(),
		CameraDriver: struct {
			Remote remote.Config
			Rpi    rpi.Config
		}{
			Remote: remote.NewConfig(),
			Rpi:    rpi.NewConfig(),
		},
	}
}

func decoderConfig(dc *mapstructure.DecoderConfig) {
	dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(dc.DecodeHook,
		camera.Int64ToRotationHookFunc(),
		reporting.StringToBarometricMeasurementHookFunc(),
	)
}

func init() {
	viper.SetConfigName("weather")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.config/weather")
	viper.AddConfigPath("/etc/weather")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./local")
}

func ReadConfig(path string) error {
	if path == "" {
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("read config: %w", err)
		}
	} else {
		fi, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer func() { _ = fi.Close() }()

		if err := viper.ReadConfig(fi); err != nil {
			return fmt.Errorf("read config: %w", err)
		}
	}

	return nil
}
