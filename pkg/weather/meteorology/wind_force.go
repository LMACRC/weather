package meteorology

import (
	"github.com/martinlindhe/unit"
)

type WindForce int

const (
	WindForceCalm WindForce = iota + 1
	WindForceLightAir
	WindForceLightBreeze
	WindForceGentleBreeze
	WindForceModerateBreeze
	WindForceFreshBreeze
	WindForceStrongBreeze
	WindForceModerateGale
	WindForceGale
	WindForceSevereGale
	WindForceStorm
	WindForceViolentStorm
	WindForceHurricane
)

func (wf WindForce) ToInt() int { return int(wf) }

// SpeedToWindForce calculates the wind force of speed.
// See https://en.wikipedia.org/wiki/Beaufort_scale#Modern_scale
func SpeedToWindForce(speed unit.Speed) WindForce {
	mps := speed.MetersPerSecond()
	switch {
	case mps < 0.5:
		return WindForceCalm
	case mps <= 1.5:
		return WindForceLightAir
	case mps <= 3.3:
		return WindForceLightBreeze
	case mps <= 5.5:
		return WindForceGentleBreeze
	case mps <= 7.9:
		return WindForceModerateBreeze
	case mps <= 10.7:
		return WindForceFreshBreeze
	case mps <= 13.8:
		return WindForceStrongBreeze
	case mps <= 17.1:
		return WindForceModerateGale
	case mps <= 20.7:
		return WindForceGale
	case mps <= 24.4:
		return WindForceSevereGale
	case mps <= 28.4:
		return WindForceStorm
	case mps <= 32.6:
		return WindForceViolentStorm
	case mps > 32.6:
		fallthrough
	default:
		return WindForceHurricane
	}
}
