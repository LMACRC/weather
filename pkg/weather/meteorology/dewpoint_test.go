package meteorology

import (
	"testing"

	"github.com/martinlindhe/unit"
	"github.com/stretchr/testify/assert"
)

func TestDewPoint(t *testing.T) {
	got := DewPoint(unit.FromCelsius(11.1), 84)
	assert.InEpsilon(t, 8.4, got.Celsius(), 0.1)
}
