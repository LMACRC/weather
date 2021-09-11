package camera

import (
	"fmt"
	"image/color"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	Cron          string
	Driver        string
	LocalDir      string        `toml:"local_dir" mapstructure:"local_dir"`
	RemoteDir     string        `toml:"remote_dir" mapstructure:"remote_dir"`
	Filename      string        `toml:"filename" mapstructure:"filename"`
	CaptureParams CaptureParams `toml:"capture_params" mapstructure:"capture_params"`
	OutputParams  OutputParams  `toml:"output_params" mapstructure:"output_params"`
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
		OutputParams: OutputParams{
			TextColor: Color{R: 255, G: 100, B: 0},
		},
	}
}

func InitViper(v *viper.Viper) {
	cfg := NewConfig()
	v.SetDefault("camera.cron", cfg.Cron)
	v.SetDefault("camera.driver", cfg.Driver)
	v.SetDefault("camera.local_dir", cfg.LocalDir)

	v.SetDefault("camera.capture_params.width", cfg.CaptureParams.Width)
	v.SetDefault("camera.capture_params.height", cfg.CaptureParams.Height)
	v.SetDefault("camera.capture_params.rotate", cfg.CaptureParams.Rotate)

	v.SetDefault("camera.output_params.text_color", cfg.OutputParams.TextColor.String())
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

func ParseHexColor(s string) (col Color, err error) {
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &col.R, &col.G, &col.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &col.R, &col.G, &col.B)
		col.R *= 17
		col.G *= 17
		col.B *= 17
	default:
		err = fmt.Errorf("invalid color: length must be 7 or 4 characters")
	}
	return
}

type OutputParams struct {
	TextColor Color `toml:"text_color" mapstructure:"text_color"`
}

type Color struct {
	R, G, B byte
}

func (c *Color) UnmarshalText(text []byte) error {
	cc, err := ParseHexColor(string(text))
	if err != nil {
		return err
	}
	*c = cc
	return nil
}

func (c Color) ToRGBA() color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: 0xff}
}

func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}
