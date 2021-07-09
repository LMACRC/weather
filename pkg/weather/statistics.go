package weather

import (
	"time"

	"github.com/lmacrc/weather/pkg/weather/meteorology"
	"github.com/lmacrc/weather/pkg/xunit"
	"github.com/martinlindhe/unit"
)

type Statistics struct {
	Timestamp            time.Time             // 01 - Date dd/mm/yy
	OutdoorTemperature   unit.Temperature      // 03 - outside temperature
	OutdoorHumidity      int                   // 04 - relative humidity http://en.wikipedia.org/wiki/Relative_humidity
	DewPoint             unit.Temperature      // 05 - dew point http://en.wikipedia.org/wiki/Dewpoint
	WindSpeedAvg         unit.Speed            // 06 - wind speed (average for current 24-hour period)
	WindSpeedLast        unit.Speed            // 07 - latest wind speed reading
	WindBearing          unit.Angle            // 08 - wind bearing (degrees)
	RainRate             unit.Length           // 09 - current rain rate (per hour)
	RainfallToday        unit.Length           // 10 - rain for current 24-hour period
	BarometricPressure   unit.Pressure         // 11 - barometer (The sea level pressure)
	WindDirection        meteorology.Direction // 12 - current wind direction (compass point)
	WindForce            int                   // 13 - current wind speed as Beaufort wind force per https://en.wikipedia.org/wiki/Beaufort_scale
	WindUnits            string                // 14 - wind units – m/s, mph, km/h, kts
	TempUnits            string                // 15 - temperature units – C, F
	PressureUnits        string                // 16 - pressure units – mb, hPa, in
	RainUnits            string                // 17 - rain units – mm, in
	WindRun              unit.Length           // 18 - wind run for current 24-hour period per https://cumuluswiki.org/a/Windrun
	PressureTrend        unit.Pressure         // 19 - average rate of pressure change over the last three hours
	MonthlyRainfall      unit.Length           // 20 - monthly rainfall
	YearlyRainfall       unit.Length           // 21 - yearly rainfall
	YesterdayRainfall    unit.Length           // 22 - yesterday's rainfall
	IndoorTemp           unit.Temperature      // 23 - inside temperature
	IndoorHumidity       int                   // 24 - inside humidity http://en.wikipedia.org/wiki/Humidity
	WindChill            unit.Temperature      // 25 - Wind chill per https://en.wikipedia.org/wiki/Wind_chill
	TempTrend            unit.Temperature      // 26 - Average rate of temperature change over the last three hours
	TodayTempHi          unit.Temperature      // 27 - today's high temp
	TodayTempHiTime      time.Time             // 28 - time of today's high temp (hh:mm)
	TodayTempLo          unit.Temperature      // 29 - today's low temp
	TodayTempLoTime      time.Time             // 30 - time of today's low temp (hh:mm)
	TodayWindHi          unit.Speed            // 31 - today's high wind speed (with multiplier? https://cumuluswiki.org/a/Wind_measurement#Weather_Stations_and_Cumulus)
	TodayWindHiTime      time.Time             // 32 - time of today's high wind speed (average) (hh:mm)
	TodayWindGustHi      unit.Speed            // 33 - today's high wind gust
	TodayWindGustHiTime  time.Time             // 34 - time of today's high wind gust (hh:mm)
	TodayPressureHi      unit.Pressure         // 35 - today's high pressure
	TodayPressureHiTime  time.Time             // 36 - time of today's high pressure (hh:mm)
	TodayPressureLo      unit.Pressure         // 37 - today's low pressure
	TodayPressureLoTime  time.Time             // 38 - time of today's low pressure (hh:mm)
	CumulusVersion       string                // 39 - Cumulus Versions (the specific version in use)
	CumulusBuildNumber   int                   // 40 - Cumulus build number
	TenMinGustHi         unit.Speed            // 41 - 10-minute high gust
	HeatIndex            unit.Temperature      // 42 - Heat index https://cumuluswiki.org/a/Heat_index
	Humidex              unit.Temperature      // 43 - https://cumuluswiki.org/a/Humidex
	UVIndex              int                   // 44 - http://en.wikipedia.org/wiki/Uv_index
	Evapotranspiration   float64               // 45 - evapotranspiration today http://en.wikipedia.org/wiki/Evapotranspiration
	SolarRadiation       xunit.Irradiance      // 46 - solar radiation W/m2 http://en.wikipedia.org/wiki/Solar_radiation
	TenMinWindBearingAvg unit.Angle            // 47 - 10-minute average wind bearing (degrees)
	RainfallLastHour     unit.Length           // 48 - rainfall last hour
	ZambrettiForecast    int                   // 49 - The number of the current (Zambretti) forecast
	IsDaylight           bool                  // 50 - Flag to indicate that the location of the station is currently in daylight (1 = yes, 0 = No)
	SensorContactLost    bool                  // 51 - If station has lost contact (1 = Yes, 0 = No)
	WindDirectionAvg     meteorology.Direction // 52 - Average wind direction
	CloudBase            int                   // 53 - Cloud base
	CloudBaseUnits       string                // 54 - Cloud base units (m, ft)
	ApparentTemp         unit.Temperature      // 55 - Apparent temperature https://cumuluswiki.org/a/Apparent_temperature
	SunshineHoursToday   time.Duration         // 56 - Sunshine hours so far today
	CurrentSolarMax      xunit.Irradiance      // 57 - Current theoretical max solar radiation
	IsSunny              bool                  // 58 - Is it sunny? 1 if the sun is shining, otherwise 0 (above or below threshold) https://cumuluswiki.org/a/Cumulus.ini_(Cumulus_1)#Section:_Solar
	TempFeelsLike        unit.Temperature      // 59 - Feels Like
}
