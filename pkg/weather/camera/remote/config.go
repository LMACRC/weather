package remote

type Config struct {
	RemotePath string `toml:"remote_path" mapstructure:"remote_path"`
}

func NewConfig() Config {
	return Config{}
}
