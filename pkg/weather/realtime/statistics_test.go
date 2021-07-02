package realtime

import (
	"testing"
	"time"

	"github.com/lmacrc/weather/pkg/weather/meteorology"
	"github.com/martinlindhe/unit"
	"github.com/stretchr/testify/assert"
)

func TestStatistics_MarshalText(t *testing.T) {
	stats := Statistics{
		Timestamp:            time.Date(2008, 10, 18, 16, 3, 45, 0, time.UTC),
		OutdoorTemperature:   unit.FromCelsius(8.4),
		OutsideHumidity:      84,
		DewPoint:             unit.FromCelsius(5.8),
		WindSpeedAvg:         24.2 * unit.KilometersPerHour,
		WindSpeedLast:        33.0 * unit.KilometersPerHour,
		WindBearing:          261 * unit.Degree,
		RainRate:             0.0 * unit.Millimeter,
		RainfallToday:        1.0 * unit.Millimeter,
		BarometricPressure:   999.7 * unit.Hectopascal,
		WindDirection:        meteorology.CardinalDirection(270),
		WindForce:            6,
		WindUnits:            "kph",
		TempUnits:            "C",
		PressureUnits:        "hPa",
		RainUnits:            "mm",
		WindRun:              146.6 * unit.Kilometer,
		PressureTrend:        0.1 * unit.Hectopascal,
		MonthlyRainfall:      85.2 * unit.Millimeter,
		YearlyRainfall:       588.4 * unit.Millimeter,
		YesterdayRainfall:    11.6 * unit.Millimeter,
		InsideTemp:           unit.FromCelsius(20.3),
		InsideHumidity:       57,
		WindChill:            unit.FromCelsius(3.6),
		TempTrend:            unit.FromCelsius(-0.7),
		TodayTempHi:          unit.FromCelsius(10.9),
		TodayTempHiTime:      time.Date(2008, 10, 18, 12, 0, 20, 0, time.UTC),
		TodayTempLo:          unit.FromCelsius(7.8),
		TodayTempLoTime:      time.Date(2008, 10, 18, 14, 41, 30, 0, time.UTC),
		TodayWindHi:          37.4 * unit.KilometersPerHour,
		TodayWindHiTime:      time.Date(2008, 10, 18, 14, 38, 45, 0, time.UTC),
		TodayWindGustHi:      44.0 * unit.KilometersPerHour,
		TodayWindGustHiTime:  time.Date(2008, 10, 18, 14, 28, 11, 0, time.UTC),
		TodayPressureHi:      999.8 * unit.Hectopascal,
		TodayPressureHiTime:  time.Date(2008, 10, 18, 16, 1, 12, 0, time.UTC),
		TodayPressureLo:      998.4 * unit.Hectopascal,
		TodayPressureLoTime:  time.Date(2008, 10, 18, 12, 6, 3, 0, time.UTC),
		CumulusVersion:       "1.8.2",
		CumulusBuildNumber:   448,
		TenMinGustHi:         36.0 * unit.KilometersPerHour,
		HeatIndex:            10.3,
		Humidex:              10.5,
		UVIndex:              1,
		Evapotranspiration:   1,
		SolarRadiation:       1,
		TenMinWindBearingAvg: 234 * unit.Degree,
		RainfallLastHour:     2.5 * unit.Millimeter,
		ZambrettiForecast:    5,
		IsDaylight:           true,
		SensorContactLost:    false,
		WindDirectionAvg:     "NNW",
		CloudBase:            2040,
		CloudBaseUnits:       "ft",
		ApparentTemp:         unit.FromCelsius(12.3),
		SunshineHoursToday:   time.Duration(684) * time.Minute,
		CurrentSolarMax:      420,
		IsSunny:              true,
		TempFeelsLike:        unit.FromCelsius(8.4),
	}

	got, err := stats.MarshalText()
	assert.NoError(t, err)
	const exp = `18/10/08 16:03:45 8.4 84 5.8 24.2 33.0 261 0.0 1.0 999.7 W 6 kph C hPa mm 146.6 +0.1 85.2 588.4 11.6 20.3 57 3.6 -0.7 10.9 12:00 7.8 14:41 37.4 14:38 44.0 14:28 999.8 16:01 998.4 12:06 1.8.2 448 36.0 10.3 10.5 1 1.0 1.0 234 2.5 5 1 0 NNW 2040 ft 12.3 11.4 420 1 8.4`
	assert.Equal(t, exp, string(got))
}
