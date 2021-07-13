package meteorology

import (
	"math"

	"github.com/martinlindhe/unit"
)

// Humidex calculates the humidity index using temp and rh.
// See https://en.wikipedia.org/wiki/Humidex#Computation_formula
func Humidex(temp unit.Temperature, rh int) unit.Temperature {
	dew := DewPoint(temp, rh)

	tempAir := temp.Celsius()
	H := tempAir + 5/9*(6.11*math.Exp(5417.7530*(1/273.16-1/(273.15+dew.Celsius())))-10)
	return unit.FromCelsius(H)
}
