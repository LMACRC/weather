package db

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/lmacrc/weather/pkg/weather/service/camera"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func newGetImageCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get-image",
		Short: "Capture image",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := zap.NewNop()
			cam, err := camera.New(log, viper.GetViper(), nil)
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			name, err := cam.CaptureImage(ctx, time.Now())
			if err != nil {
				return err
			}

			fmt.Printf("Captured image: %s\n", name)

			return nil
		},
	}
}
