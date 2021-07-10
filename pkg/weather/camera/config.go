package camera

type Config struct {
	Cron           string
	UploadPath     string `toml:"upload_path"`
	UploadFilename string `toml:"upload_filename"`
}
