package meteorology

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardinalDirection(t *testing.T) {
	tests := []struct {
		degrees float64
		want    string
	}{
		{0, "N"},
		{15, "NNE"},
		{35, "NE"},
		{65, "ENE"},

		{90, "E"},
		{105, "ESE"},
		{125, "SE"},
		{155, "SSE"},

		{180, "S"},
		{195, "SSW"},
		{215, "SW"},
		{245, "WSW"},

		{270, "W"},
		{285, "WNW"},
		{305, "NW"},
		{335, "NNW"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%03.0f_%s", tt.degrees, tt.want), func(t *testing.T) {
			assert.Equal(t, tt.want, CardinalDirection(tt.degrees))
		})
	}
}
