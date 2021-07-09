package db

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var getLastCommand = &cobra.Command{
	Use:   "get-last",
	Short: "Get last reading",
	RunE: func(cmd *cobra.Command, args []string) error {
		res := db.LastObservation(time.Now())
		if res != nil {
			fmt.Println(res)
		}
		return nil
	},
}
