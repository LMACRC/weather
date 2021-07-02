package realtime

import (
	"strconv"
	"time"

	"github.com/lmacrc/weather/pkg/weather"
	"github.com/lmacrc/weather/pkg/xunit"
	"github.com/martinlindhe/unit"
)

type Statistics weather.Statistics

type buffer struct {
	b []byte
}

func (b *buffer) Bytes() []byte {
	if len(b.b) > 0 && b.b[len(b.b)-1] == ' ' {
		return b.b[:len(b.b)-1]
	}
	return b.b
}

func (b *buffer) Float1(f float64) {
	b.Float(f, 1)
}

func (b *buffer) Float(f float64, prec int) {
	b.b = strconv.AppendFloat(b.b, f, 'f', prec, 64)
	b.b = append(b.b, ' ')
}

func (b *buffer) SignedFloat(f float64, prec int) {
	if f >= 0 {
		b.Append('+')
	}
	b.b = strconv.AppendFloat(b.b, f, 'f', prec, 64)
	b.b = append(b.b, ' ')
}

func (b *buffer) Int(i int) {
	b.b = strconv.AppendInt(b.b, int64(i), 10)
	b.b = append(b.b, ' ')
}

func (b *buffer) Timestamp(t time.Time, layout string) {
	b.b = t.AppendFormat(b.b, layout)
	b.b = append(b.b, ' ')
}

func (b *buffer) Date(t time.Time) {
	b.b = t.AppendFormat(b.b, "02/01/06")
	b.b = append(b.b, ' ')
}

func (b *buffer) Time(t time.Time) {
	b.b = t.AppendFormat(b.b, "15:04:05")
	b.b = append(b.b, ' ')
}

func (b *buffer) Temp(t unit.Temperature) {
	b.Float(t.Celsius(), 1)
}

func (b *buffer) Hours(t time.Duration) {
	b.Float(t.Hours(), 1)
}

func (b *buffer) Pressure(t unit.Pressure) {
	b.Float(t.Hectopascals(), 1)
}

func (b *buffer) SignedPressure(t unit.Pressure) {
	b.SignedFloat(t.Hectopascals(), 1)
}

func (b *buffer) Bearing(v unit.Angle) {
	b.Float(v.Degrees(), 0)
}

func (b *buffer) Irradiance(v xunit.Irradiance) {
	b.Float(v.WattsPerSquareMetre(), 1)
}

func (b *buffer) SpeedKph(v unit.Speed) {
	b.Float(v.KilometersPerHour(), 1)
}

func (b *buffer) LengthKm(v unit.Length) {
	b.Float(v.Kilometers(), 1)
}

func (b *buffer) LengthMm(v unit.Length) {
	b.Float(v.Millimeters(), 1)
}

func (b *buffer) Bool(v bool) {
	if v {
		b.Int(1)
	} else {
		b.Int(0)
	}
}

func (b *buffer) SignedTemp(t unit.Temperature) {
	b.SignedFloat(t.Celsius(), 1)
}

func (b *buffer) ShortTime(t time.Time) {
	b.b = t.AppendFormat(b.b, "15:04")
	b.b = append(b.b, ' ')
}

func (b *buffer) String(s string) {
	b.b = append(b.b, s...)
	b.b = append(b.b, ' ')
}

func (b *buffer) Append(d ...byte) {
	b.b = append(b.b, d...)
}

func (s Statistics) MarshalText() (text []byte, err error) {
	b := buffer{b: make([]byte, 0, 512)}

	b.Date(s.Timestamp)                                 // 01
	b.Time(s.Timestamp)                                 // 02
	b.Temp(s.OutdoorTemperature)                        // 03
	b.Int(s.OutsideHumidity)                            // 04
	b.Temp(s.DewPoint)                                  // 05
	b.SpeedKph(s.WindSpeedAvg)                          // 06
	b.SpeedKph(s.WindSpeedLast)                         // 07
	b.Bearing(s.WindBearing)                            // 08
	b.LengthMm(s.RainRate)                              // 09
	b.LengthMm(s.RainfallToday)                         // 10
	b.Pressure(s.BarometricPressure)                    // 11
	b.String(string(s.WindDirection))                   // 12
	b.Int(s.WindForce)                                  // 13
	b.String(s.WindUnits)                               // 14
	b.String(s.TempUnits)                               // 15
	b.String(s.PressureUnits)                           // 16
	b.String(s.RainUnits)                               // 17
	b.LengthKm(s.WindRun)                               // 18
	b.SignedPressure(s.PressureTrend)                   // 19
	b.LengthMm(s.MonthlyRainfall)                       // 20
	b.LengthMm(s.YearlyRainfall)                        // 21
	b.LengthMm(s.YesterdayRainfall)                     // 22
	b.Temp(s.InsideTemp)                                // 23
	b.Int(s.InsideHumidity)                             // 24
	b.Temp(s.WindChill)                                 // 25
	b.SignedTemp(s.TempTrend)                           // 26
	b.Temp(s.TodayTempHi)                               // 27
	b.ShortTime(s.TodayTempHiTime)                      // 28
	b.Temp(s.TodayTempLo)                               // 29
	b.ShortTime(s.TodayTempLoTime)                      // 30
	b.SpeedKph(s.TodayWindHi)                           // 31
	b.ShortTime(s.TodayWindHiTime)                      // 32
	b.SpeedKph(s.TodayWindGustHi)                       // 33
	b.ShortTime(s.TodayWindGustHiTime)                  // 34
	b.Pressure(s.TodayPressureHi)                       // 35
	b.ShortTime(s.TodayPressureHiTime)                  // 36
	b.Pressure(s.TodayPressureLo)                       // 37
	b.ShortTime(s.TodayPressureLoTime)                  // 38
	b.String(s.CumulusVersion)                          // 39
	b.Int(s.CumulusBuildNumber)                         // 40
	b.SpeedKph(s.TenMinGustHi)                          // 41
	b.Float1(s.HeatIndex)                               // 42
	b.Float1(s.Humidex)                                 // 43
	b.Int(s.UVIndex)                                    // 44
	b.Float1(s.Evapotranspiration)                      // 45
	b.Irradiance(s.SolarRadiation)                      // 46
	b.Bearing(s.TenMinWindBearingAvg)                   // 47
	b.LengthMm(s.RainfallLastHour)                      // 48
	b.Int(s.ZambrettiForecast)                          // 49
	b.Bool(s.IsDaylight)                                // 50
	b.Bool(s.SensorContactLost)                         // 51
	b.String(string(s.WindDirectionAvg))                // 52
	b.Int(s.CloudBase)                                  // 53
	b.String(s.CloudBaseUnits)                          // 54
	b.Temp(s.ApparentTemp)                              // 55
	b.Hours(s.SunshineHoursToday)                       // 56
	b.Int(int(s.CurrentSolarMax.WattsPerSquareMetre())) // 57
	b.Bool(s.IsSunny)                                   // 58
	b.Temp(s.TempFeelsLike)                             // 59

	return b.Bytes(), nil
}
