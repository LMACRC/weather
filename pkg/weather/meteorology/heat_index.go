package meteorology

import (
	"github.com/martinlindhe/unit"
)

// HeatIndex calculates the heat index (HI) using temp and rh.
// Based on https://en.wikipedia.org/wiki/Heat_index#Formula
func HeatIndex(temp unit.Temperature, rh int) unit.Temperature {
	tempC := temp.Celsius()

	if tempC < 27 {
		return temp
	}

	T := tempC
	R := float64(rh)

	T2 := tempC * tempC
	R2 := float64(rh * rh)

	// coefficients for °C
	const (
		c1 = -8.78469475556
		c2 = 1.61139411
		c3 = 2.33854883889
		c4 = -0.14611605
		c5 = -0.012308094
		c6 = -0.0164248277778
		c7 = 0.002211732
		c8 = 0.00072546
		c9 = -0.000003582
	)

	resultC := c1 + c2*T + c3*R + c4*T*R + c5*T2 + c6*R2 + c7*T2*R + c8*T*R2 + c9*T2*R2

	return unit.FromCelsius(resultC)
}
