package xunit

import (
	"github.com/martinlindhe/unit"
)

type Irradiance unit.Unit

const (
	WattPerSquareMetre = 1e0
)

func (v Irradiance) WattsPerSquareMetre() float64 { return float64(v) }
