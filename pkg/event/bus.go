package event

import (
	"fmt"
	"reflect"
	"sync"
)

// Topic is a type for specifying names which may be used to publish or subscribe to events.
type Topic struct{ string }

// T constructs a topic using the specified name.
func T(name string) Topic { return Topic{name} }

func (t Topic) String() string { return t.string }

type Bus struct {
	mu          sync.RWMutex
	subscribers map[Topic][]*subscriber
}

// New returns new Bus.
func New() *Bus {
	return &Bus{
		subscribers: make(map[Topic][]*subscriber),
	}
}

// Subscribe subscribes fn to the specified topic.
// Returns error if sub is not a function.
func (bus *Bus) Subscribe(topic Topic, sub interface{}) error {
	fnv := reflect.ValueOf(sub)
	if fnv.Kind() != reflect.Func {
		return fmt.Errorf("%s is not of type reflect.Func", fnv.Type().Name())
	}

	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.subscribers[topic] = append(bus.subscribers[topic], newSubscriberValue(fnv))

	return nil
}

// HasSubscribers returns true if topic has any subscribers.
func (bus *Bus) HasSubscribers(topic Topic) bool {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	return len(bus.subscribers[topic]) > 0
}

// Unsubscribe removes sub from the list of subscribers to the specified topic.
func (bus *Bus) Unsubscribe(topic Topic, sub interface{}) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if len(bus.subscribers[topic]) > 0 {
		bus.removeSubscriber(topic, bus.findSubscriberIndex(topic, reflect.ValueOf(sub)))
	}
}

// Publish executes callback defined for a topic. Any additional arguments will be transferred to the callback.
func (bus *Bus) Publish(topic Topic, args ...interface{}) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	if handlers := bus.subscribers[topic]; len(handlers) > 0 {
		// copy, in case handlers are modified during callbacks
		handlers = append([]*subscriber(nil), handlers...)
		for _, handler := range handlers {
			handler.Call(args)
		}
	}
}

func (bus *Bus) removeSubscriber(topic Topic, idx int) {
	l := len(bus.subscribers[topic])

	if idx < 0 || idx >= l {
		return
	}

	copy(bus.subscribers[topic][idx:], bus.subscribers[topic][idx+1:])
	bus.subscribers[topic][l-1] = nil // or the zero value of T
	bus.subscribers[topic] = bus.subscribers[topic][:l-1]
}

func (bus *Bus) findSubscriberIndex(topic Topic, fn reflect.Value) int {
	for idx, handler := range bus.subscribers[topic] {
		if hfn := handler.fn; hfn.Type() == fn.Type() && hfn.Pointer() == fn.Pointer() {
			return idx
		}
	}
	return -1
}

type subscriber struct {
	fn   reflect.Value
	args []reflect.Value
}

func newSubscriberValue(fn reflect.Value) *subscriber {
	h := &subscriber{fn: fn, args: make([]reflect.Value, fn.Type().NumIn())}
	h.zeroArgs()
	return h
}

func (h *subscriber) Call(args []interface{}) {
	h.buildArgs(args)
	h.fn.Call(h.args)
	h.zeroArgs()
}

func (h *subscriber) buildArgs(args []interface{}) {
	for i, v := range args {
		if v != nil {
			h.args[i] = reflect.ValueOf(v)
		}
	}
}

func (h *subscriber) zeroArgs() {
	typ := h.fn.Type()
	for i := range h.args {
		h.args[i] = reflect.Zero(typ.In(i))
	}
}
