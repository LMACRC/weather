package influxdb

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server          string
	Database        string
	RetentionPolicy string `toml:"retention_policy" mapstructure:"retention_policy"`
	Measurement     string
}

func NewConfig() Config {
	return Config{
		Server:          "http://localhost:8086",
		RetentionPolicy: "auto",
		Measurement:     "observations",
	}
}

func InitViper(v *viper.Viper) {
	cfg := NewConfig()
	v.SetDefault("influxdb.server", cfg.Server)
	v.SetDefault("influxdb.retention_policy", cfg.RetentionPolicy)
	v.SetDefault("influxdb.measurement", cfg.Measurement)
}
