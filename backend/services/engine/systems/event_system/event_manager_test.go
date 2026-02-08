package eventsystem

import (
	"sync"
	"testing"
)

func TestEventManager_RegisterAndEmit(t *testing.T) {
	em := NewEventManager[int]()
	var mu sync.Mutex
	calls := []int{}

	handler1 := func(e Event[int]) {
		mu.Lock()
		calls = append(calls, e.Data)
		mu.Unlock()
	}
	handler2 := func(e Event[int]) {
		mu.Lock()
		calls = append(calls, e.Data*10)
		mu.Unlock()
	}

	em.Register("test", handler1)
	em.Register("test", handler2)

	em.Emit(Event[int]{EventType: "test", Data: 5})

	mu.Lock()
	defer mu.Unlock()
	if len(calls) != 2 {
		t.Fatalf("expected 2 handlers to be called, got %d", len(calls))
	}
	if calls[0] != 5 || calls[1] != 50 {
		t.Fatalf("unexpected handler results: %v", calls)
	}
}

func TestEventManager_Unregister(t *testing.T) {
	em := NewEventManager[int]()
	var mu sync.Mutex
	calls := []int{}

	handler1 := func(e Event[int]) {
		mu.Lock()
		calls = append(calls, 1)
		mu.Unlock()
	}
	handler2 := func(e Event[int]) {
		mu.Lock()
		calls = append(calls, 2)
		mu.Unlock()
	}

	em.Register("test", handler1)
	em.Register("test", handler2)

	removed := em.Unregister("test", handler1)
	if !removed {
		t.Fatalf("handler1 was not removed")
	}

	em.Emit(Event[int]{EventType: "test", Data: 0})

	mu.Lock()
	defer mu.Unlock()
	if len(calls) != 1 || calls[0] != 2 {
		t.Fatalf("unexpected calls after unregister: %v", calls)
	}
}

func TestEventManager_UnregisterNotFound(t *testing.T) {
	em := NewEventManager[int]()

	handler1 := func(e Event[int]) {}

	removed := em.Unregister("nonexistent", handler1)
	if removed {
		t.Fatalf("Unregister returned true for a handler that was never registered")
	}
}
