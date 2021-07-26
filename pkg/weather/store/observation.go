package store

import (
	"github.com/lmacrc/weather/pkg/sql/driver/sqlite"
	"github.com/lmacrc/weather/pkg/weather/model"
	"github.com/lmacrc/weather/pkg/xunit"
	"github.com/martinlindhe/unit"
)

type Observation struct {
	ID                 uint             `gorm:"primarykey" csv:"id"`
	Timestamp          sqlite.Timestamp `gorm:"index:idx_timestamp,sort:desc,priority:1" csv:"timestamp"`
	BarometricAbsHpa   float64          `csv:"barometric_abs_hpa"`
	BarometricRelHpa   float64          `csv:"barometric_rel_hpa"`
	HourlyRainMm       float64          `csv:"hourly_rain_mm"`
	DailyRainMm        float64          `csv:"daily_rain_mm"`
	WeeklyRainMm       float64          `csv:"weekly_rain_mm"`
	MonthlyRainMm      float64          `csv:"monthly_rain_mm"`
	TotalRainMm        float64          `csv:"total_rain_mm"`
	EventRainMm        float64          `csv:"event_rain_mm"`
	RainRatePerHourMm  float64          `csv:"rain_rate_per_hour_mm"`
	HumidityOutdoorPct float64          `csv:"humidity_outdoor_pct"`
	HumidityIndoorPct  float64          `csv:"humidity_indoor_pct"`
	WindDirDeg         float64          `csv:"wind_dir_deg"`
	WindGustKph        float64          `csv:"wind_gust_kph"`
	WindSpeedKph       float64          `csv:"wind_speed_kph"`
	MaxDailyGustKph    float64          `csv:"max_daily_gust_kph"`
	SolarRadiationWm2  float64          `csv:"solar_radiation_wm_2"`
	TempOutdoorC       float64          `csv:"temp_outdoor_c"`
	TempIndoorC        float64          `csv:"temp_indoor_c"`
	UltravioletIndex   int              `csv:"ultraviolet_index"`
}

func (m *Observation) FromObservation(wo model.Observation) {
	*m = Observation{
		ID:                 wo.ID,
		Timestamp:          sqlite.Timestamp{wo.Timestamp},
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
		SolarRadiationWm2:  wo.SolarRadiation.WattsPerSquareMetre(),
		TempOutdoorC:       wo.TempOutdoor.Celsius(),
		TempIndoorC:        wo.TempIndoor.Celsius(),
		UltravioletIndex:   wo.UltravioletIndex,
	}
}

func (m Observation) ToObservation() *model.Observation {
	return &model.Observation{
		ID:               m.ID,
		Timestamp:        m.Timestamp.Time,
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
		SolarRadiation:   xunit.Irradiance(m.SolarRadiationWm2) * xunit.WattPerSquareMetre,
		TempOutdoor:      unit.FromCelsius(m.TempOutdoorC),
		TempIndoor:       unit.FromCelsius(m.TempIndoorC),
		UltravioletIndex: m.UltravioletIndex,
	}
}
