package reporting

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type BarometricMeasurementType int

const (
	BarometricMeasurementTypeAbsolute BarometricMeasurementType = iota
	BarometricMeasurementTypeRelative
)

func (b *BarometricMeasurementType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "relative":
		*b = BarometricMeasurementTypeRelative
	case "absolute":
		*b = BarometricMeasurementTypeAbsolute
	default:
		return fmt.Errorf("invalid barometric measurement type: %s", string(text))
	}
	return nil
}

func StringToBarometricMeasurementHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(BarometricMeasurementType(5)) {
			return data, nil
		}

		var r BarometricMeasurementType
		return r, r.UnmarshalText([]byte(data.(string)))
	}
}

type Config struct {
	BarometricMeasurement BarometricMeasurementType `toml:"barometric_measurement" mapstructure:"barometric_measurement"`
}

func NewConfig() Config {
	return Config{
		BarometricMeasurement: BarometricMeasurementTypeRelative,
	}
}
