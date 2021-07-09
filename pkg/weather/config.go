package weather

import (
	"github.com/lmacrc/weather/pkg/weather/ftp"
)

type Config struct {
	Latitude  float64
	Longitude float64
	Ftp       ftp.Config
}
