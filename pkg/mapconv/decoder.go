package mapconv

import (
	"reflect"
	"strconv"

	"github.com/lmacrc/weather/pkg/xunit"
	"github.com/martinlindhe/unit"
	"github.com/mitchellh/mapstructure"
)

func StringToLengthHookFunc(base unit.Length) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(unit.Length(5)) {
			return data, nil
		}

		v, err := strconv.ParseFloat(data.(string), 64)
		if err != nil {
			return nil, err
		}
		return unit.Length(v) * base, nil
	}
}

func StringToPressureHookFunc(base unit.Pressure) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(unit.Pressure(5)) {
			return data, nil
		}

		v, err := strconv.ParseFloat(data.(string), 64)
		if err != nil {
			return nil, err
		}
		return unit.Pressure(v) * base, nil
	}
}

func StringToSpeedHookFunc(base unit.Speed) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(unit.Speed(5)) {
			return data, nil
		}

		v, err := strconv.ParseFloat(data.(string), 64)
		if err != nil {
			return nil, err
		}
		return unit.Speed(v) * base, nil
	}
}

func StringToAngleFunc(base unit.Angle) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(unit.Angle(5)) {
			return data, nil
		}

		v, err := strconv.ParseFloat(data.(string), 64)
		if err != nil {
			return nil, err
		}
		return unit.Angle(v) * base, nil
	}
}

func StringToIrradianceFunc(base xunit.Irradiance) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(xunit.Irradiance(5)) {
			return data, nil
		}

		v, err := strconv.ParseFloat(data.(string), 64)
		if err != nil {
			return nil, err
		}
		return xunit.Irradiance(v) * base, nil
	}
}

func StringToTemperatureHookFunc(conv func(t float64) unit.Temperature) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(unit.Temperature(5)) {
			return data, nil
		}

		v, err := strconv.ParseFloat(data.(string), 64)
		if err != nil {
			return nil, err
		}
		return conv(v), nil
	}
}
