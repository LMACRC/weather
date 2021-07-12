package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/lmacrc/weather/pkg/weather/meteorology"
	"github.com/lmacrc/weather/pkg/weather/reporting"
	"github.com/lmacrc/weather/pkg/xunit"
	"github.com/martinlindhe/unit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func newGetStatsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get-stats",
		Short: "Get latest statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := reporting.New(zap.NewNop(), viper.GetViper(), db)
			if err != nil {
				return err
			}

			res := r.Generate(time.Now())
			s := structs.New(res)
			fields := s.Fields()
			ftoa := func(v float64) string {
				return strconv.FormatFloat(v, 'f', 1, 64)
			}

			for _, f := range fields {
				var s string
				switch v := f.Value().(type) {
				case unit.Temperature:
					s = ftoa(v.Celsius()) + " °C"
				case unit.Speed:
					s = ftoa(v.KilometersPerHour()) + " km/h"
				case unit.Pressure:
					s = ftoa(v.Hectopascals()) + " hPa"
				case unit.Length:
					s = ftoa(v.Millimeters()) + " mm"
				case xunit.Irradiance:
					s = ftoa(v.WattsPerSquareMetre()) + " w/m²"
				case unit.Angle:
					s = ftoa(v.Degrees()) + "°"
				case string:
					s = v
				case bool:
					if v {
						s = "✓"
					} else {
						s = "𐄂"
					}
				case int:
					s = strconv.Itoa(v)
				case time.Time:
					s = v.Format(time.RFC850)
				case meteorology.Direction:
					s = string(v)
				case float64:
					s = ftoa(v)
				default:
					s = "<no conversion>"
				}

				fmt.Printf("%-20s: %s\n", f.Name(), s)
			}

			return nil
		},
	}
}
