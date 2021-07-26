package event_test

import (
	"testing"

	"github.com/lmacrc/weather/pkg/event"
	"github.com/stretchr/testify/assert"
)

func TestHasCallback(t *testing.T) {
	topic := event.T("topic")
	eb := event.New()
	_ = eb.Subscribe(topic, func() {})
	assert.True(t, eb.HasSubscribers(topic))
	assert.False(t, eb.HasSubscribers(event.T("foo")))
}

func TestSubscribe(t *testing.T) {
	eb := event.New()
	assert.NoError(t, eb.Subscribe(event.T("topic"), func() {}))
	assert.Error(t, eb.Subscribe(event.T("topic"), 0))
}

func TestManySubscribe(t *testing.T) {
	eb := event.New()
	topic := event.T("topic")
	count := 0
	fn := func() { count++ }
	_ = eb.Subscribe(topic, fn)
	_ = eb.Subscribe(topic, fn)
	eb.Publish(topic)

	assert.Equal(t, 2, count)
}

func TestUnsubscribe(t *testing.T) {
	eb := event.New()
	count := 0
	handler := func() { count++ }
	_ = eb.Subscribe(event.T("topic"), handler)
	eb.Publish(event.T("topic"))
	assert.Equal(t, 1, count)
	eb.Unsubscribe(event.T("topic"), handler)
	eb.Publish(event.T("topic"))
	assert.Equal(t, 1, count)
}

func TestPublish(t *testing.T) {
	eb := event.New()
	_ = eb.Subscribe(event.T("topic"), func(a int, err error) {
		assert.Equal(t, 10, a)
		assert.Nil(t, err)
	})
	eb.Publish(event.T("topic"), 10, nil)
}
