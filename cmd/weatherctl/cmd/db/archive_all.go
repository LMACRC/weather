package db

import (
	"github.com/lmacrc/weather/pkg/weather/service/archive"
	"github.com/lmacrc/weather/pkg/weather/service/ftp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func newArchiveAllCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "archive-all",
		Short: "Generate CSV archives for all historical data",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log := zap.NewNop()

			vp := viper.GetViper()

			var ftpSvc *ftp.Service
			if vp.GetBool("ftp.enabled") {
				ftpSvc, err = ftp.New(log, db, vp)
				if err != nil {
					return err
				}
			}

			arSvc, _ := archive.New(zap.NewNop(), db, vp, st, ftpSvc)
			return arSvc.ArchiveAll()
		},
	}
}
