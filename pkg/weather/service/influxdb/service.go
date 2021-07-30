package influxdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/influxdata/influxdb1-client/v2"
	"github.com/lmacrc/weather/pkg/event"
	"github.com/lmacrc/weather/pkg/weather/reporting"
	"github.com/lmacrc/weather/pkg/weather/service/realtime"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Service struct {
	log    *zap.Logger
	client client.Client

	database    string
	rp          string
	measurement string
}

func New(log *zap.Logger, v *viper.Viper, bus *event.Bus) (*Service, error) {
	var cfg Config
	if err := v.UnmarshalKey("influxdb", &cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	if cfg.Database == "" {
		return nil, errors.New("database cannot be empty")
	}

	if cfg.Measurement == "" {
		return nil, errors.New("measurement cannot be empty")
	}

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: cfg.Server,
	})
	if err != nil {
		return nil, err
	}

	s := &Service{
		log:         log.With(zap.String("service", "influxdb")),
		client:      influxClient,
		database:    cfg.Database,
		rp:          cfg.RetentionPolicy,
		measurement: cfg.Measurement,
	}

	bus.MustSubscribe(realtime.NewStatistics, s.HandleStatistics)

	return s, nil
}

func (s *Service) Run(ctx context.Context) {
	s.log.Info("Starting.")
	<-ctx.Done()
	s.log.Info("Shutting down.")
}

func (s *Service) HandleStatistics(stats *reporting.Statistics) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision:       "s",
		Database:        s.database,
		RetentionPolicy: s.rp,
	})
	if err != nil {
		return
	}

	boolToString := func(b bool) string {
		if b {
			return "Y"
		}
		return "N"
	}

	tags := map[string]string{
		"wind_units":     stats.WindUnits,
		"temp_units":     stats.TempUnits,
		"pressure_units": stats.PressureUnits,
		"rain_units":     stats.RainUnits,
	}

	fields := map[string]interface{}{
		"outdoor_temperature_c":   stats.OutdoorTemperature.Celsius(),
		"outdoor_humidity":        stats.OutdoorHumidity,
		"dew_point_c":             stats.DewPoint.Celsius(),
		"wind_speed_last_mps":     stats.WindSpeedLast.MetersPerSecond(),
		"wind_bearing_deg":        stats.WindBearing.Degrees(),
		"rain_rate_mm_per_hour":   stats.RainRate.Millimeters(),
		"barometric_pressure_hpa": stats.BarometricPressure.Hectopascals(),
		"wind_force":              stats.WindForce,
		"wind_run_m":              stats.WindRun.Meters(),
		"indoor_temp_c":           stats.IndoorTemp.Celsius(),
		"indoor_humidity":         stats.IndoorHumidity,
		"wind_chill_c":            stats.WindChill.Celsius(),
		"ten_min_gust_hi_mps":     stats.TenMinGustHi.MetersPerSecond(),
		"heat_index_c":            stats.HeatIndex.Celsius(),
		"humidex_c":               stats.Humidex.Celsius(),
		"uv_index":                stats.UVIndex,
		"solar_radiation_wsm":     stats.SolarRadiation.WattsPerSquareMetre(),
		"apparent_temp_c":         stats.ApparentTemp.Celsius(),
		"wind_direction":          stats.WindDirection.String(),
		"is_daylight":             boolToString(stats.IsDaylight),
	}

	p, err := client.NewPoint(s.measurement, tags, fields, stats.Timestamp)
	if err != nil {
		return
	}

	bp.AddPoint(p)

	err = s.client.Write(bp)
	if err != nil {
		return
	}
}
