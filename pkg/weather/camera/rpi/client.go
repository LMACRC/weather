package rpi

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dhowden/raspicam"
	"github.com/hashicorp/go-multierror"
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

func New(_ *viper.Viper) (*Client, error) {
	return &Client{}, nil
}

func (c *Client) Capture(params camera.CaptureParams) (string, error) {
	s := raspicam.NewStill()
	s.Width = params.Width
	s.Height = params.Height
	s.Camera.Rotation = params.Rotate.ToInt()

	f, err := ioutil.TempFile("", "webcam.jpg")
	if err != nil {
		return "", fmt.Errorf("create: %w", err)
	}
	defer func() { _ = f.Close() }()

	var errs error
	errCh := make(chan error)
	go func() {
		for x := range errCh {
			errs = multierror.Append(errs, x)
		}
	}()

	raspicam.Capture(s, f, errCh)

	_ = f.Close()

	if errs != nil {
		_ = os.Remove(f.Name())
		return "", errs
	}

	return f.Name(), nil
}
