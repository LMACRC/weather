package meteorology

import (
	"math"

	"github.com/martinlindhe/unit"
)

// ApparentTemperature calculates the apparent temperature using temp, wind and rh.
// See http://www.bom.gov.au/info/thermal_stress/#atapproximation
func ApparentTemperature(temp unit.Temperature, wind unit.Speed, rh int) unit.Temperature {
	tempC := temp.Celsius()
	vapourPressure := (float64(rh) / 100.0) * 6.105 * math.Exp(17.27*tempC/(237.7+tempC))

	return unit.FromCelsius(tempC + (0.33 * vapourPressure) - (0.7 * wind.MetersPerSecond()) - 4.0)
}
