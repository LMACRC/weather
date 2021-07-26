package archive

import (
	"fmt"
)

type Config struct {
	Enabled     bool
	Compression Compression
	LocalDir    string `toml:"local_dir" mapstructure:"local_dir"`
	RemoteDir   string `toml:"remote_dir" mapstructure:"remote_dir"`
	Filename    string
}

func NewConfig() Config {
	return Config{
		Enabled:     true,
		Compression: "gzip",
		LocalDir:    ".",
		Filename:    "archive_{{ strftime \"%Y%m%d\" .Now }}.csv",
	}
}

type Compression string

func (c *Compression) UnmarshalText(text []byte) error {
	switch string(text) {
	case "gzip":
		*c = CompressionGzip
	case "brotli":
		*c = CompressionBrotli
	default:
		return fmt.Errorf("invalid compression %s: expect brotli,gzip", string(text))
	}
	return nil
}

const (
	CompressionGzip   Compression = "gzip"
	CompressionBrotli Compression = "brotli"
)
