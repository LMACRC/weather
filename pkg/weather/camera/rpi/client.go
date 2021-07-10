package rpi

import (
	"github.com/dhowden/raspicam"
)

type Client struct {
}

func (c *Client) Capture() {
	s := raspicam.NewStill()

	errCh := make(chan error)
	raspicam.Capture(s, nil, errCh)
}
