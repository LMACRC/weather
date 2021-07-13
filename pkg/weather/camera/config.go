package camera

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Cron          string
	Driver        string
	LocalDir      string        `toml:"local_dir" mapstructure:"local_dir"`
	RemoteDir     string        `toml:"remote_dir" mapstructure:"remote_dir"`
	Filename      string        `toml:"filename" mapstructure:"filename"`
	CaptureParams CaptureParams `toml:"capture_params" mapstructure:"capture_params"`
}

func NewConfig() Config {
	return Config{
		Cron:     "*/5 * * * *",
		Driver:   "rpi",
		LocalDir: ".",
		Filename: "webcam.jpg",
		CaptureParams: CaptureParams{
			Width:  640,
			Height: 480,
		},
	}
}

func Int64ToRotationHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Int64 {
			return data, nil
		}
		if t != reflect.TypeOf(Rotation(5)) {
			return data, nil
		}

		val := data.(int64)
		switch val {
		case 0:
			return Rotation000, nil
		case 90:
			return Rotation090, nil
		case 180:
			return Rotation180, nil
		case 270:
			return Rotation270, nil
		default:
			return data, fmt.Errorf("invalid rotation %d, expect [0, 90, 180, 270]", val)
		}
	}
}

type Rotation int

func (r Rotation) ToInt() int { return int(r) }

const (
	Rotation000 Rotation = 0
	Rotation090 Rotation = 90
	Rotation180 Rotation = 180
	Rotation270 Rotation = 270
)

type CaptureParams struct {
	Width, Height int
	Rotate        Rotation
}
