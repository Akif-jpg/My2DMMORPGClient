package eventsystem

import (
	"reflect"
	"sync"
)

// EventManager is a generic event bus that stores callbacks for each
// event type identified by a string key. The generic type parameter T
// represents the data type carried by the event.
type EventManager[T any] struct {
	mu     sync.RWMutex
	events map[string][]EventHandler[T]
}

// NewEventManager creates a new EventManager with an initialized
// internal map. It is safe to use concurrently.
func NewEventManager[T any]() *EventManager[T] {
	return &EventManager[T]{
		events: make(map[string][]EventHandler[T]),
	}
}

// Register registers an event handler for the specified event type.
// Handlers are invoked in the order they were registered.
func (m *EventManager[T]) Register(eventType string, handler EventHandler[T]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events[eventType] = append(m.events[eventType], handler)
}

// Unregister removes a previously registered handler for the specified
// event type. It returns true if a handler was removed, or false if
// no matching handler was found. Because Go function values are not
// directly comparable, the comparison uses reflect.ValueOf(handler).Pointer()
// which retrieves the underlying function pointer. This works for
// function literals and named functions but may not be reliable for
// certain closure patterns.
func (m *EventManager[T]) Unregister(eventType string, handler EventHandler[T]) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	handlers := m.events[eventType]
	for i, h := range handlers {
		if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
			// Remove the handler by slicing.
			m.events[eventType] = append(handlers[:i], handlers[i+1:]...)
			return true
		}
	}
	return false
}

// Emit broadcasts an event to all handlers registered for its type.
// Handlers are executed sequentially in a goroutine-safety context.
func (m *EventManager[T]) Emit(event Event[T]) {
	m.mu.RLock()
	handlers := m.events[event.EventType]
	m.mu.RUnlock()

	for _, h := range handlers {
		// Call each handler synchronously. If asynchronous behavior is
		// desired, the handler itself can spawn goroutines.
		h(event)
	}
}
