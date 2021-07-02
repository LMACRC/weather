package service

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/lmacrc/weather/pkg/mapconv"
	"github.com/lmacrc/weather/pkg/weather"
	"github.com/lmacrc/weather/pkg/xunit"
	"github.com/martinlindhe/unit"
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	decoderHookFn = mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeHookFunc("2006-01-02 15:04:05"),
		mapconv.StringToLengthHookFunc(unit.Inch),
		mapconv.StringToPressureHookFunc(unit.InchOfMercury),
		mapconv.StringToSpeedHookFunc(unit.MilesPerHour),
		mapconv.StringToAngleFunc(unit.Degree),
		mapconv.StringToIrradianceFunc(xunit.WattPerSquareMetre),
		mapconv.StringToTemperatureHookFunc(unit.FromFahrenheit),
	)
)

var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "weather",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "The total number of processed weather requests",
	}, []string{"status"})
)

type ObservationWriter interface {
	Write(value interface{}) (int64, error)
}

type Handler struct {
	Path  string
	Store ObservationWriter
}

func (h *Handler) Handle(mux *http.ServeMux) {
	mux.Handle(h.Path, h)
}

type ecowitt struct {
	Timestamp        time.Time        `mapstructure:"dateutc"`
	BarometricAbs    unit.Pressure    `mapstructure:"baromabsin"`
	BarometricRel    unit.Pressure    `mapstructure:"baromrelin"`
	HourlyRain       unit.Length      `mapstructure:"hourlyrainin,inch"`
	DailyRain        unit.Length      `mapstructure:"dailyrainin"`
	WeeklyRain       unit.Length      `mapstructure:"weeklyrainin"`
	MonthlyRain      unit.Length      `mapstructure:"monthlyrainin"`
	TotalRain        unit.Length      `mapstructure:"totalrainin"`
	EventRain        unit.Length      `mapstructure:"eventrainin"`
	RainRatePerHour  unit.Length      `mapstructure:"rainratein"`
	HumidityOutdoor  int              `mapstructure:"humidity"`
	HumidityIndoor   int              `mapstructure:"humidityin"`
	WindDir          unit.Angle       `mapstructure:"winddir"`
	WindGust         unit.Speed       `mapstructure:"windgustmph"`
	WindSpeed        unit.Speed       `mapstructure:"windspeedmph"`
	MaxDailyGust     unit.Speed       `mapstructure:"maxdailygust"`
	Model            string           `mapstructure:"model"`
	StationType      string           `mapstructure:"stationtype"`
	SolarRadiation   xunit.Irradiance `mapstructure:"solarradiation"`
	TempOutdoor      unit.Temperature `mapstructure:"tempf"`
	TempIndoor       unit.Temperature `mapstructure:"tempinf"`
	UltravioletIndex int              `mapstructure:"uv"`
}

func (e ecowitt) ToObservation() weather.Observation {
	return weather.Observation{
		Timestamp:          e.Timestamp,
		BarometricAbsHpa:   e.BarometricAbs.Hectopascals(),
		BarometricRelHpa:   e.BarometricRel.Hectopascals(),
		HourlyRainMm:       e.HourlyRain.Millimeters(),
		DailyRainMm:        e.DailyRain.Millimeters(),
		WeeklyRainMm:       e.WeeklyRain.Millimeters(),
		MonthlyRainMm:      e.MonthlyRain.Millimeters(),
		TotalRainMm:        e.TotalRain.Millimeters(),
		EventRainMm:        e.EventRain.Millimeters(),
		RainRatePerHourMm:  e.RainRatePerHour.Millimeters(),
		HumidityOutdoorPct: float64(e.HumidityOutdoor) / 100.0,
		HumidityIndoorPct:  float64(e.HumidityIndoor) / 100.0,
		WindDirDeg:         e.WindDir.Degrees(),
		WindGustKph:        e.WindGust.KilometersPerHour(),
		WindSpeedKph:       e.WindSpeed.KilometersPerHour(),
		MaxDailyGustKph:    e.MaxDailyGust.KilometersPerHour(),
		Model:              e.Model,
		StationType:        e.StationType,
		SolarRadiationWm2:  e.SolarRadiation.WattsPerSquareMetre(),
		TempOutdoorC:       e.TempOutdoor.Celsius(),
		TempIndoorC:        e.TempIndoor.Celsius(),
		UltravioletIndex:   e.UltravioletIndex,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var status int
	defer func() {
		httpRequests.WithLabelValues(http.StatusText(status)).Inc()
	}()

	err := req.ParseForm()
	if err != nil {
		status = http.StatusBadRequest
		http.Error(w, fmt.Sprintf("invalid form data: %s", err), status)
		return
	}

	keys := make([]string, 0, len(req.PostForm))
	for k := range req.PostForm {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	d := make(map[string]string, len(keys))
	for _, k := range keys {
		d[k] = req.PostForm.Get(k)
	}

	obj, err := h.decode(d)
	if err != nil {
		log.Printf("Error decoding ecowitt data: %s", err)
		status = http.StatusInternalServerError
		http.Error(w, fmt.Sprintf("Error decoding ecowitt data: %s", err), status)
		return
	}

	obs := obj.ToObservation()
	_, err = h.Store.Write(&obs)
	if err != nil {
		log.Printf("Error writing observation: %s", err)
		status = http.StatusInternalServerError
		http.Error(w, fmt.Sprintf("Error writing observation: %s", err), status)
		return
	}

	status = http.StatusOK
	w.WriteHeader(status)
}

func (h *Handler) decode(d map[string]string) (ecowitt, error) {
	var obj ecowitt
	cfg := mapstructure.DecoderConfig{
		DecodeHook:       decoderHookFn,
		WeaklyTypedInput: true,
		Result:           &obj,
	}
	dec, err := mapstructure.NewDecoder(&cfg)
	if err != nil {
		// an error here indicates a program error
		panic(err)
	}

	err = dec.Decode(d)
	return obj, err
}
