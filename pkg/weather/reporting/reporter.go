package reporting

import (
	"database/sql"
	"math"
	"time"

	"github.com/jinzhu/now"
	"github.com/kelvins/sunrisesunset"
	"github.com/lmacrc/weather/pkg/weather/meteorology"
	"github.com/lmacrc/weather/pkg/weather/store"
	"github.com/martinlindhe/unit"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Reporter struct {
	log       *zap.Logger
	store     *store.Store
	lat, long float64
}

func New(log *zap.Logger, store *store.Store, lat, long float64) *Reporter {
	return &Reporter{
		log:   log.With(zap.String("service", "reporter")),
		store: store,
		lat:   lat,
		long:  long,
	}
}

func (r *Reporter) Generate(ts time.Time) *Statistics {
	r.log.Info("Starting report generation.")

	s := &Statistics{
		Timestamp:          ts,
		WindUnits:          "km/h",
		TempUnits:          "C",
		PressureUnits:      "hPa",
		RainUnits:          "mm",
		CloudBaseUnits:     "m",
		CumulusVersion:     "1.8.2",
		CumulusBuildNumber: 1,
	}

	r.calcLastObservation(ts, s)
	r.calcDewPoint(ts, s)
	r.calc24HourAverages(ts, s)
	r.calcWindDirection(ts, s)
	r.calcWindForce(ts, s)
	r.calcWindRun(ts, s)
	r.calcTrends(ts, s)
	r.calcLimitsForCurrent24HourPeriod(ts, s)
	r.calcTenMinuteStats(ts, s)
	r.calcIndices(ts, s)
	r.calcRainfall(ts, s)
	r.calcIsDaylight(ts, s)
	r.calcApparentTemp(ts, s)

	r.log.Info("Completed report generation.")

	return s
}

func (r *Reporter) calcLastObservation(ts time.Time, s *Statistics) {
	o := r.store.LastObservation(ts.UTC())
	s.BarometricPressure = o.BarometricAbs
	s.RainfallLastHour = o.HourlyRain
	s.RainfallToday = o.DailyRain
	s.MonthlyRainfall = o.MonthlyRain
	s.YearlyRainfall = o.TotalRain
	s.RainRate = o.RainRatePerHour
	s.OutdoorHumidity = o.HumidityOutdoor
	s.IndoorHumidity = o.HumidityIndoor
	s.WindBearing = o.WindDir
	s.WindSpeedLast = o.WindSpeed
	s.SolarRadiation = o.SolarRadiation
	s.OutdoorTemperature = o.TempOutdoor
	s.IndoorTemp = o.TempIndoor
	s.UVIndex = o.UltravioletIndex
}

func (r *Reporter) calcDewPoint(_ time.Time, s *Statistics) {
	s.DewPoint = meteorology.DewPoint(s.OutdoorTemperature, s.OutdoorHumidity)
}

func (r *Reporter) calc24HourAverages(ts time.Time, s *Statistics) {
	db := r.store.DB().Scopes(last24Hours(ts))

	res := struct {
		Speed, Dir float64
	}{}
	db.Table("observations").Select("AVG(wind_speed_kph) as speed, AVG(wind_dir_deg) as dir").Find(&res)
	s.WindSpeedAvg = unit.Speed(res.Speed) * unit.KilometersPerHour
	s.WindDirectionAvg = meteorology.CardinalDirection(res.Dir)
}

func (r *Reporter) calcWindDirection(_ time.Time, s *Statistics) {
	s.WindDirection = meteorology.CardinalDirection(s.WindBearing.Degrees())
}

func (r *Reporter) calcWindForce(_ time.Time, s *Statistics) {
	s.WindForce = meteorology.SpeedToWindForce(s.WindSpeedAvg).ToInt()
}

func (r *Reporter) calcWindRun(ts time.Time, s *Statistics) {
	db := r.store.DB()
	subQuery := db.Scopes(last24Hours(ts)).
		Model(&store.Observation{}).
		Select(
			"strftime('%s', timestamp) - lag(strftime('%s', timestamp), -1) over (order by timestamp desc) as diff_secs",
			"wind_speed_kph * 0.277778 as \"wind_speed_mps\"",
		)
	var windRunMetres float64
	db.Table("(?) as wind_speed_values", subQuery).
		Where("diff_secs IS NOT NULL").
		Select("SUM(diff_secs * wind_speed_mps) as wind_run_m").
		Find(&windRunMetres)
	s.WindRun = unit.Length(windRunMetres) * unit.Meter
}

func (r *Reporter) calcTrends(ts time.Time, s *Statistics) {
	s.PressureTrend = unit.Pressure(r.calcLinearRegression("barometric_rel_hpa", ts, -3*time.Hour)) * unit.Hectopascal
	s.TempTrend = unit.FromCelsius(r.calcLinearRegression("temp_outdoor_c", ts, -3*time.Hour))
}

func (r *Reporter) calcLinearRegression(col string, now time.Time, d time.Duration) float64 {
	db := r.store.DB()
	start := now.Add(d)
	end := now

	// dependent variable:   col
	// independent variable: timestamp (seconds)

	subQuery := db.Model(&store.Observation{}).
		Where("timestamp >= ? AND timestamp <= ?", start.UTC(), end.UTC()).
		Select("AVG("+col+") over () as ybar, "+col+" as y, AVG((STRFTIME('%s', timestamp) - @start)) OVER () as xbar, (STRFTIME('%s', timestamp) - @start) as x",
			sql.Named("start", start.Unix()))

	var slope sql.NullFloat64
	db.Table("(?) as calculations", subQuery).
		Select("SUM((x-xbar)*(y-ybar)) / SUM((x-xbar)*(x-xbar)) as beta").
		Find(&slope)

	// Slope is in column_units / seconds, therefore we adjust slope to establish the trend over the entire duration
	if slope.Valid {
		return slope.Float64 * math.Abs(d.Seconds())
	}

	return 0
}

func (r *Reporter) calcRainfall(ts time.Time, s *Statistics) {
	// limit for previous hour, as hourly_rain_mm resets to zero after each hour
	_, val := r.calcLimitAndTimeForPeriod("hourly_rain_mm", limitMax, now.With(ts).BeginningOfHour(), -1*time.Hour)
	s.RainfallLastHour = unit.Length(val) * unit.Millimeter
}

func (r *Reporter) calcLimitsForCurrent24HourPeriod(ts time.Time, s *Statistics) {
	start := now.With(ts).BeginningOfDay()
	dur := 24 * time.Hour
	var val float64

	s.TodayTempHiTime, val = r.calcLimitAndTimeForPeriod("temp_outdoor_c", limitMax, start, dur)
	s.TodayTempHi = unit.FromCelsius(val)
	s.TodayTempLoTime, val = r.calcLimitAndTimeForPeriod("temp_outdoor_c", limitMin, start, dur)
	s.TodayTempLo = unit.FromCelsius(val)
	s.TodayWindHiTime, val = r.calcLimitAndTimeForPeriod("wind_speed_kph", limitMax, start, dur)
	s.TodayWindHi = unit.Speed(val) * unit.KilometersPerHour
	s.TodayWindGustHiTime, val = r.calcLimitAndTimeForPeriod("wind_gust_kph", limitMax, start, dur)
	s.TodayWindGustHi = unit.Speed(val) * unit.KilometersPerHour
	s.TodayPressureHiTime, val = r.calcLimitAndTimeForPeriod("barometric_abs_hpa", limitMax, start, dur)
	s.TodayPressureHi = unit.Pressure(val) * unit.Hectopascal
	s.TodayPressureLoTime, val = r.calcLimitAndTimeForPeriod("barometric_abs_hpa", limitMin, start, dur)
	s.TodayPressureLo = unit.Pressure(val) * unit.Hectopascal
}

type limit int

const (
	limitMin limit = iota
	limitMax
)

func (r *Reporter) calcLimitAndTimeForPeriod(col string, limit limit, now time.Time, d time.Duration) (time.Time, float64) {
	var (
		start, end time.Time
	)

	if d < 0 {
		start = now.Add(d)
		end = now
	} else {
		start = now
		end = now.Add(d)
	}

	var res struct {
		Timestamp time.Time
		Value     float64
	}
	order := clause.OrderByColumn{Column: clause.Column{Name: col}}
	if limit == limitMax {
		order.Desc = true
	}

	r.store.DB().Model(&store.Observation{}).
		Where("timestamp >= ? AND timestamp <= ?", start.UTC(), end.UTC()).
		Select("timestamp, " + col + " as value").
		Order(order).Order("timestamp").
		Find(&res)

	return res.Timestamp.In(now.Location()), res.Value
}

func (r *Reporter) calcTenMinuteStats(ts time.Time, s *Statistics) {
	s.TenMinGustHi = unit.Speed(r.calcStatForPeriod("wind_gust_kph", "MAX", ts, -10*time.Minute)) * unit.KilometersPerHour
	s.TenMinWindBearingAvg = unit.Angle(r.calcStatForPeriod("wind_dir_deg", "AVG", ts, -10*time.Minute)) * unit.Degree
}

func (r *Reporter) calcStatForPeriod(col, stat string, now time.Time, d time.Duration) float64 {
	var (
		start, end time.Time
	)

	if d < 0 {
		start = now.Add(d)
		end = now
	} else {
		start = now
		end = now.Add(d)
	}

	var res float64

	r.store.DB().Model(&store.Observation{}).
		Where("timestamp >= ? AND timestamp <= ?", start.UTC(), end.UTC()).
		Select(stat + "(" + col + ") as value").
		Find(&res)

	return res
}

func (r *Reporter) calcIndices(_ time.Time, s *Statistics) {
	s.HeatIndex = meteorology.HeatIndex(s.OutdoorTemperature, s.OutdoorHumidity)
	s.Humidex = meteorology.Humidex(s.OutdoorTemperature, s.OutdoorHumidity)

}

func (r *Reporter) calcIsDaylight(ts time.Time, s *Statistics) {
	_, ofs := ts.Zone()
	utcOffset := time.Duration(ofs) * time.Second
	p := sunrisesunset.Parameters{
		Latitude:  r.lat,
		Longitude: r.long,
		UtcOffset: utcOffset.Hours(),
		Date:      ts,
	}
	sunrise, sunset, err := p.GetSunriseSunset()
	if err != nil {
		return
	}
	s.IsDaylight = ts.After(sunrise) && ts.Before(sunset)
}

func (r *Reporter) calcApparentTemp(_ time.Time, s *Statistics) {
	// Wind chill in Australia also uses the apparent temperature
	//   per https://en.wikipedia.org/wiki/Wind_chill
	s.WindChill = meteorology.ApparentTemperature(s.OutdoorTemperature, s.WindSpeedLast, s.OutdoorHumidity)
	s.ApparentTemp = meteorology.ApparentTemperature(s.OutdoorTemperature, s.WindSpeedLast, s.OutdoorHumidity)
	s.TempFeelsLike = meteorology.ApparentTemperature(s.OutdoorTemperature, s.WindSpeedLast, s.OutdoorHumidity)
}

func last24Hours(d time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		lo := now.With(d).BeginningOfDay()
		hi := now.With(d).BeginningOfDay().Add(24 * time.Hour)
		return db.Where("timestamp >= ? AND timestamp < ?", lo.UTC(), hi.UTC())
	}
}
