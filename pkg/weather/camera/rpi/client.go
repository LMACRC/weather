package rpi

import (
	"errors"

	"github.com/dhowden/raspicam"
	"github.com/lmacrc/weather/pkg/weather/camera"
	"github.com/spf13/viper"
)

type Config struct{}

func NewConfig() Config {
	return Config{}
}

type Client struct {
}

func init() {
	camera.Register("rpi", func(v *viper.Viper) (camera.Capturer, error) {
		return New(v)
	})
}

func New(v *viper.Viper) (*Client, error) {
	return nil, errors.New("bad config")
}

func (c *Client) Capture(params camera.CaptureParams) (string, error) {
	s := raspicam.NewStill()
	s.Width = params.Width
	s.Height = params.Height
	s.Camera.Rotation = params.Rotate.ToInt()

	errCh := make(chan error)
	raspicam.Capture(s, nil, errCh)

	return "", nil
}
