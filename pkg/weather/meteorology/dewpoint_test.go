package meteorology

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDewPoint(t *testing.T) {
	got := DewPoint(11.1, 84)
	assert.InEpsilon(t, 8.4, got, 0.1)
}
