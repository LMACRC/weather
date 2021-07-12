package camera

import (
	"github.com/spf13/viper"
)

type DriverFn func(v *viper.Viper) (Capturer, error)

type Capturer interface {
	Capture(params CaptureParams) (string, error)
}

var (
	drivers = make(map[string]DriverFn)
)

func Register(name string, driver DriverFn) {
	drivers[name] = driver
}

func hasDriver(name string) bool {
	_, ok := drivers[name]
	return ok
}

func driverList() []string {
	names := make([]string, 0, len(drivers))
	for k := range drivers {
		names = append(names, k)
	}
	return names
}
