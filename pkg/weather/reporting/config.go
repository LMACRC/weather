package reporting

import (
	"fmt"
)

type BarometricMeasurementType int

const (
	BarometricMeasurementTypeAbsolute BarometricMeasurementType = iota
	BarometricMeasurementTypeRelative
)

type BarometricMeasurement struct {
	BarometricMeasurementType
}

func (b *BarometricMeasurement) UnmarshalText(text []byte) error {
	switch string(text) {
	case "relative":
		(*b).BarometricMeasurementType = BarometricMeasurementTypeRelative
	case "absolute":
		(*b).BarometricMeasurementType = BarometricMeasurementTypeAbsolute
	default:
		return fmt.Errorf("invalid barometric measurement type: %s", string(text))
	}
	return nil
}

type Config struct {
	BarometricMeasurement BarometricMeasurement `toml:"barometric_measurement"`
}
