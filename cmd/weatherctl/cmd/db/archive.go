package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
	"github.com/lmacrc/weather/pkg/weather/service/archive"
	"github.com/lmacrc/weather/pkg/weather/service/ftp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func newArchiveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "archive YYYYMMDD,[YYYYMMDD]",
		Short: "Generate CSV archives for historical data on specific dates",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log := zap.NewNop()

			last := now.BeginningOfDay()

			vp := viper.GetViper()

			var ftpSvc *ftp.Service
			if vp.GetBool("ftp.enabled") {
				ftpSvc, err = ftp.New(log, db, vp)
				if err != nil {
					return err
				}
			}

			arSvc, _ := archive.New(zap.NewNop(), db, vp, st, ftpSvc)
			arSvc.ArchiveAll()

			var dates []time.Time
			for _, arg := range args {
				ts, err := time.Parse("20060102", arg)
				if err != nil {
					return fmt.Errorf("invalid date %q: must be YYYYMMDD", arg)
				} else if ts.Equal(last) || ts.After(last) {
					return fmt.Errorf("archive dates must be prior to today")
				}
				dates = append(dates, ts)
			}

			for _, ts := range dates {
				path, err := arSvc.Archive(ts)
				d := ts.Format("02 Jan 2006")
				if err != nil {
					fmt.Printf("Error archiving %s: %s\n", d, err)
				} else {
					fmt.Printf("Archived %s to %s\n", d, path)
				}
			}

			return nil
		},
	}
}
