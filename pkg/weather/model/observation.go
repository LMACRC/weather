package model

import (
	"time"

	"github.com/lmacrc/weather/pkg/xunit"
	"github.com/martinlindhe/unit"
)

type Observation struct {
	ID               uint
	Timestamp        time.Time
	BarometricAbs    unit.Pressure
	BarometricRel    unit.Pressure
	HourlyRain       unit.Length
	DailyRain        unit.Length
	WeeklyRain       unit.Length
	MonthlyRain      unit.Length
	TotalRain        unit.Length
	EventRain        unit.Length
	RainRatePerHour  unit.Length
	HumidityOutdoor  int
	HumidityIndoor   int
	WindDir          unit.Angle
	WindGust         unit.Speed
	WindSpeed        unit.Speed
	MaxDailyGust     unit.Speed
	SolarRadiation   xunit.Irradiance
	TempOutdoor      unit.Temperature
	TempIndoor       unit.Temperature
	UltravioletIndex int
}
