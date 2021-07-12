package ftp

import (
	"context"
	"fmt"
	"io"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
)

type Client struct {
	config Config
}

func New(v *viper.Viper) (*Client, error) {
	var cfg Config
	if err := v.UnmarshalKey("ftp", &cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return &Client{
		config: cfg,
	}, nil
}

func (s *Client) Upload(ctx context.Context, dir, filename string, r io.Reader) error {
	conn, err := ftp.Dial(s.config.Address, ftp.DialWithContext(ctx))
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	defer func() { _ = conn.Quit() }()

	if err := conn.Login(s.config.Username, s.config.Password); err != nil {
		return fmt.Errorf("login error: %w", err)
	}

	if err := conn.ChangeDir(dir); err != nil {
		return fmt.Errorf("unable to change to directory %q: %w", dir, err)
	}

	if err := conn.Stor(filename, r); err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	return nil
}
