package weather

import (
	"embed"
)

//go:embed etc/weather.service
var Content embed.FS
