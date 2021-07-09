package weather

import (
	"github.com/lmacrc/weather/pkg/weather/ftp"
	"github.com/lmacrc/weather/pkg/weather/realtime"
)

type Config struct {
	Latitude  float64
	Longitude float64
	Ftp       ftp.Config
	Realtime  realtime.Config
}
