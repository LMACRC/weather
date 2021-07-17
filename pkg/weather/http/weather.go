package http

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/lmacrc/weather/pkg/mapconv"
	"github.com/lmacrc/weather/pkg/weather/model"
	"github.com/lmacrc/weather/pkg/xunit"
	"github.com/martinlindhe/unit"
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

	httpForwardRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "weather",
		Subsystem: "http",
		Name:      "forward_requests_total",
		Help:      "The total number of forwarded weather requests",
	}, []string{"status"})
)

type ObservationWriter interface {
	WriteObservation(o model.Observation) (*model.Observation, error)
}

type Handler struct {
	log       *zap.Logger
	path      string
	forwardTo string
	store     ObservationWriter
}

// InitViper sets any default values for vp.
func InitViper(vp *viper.Viper) {
	vp.SetDefault("http.path", "/weather")
}

func New(log *zap.Logger, vp *viper.Viper, store ObservationWriter) (*Handler, error) {
	log = log.With(zap.String("service", "http_handler"))

	var forwardTo string
	if vp.IsSet("http.dev.forward_to") {
		forwardTo = vp.GetString("http.dev.forward_to")
		if _, err := url.Parse(forwardTo); err != nil {
			log.Error("Unable to forward HTTP requests, invalid url.", zap.Error(err))
			return nil, fmt.Errorf("forward_url: %w", err)
		} else {
			log.Info("Forwarding HTTP requests.", zap.String("forward_to", forwardTo))
		}
	}

	return &Handler{
		log:       log,
		path:      vp.GetString("http.path"),
		forwardTo: forwardTo,
		store:     store,
	}, nil
}

func (h *Handler) Handle(mux *http.ServeMux) {
	mux.Handle(h.path, h)
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

func (e ecowitt) ToObservation() model.Observation {
	return model.Observation{
		Timestamp:        e.Timestamp,
		BarometricAbs:    e.BarometricAbs,
		BarometricRel:    e.BarometricRel,
		HourlyRain:       e.HourlyRain,
		DailyRain:        e.DailyRain,
		WeeklyRain:       e.WeeklyRain,
		MonthlyRain:      e.MonthlyRain,
		TotalRain:        e.TotalRain,
		EventRain:        e.EventRain,
		RainRatePerHour:  e.RainRatePerHour,
		HumidityOutdoor:  e.HumidityOutdoor,
		HumidityIndoor:   e.HumidityIndoor,
		WindDir:          e.WindDir,
		WindGust:         e.WindGust,
		WindSpeed:        e.WindSpeed,
		MaxDailyGust:     e.MaxDailyGust,
		Model:            e.Model,
		StationType:      e.StationType,
		SolarRadiation:   e.SolarRadiation,
		TempOutdoor:      e.TempOutdoor,
		TempIndoor:       e.TempIndoor,
		UltravioletIndex: e.UltravioletIndex,
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

	if h.forwardTo != "" {
		freq, err := http.NewRequest(http.MethodPost, h.forwardTo, strings.NewReader(req.PostForm.Encode()))
		if err != nil {
			h.log.Warn("Failed to create forward HTTP request", zap.Error(err))
		}
		freq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		fres, err := defaultClient.Do(freq)
		if err != nil {
			httpForwardRequests.WithLabelValues("error").Inc()
			h.log.Warn("Failed to forward HTTP request", zap.Error(err))
		} else if fres != nil {
			httpForwardRequests.WithLabelValues("ok").Inc()
			_ = fres.Body.Close()
		}
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
		h.log.Error("Error decoding ecowitt data.", zap.Error(err))
		status = http.StatusInternalServerError
		http.Error(w, fmt.Sprintf("Error decoding ecowitt data: %s", err), status)
		return
	}

	obs := obj.ToObservation()
	_, err = h.store.WriteObservation(obs)
	if err != nil {
		h.log.Error("Error writing observation", zap.Error(err))
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
