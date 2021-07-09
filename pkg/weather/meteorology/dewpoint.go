package meteorology

import (
	"math"

	"github.com/martinlindhe/unit"
)

func dewPointGamma(temp unit.Temperature, rh int) float64 {
	const (
		b = 18.678
		c = 257.14 // °C
		d = 234.5  // °C
	)

	T := temp.Celsius()

	return math.Log(float64(rh) / 100 * math.Exp((b-T/d)*(T/(c+T))))
}

// DewPoint determines the dew point of t and rh.
// Formula based on https://en.wikipedia.org/wiki/Dew_point#Calculating_the_dew_point
func DewPoint(t unit.Temperature, rh int) unit.Temperature {
	const (
		b = 18.678
		c = 257.14 // °C
	)

	gamma := dewPointGamma(t, rh)
	v := c * gamma / (b - gamma)

	return unit.FromCelsius(v)
}
