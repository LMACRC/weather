package cmd

import (
	"net/http"
	"strconv"

	"github.com/lmacrc/weather/pkg/weather/service"
	"github.com/lmacrc/weather/pkg/weather/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func newServer() *cobra.Command {
	flags := struct {
		Port int
	}{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the LMACRC weather server",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := store.New()
			if err != nil {
				return err
			}

			mux := http.NewServeMux()

			mux.Handle("/metrics", promhttp.Handler())

			wh := service.Handler{Path: "/weather", Store: s}
			wh.Handle(mux)

			return http.ListenAndServe(":"+strconv.Itoa(flags.Port), mux)
		},
	}

	cmd.Flags().IntVarP(&flags.Port, "port", "p", 9876, "Port to listen for requests")

	return cmd
}
