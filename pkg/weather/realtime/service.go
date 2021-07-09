// Package realtime is responsible for generating weather statistics
package realtime

import (
	"bytes"
	"io"
	"time"

	"github.com/lmacrc/weather/pkg/weather/reporting"
	"go.uber.org/zap"
)

// https://cumuluswiki.org/a/Realtime.txt#List_of_fields_in_the_file

type Service struct {
	Log      *zap.Logger
	Reporter *reporting.Reporter
}

func (s Service) String() string {
	return "Realtime"
}

func (s Service) Data(ts time.Time) (path string, r io.ReadCloser, err error) {
	s.Log.Info("Generating realtime.txt data.")

	stats := s.Reporter.Generate(ts)
	data, err := Statistics(*stats).MarshalText()
	if err != nil {
		return "", nil, err
	}

	return "realtime.txt", io.NopCloser(bytes.NewReader(data)), nil
}
