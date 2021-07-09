package realtime

type Config struct {
	Cron       string
	UploadPath string `toml:"upload_path"`
}
