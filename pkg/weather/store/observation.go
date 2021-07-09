package store

import (
	"time"

	"github.com/lmacrc/weather/pkg/weather"
	"github.com/lmacrc/weather/pkg/xunit"
	"github.com/martinlindhe/unit"
)

type Observation struct {
	ID                 uint      `gorm:"primarykey"`
	Timestamp          time.Time `gorm:"index:idx_ts_wind_dir,sort:desc,priority:1;index:index_ts_wind_gust,sort:desc,priority:1;index:idx_ts_wind_speed,sort:desc,priority:1;index:idx_ts_outdoor_temp,sort:desc,priority:1;index:idx_ts_indoor_temp,sort:desc,priority:1;index:idx_ts_barom_abs,sort:desc,priority:1"`
	BarometricAbsHpa   float64   `gorm:"index:idx_ts_barom_abs,sort:desc,priority:2"`
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
	WindDirDeg         float64 `gorm:"index:idx_ts_wind_dir,sort:desc,priority:2"`
	WindGustKph        float64 `gorm:"index:idx_ts_wind_gust,sort:desc,priority:2"`
	WindSpeedKph       float64 `gorm:"index:idx_ts_wind_speed,sort:desc,priority:2"`
	MaxDailyGustKph    float64
	Model              string
	StationType        string
	SolarRadiationWm2  float64
	TempOutdoorC       float64 `gorm:"index:idx_ts_outdoor_temp,sort:desc,priority:2"`
	TempIndoorC        float64 `gorm:"index:idx_ts_indoor_temp,sort:desc,priority:2"`
	UltravioletIndex   int
}

func (m *Observation) FromObservation(wo weather.Observation) {
	*m = Observation{
		ID:                 wo.ID,
		Timestamp:          wo.Timestamp,
		BarometricAbsHpa:   wo.BarometricAbs.Hectopascals(),
		BarometricRelHpa:   wo.BarometricRel.Hectopascals(),
		HourlyRainMm:       wo.HourlyRain.Millimeters(),
		DailyRainMm:        wo.DailyRain.Millimeters(),
		WeeklyRainMm:       wo.WeeklyRain.Millimeters(),
		MonthlyRainMm:      wo.MonthlyRain.Millimeters(),
		TotalRainMm:        wo.TotalRain.Millimeters(),
		EventRainMm:        wo.EventRain.Millimeters(),
		RainRatePerHourMm:  wo.RainRatePerHour.Millimeters(),
		HumidityOutdoorPct: float64(wo.HumidityOutdoor) / 100.0,
		HumidityIndoorPct:  float64(wo.HumidityIndoor) / 100.0,
		WindDirDeg:         wo.WindDir.Degrees(),
		WindGustKph:        wo.WindGust.KilometersPerHour(),
		WindSpeedKph:       wo.WindSpeed.KilometersPerHour(),
		MaxDailyGustKph:    wo.MaxDailyGust.KilometersPerHour(),
		Model:              wo.Model,
		StationType:        wo.StationType,
		SolarRadiationWm2:  wo.SolarRadiation.WattsPerSquareMetre(),
		TempOutdoorC:       wo.TempOutdoor.Celsius(),
		TempIndoorC:        wo.TempIndoor.Celsius(),
		UltravioletIndex:   wo.UltravioletIndex,
	}
}

func (m Observation) ToObservation() *weather.Observation {
	return &weather.Observation{
		ID:               m.ID,
		Timestamp:        m.Timestamp,
		BarometricAbs:    unit.Pressure(m.BarometricAbsHpa) * unit.Hectopascal,
		BarometricRel:    unit.Pressure(m.BarometricRelHpa) * unit.Hectopascal,
		HourlyRain:       unit.Length(m.HourlyRainMm) * unit.Millimeter,
		DailyRain:        unit.Length(m.DailyRainMm) * unit.Millimeter,
		WeeklyRain:       unit.Length(m.WeeklyRainMm) * unit.Millimeter,
		MonthlyRain:      unit.Length(m.MonthlyRainMm) * unit.Millimeter,
		TotalRain:        unit.Length(m.TotalRainMm) * unit.Millimeter,
		EventRain:        unit.Length(m.EventRainMm) * unit.Millimeter,
		RainRatePerHour:  unit.Length(m.RainRatePerHourMm) * unit.Millimeter,
		HumidityOutdoor:  int(m.HumidityOutdoorPct * 100),
		HumidityIndoor:   int(m.HumidityIndoorPct * 100),
		WindDir:          unit.Angle(m.WindDirDeg) * unit.Degree,
		WindGust:         unit.Speed(m.WindGustKph) * unit.KilometersPerHour,
		WindSpeed:        unit.Speed(m.WindSpeedKph) * unit.KilometersPerHour,
		MaxDailyGust:     unit.Speed(m.MaxDailyGustKph) * unit.KilometersPerHour,
		Model:            m.Model,
		StationType:      m.StationType,
		SolarRadiation:   xunit.Irradiance(m.SolarRadiationWm2) * xunit.WattPerSquareMetre,
		TempOutdoor:      unit.FromCelsius(m.TempOutdoorC),
		TempIndoor:       unit.FromCelsius(m.TempIndoorC),
		UltravioletIndex: m.UltravioletIndex,
	}
}
