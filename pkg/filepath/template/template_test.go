package template

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemplate_Execute(t *testing.T) {
	tm := Must(New("file").Parse(`webcam_latest_{{ strftime "%H%M" .Now }}.jpg`))
	ts := time.Date(2004, 4, 9, 12, 13, 14, 15, time.UTC)
	var buf bytes.Buffer
	err := tm.Execute(&buf, map[string]interface{}{
		"Now": ts,
	})
	assert.NoError(t, err)
	assert.Equal(t, "webcam_latest_1213.jpg", buf.String())
}
