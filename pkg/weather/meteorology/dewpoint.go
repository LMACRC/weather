package meteorology

import (
	"math"
)

const (
	a = 17.271
	b = 237.7 // b in units of degrees
)

func DewPoint(t, rh float64) float64 {
	return (b * gamma(t, rh)) / (a - gamma(t, rh))
}

func gamma(t, rh float64) float64 {
	return (a * t / (b + t)) + math.Log(rh/100.0)
}
