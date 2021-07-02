package weather

import (
	"time"
)

type Observation struct {
	ID                 uint      `gorm:"primarykey"`
	Timestamp          time.Time `gorm:"index"`
	BarometricAbsHpa   float64
	BarometricRelHpa   float64
	HourlyRainMm       float64
	DailyRainMm        float64
	WeeklyRainMm       float64
	MonthlyRainMm      float64
	TotalRainMm        float64
	EventRainMm        float64
	RainRatePerHourMm  float64
	HumidityOutdoorPct float64
	HumidityIndoorPct  float64
	WindDirDeg         float64
	WindGustKph        float64
	WindSpeedKph       float64
	MaxDailyGustKph    float64
	Model              string
	StationType        string
	SolarRadiationWm2  float64
	TempOutdoorC       float64
	TempIndoorC        float64
	UltravioletIndex   int
}
