package service

import (
	"testing"

	"github.com/lmacrc/weather/pkg/mapconv"
	"github.com/lmacrc/weather/pkg/xunit"
	"github.com/martinlindhe/unit"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

var testData = map[string]string{
	"PASSKEY":        "6018F8D638BE61DF2E79DCF23DBACB79",
	"baromabsin":     "29.283",
	"baromrelin":     "30.033",
	"dailyrainin":    "0.000",
	"dateutc":        "2021-07-01 01:43:22",
	"eventrainin":    "0.000",
	"freq":           "433M",
	"hourlyrainin":   "0.000",
	"humidity":       "74",
	"humidityin":     "52",
	"maxdailygust":   "21.7",
	"model":          "WS2900C_V2.01.13",
	"monthlyrainin":  "0.000",
	"rainratein":     "0.000",
	"solarradiation": "309.27",
	"stationtype":    "EasyWeatherV1.5.9",
	"tempf":          "56.7",
	"tempinf":        "72.9",
	"totalrainin":    "0.469",
	"uv":             "3",
	"weeklyrainin":   "0.000",
	"wh65batt":       "0",
	"winddir":        "344",
	"windgustmph":    "6.9",
	"windspeedmph":   "5.4",
}

func TestIt(t *testing.T) {
	var obj ecowitt
	cfg := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc("2006-01-02 15:04:05"),
			mapconv.StringToLengthHookFunc(unit.Inch),
			mapconv.StringToPressureHookFunc(unit.InchOfMercury),
			mapconv.StringToSpeedHookFunc(unit.MilesPerHour),
			mapconv.StringToAngleFunc(unit.Degree),
			mapconv.StringToIrradianceFunc(xunit.WattPerSquareMetre),
			mapconv.StringToTemperatureHookFunc(unit.FromFahrenheit),
		),
		WeaklyTypedInput: true,
		Result:           &obj,
	}

	dec, err := mapstructure.NewDecoder(&cfg)
	assert.NoError(t, err)

	err = dec.Decode(testData)
	obs := obj.ToObservation()
	t.Logf("%v", obs)
	assert.NoError(t, err)
}
