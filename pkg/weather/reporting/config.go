package reporting

import (
	"fmt"
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

type Config struct {
	BarometricMeasurement BarometricMeasurementType `toml:"barometric_measurement" mapstructure:"barometric_measurement"`
}

func NewConfig() Config {
	return Config{
		BarometricMeasurement: BarometricMeasurementTypeRelative,
	}
}
