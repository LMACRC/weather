package realtime

type Config struct {
	Cron      string
	RemoteDir string `toml:"remote_dir" mapstructure:"remote_dir"`
}

func NewConfig() Config {
	return Config{
		Cron: "*/5 * * * *",
	}
}
