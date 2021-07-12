package remote

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lmacrc/weather/pkg/weather/camera"
	"github.com/spf13/viper"
)

type Client struct {
	url url.URL
	cfg Config
}

func init() {
	camera.Register("remote", func(v *viper.Viper) (camera.Capturer, error) {
		return New(v)
	})
}

func New(v *viper.Viper) (*Client, error) {
	var cfg Config
	if err := v.UnmarshalKey("camera_driver.remote", &cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	r, err := url.Parse(cfg.RemotePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing remote path: %w", err)
	}

	r.RawQuery = ""

	return &Client{
		url: *r,
		cfg: cfg,
	}, nil
}

func (c Client) Capture(params camera.CaptureParams) (string, error) {
	qry := url.Values{
		"w": []string{strconv.Itoa(params.Width)},
		"h": []string{strconv.Itoa(params.Height)},
		"r": []string{strconv.Itoa(params.Rotate.ToInt())},
	}

	req := c.url
	req.RawQuery = qry.Encode()

	res, err := http.Get(req.String())
	if err != nil {
		return "", fmt.Errorf("get: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	f, err := ioutil.TempFile("", "webcam.jpg")
	if err != nil {
		return "", fmt.Errorf("create failed: %w", err)
	}
	defer func() { _ = f.Close() }()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return "", fmt.Errorf("copy failed for path %q: %w", f.Name(), err)
	}

	return f.Name(), nil
}
