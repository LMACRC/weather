package realtime

type Config struct {
	Cron       string
	UploadPath string `toml:"upload_path" mapstructure:"upload_path"`
}

func NewConfig() Config {
	return Config{
		Cron: "*/5 * * * *",
	}
}
