package matcha

import (
	"fmt"
	gomaps "maps"
	"sync"
)

// Subscriber represents a callback function that is triggered
// whenever the Atom's value changes.
//
// Each subscriber is identified by a unique ID so that it can
// be individually added or removed.
type Subscriber[T any] struct {
	id string
	cb func(value T)
}

// Atom is a reactive state container that holds a value of type T
// and notifies its subscribers whenever the value changes.
//
// Concurrency:
//   - All reads/writes to the Atom's value and subscriber list are protected by an internal RWMutex.
//   - Callbacks are invoked outside of the lock to avoid deadlocks and blocking other operations.
//   - Each update increments a `version` counter; this version is embedded in subscriber IDs to
//     ensure stale subscribers from previous UI builds are never removed prematurely.
//
// Subscribers are stored in a map keyed by subscriber ID.
type Atom[T any] struct {
	ID          string
	Value       T
	subscribers map[string]*Subscriber[T]
	mu          sync.RWMutex
	version     int // Incremented on every update; used to isolate subscriptions between renders.
}

// subscribe adds a new subscriber to the Atom.
//
// If a subscriber with the same ID already exists, this is a no-op.
// The subscriber map is lazily initialized.
// Thread-safe.
func (a *Atom[T]) subscribe(subscriber *Subscriber[T]) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.subscribers == nil {
		a.subscribers = make(map[string]*Subscriber[T])
	}
	if _, ok := a.subscribers[subscriber.id]; ok {
		return
	}
	a.subscribers[subscriber.id] = subscriber
}

// unsubscribe removes multiple subscribers from the Atom.
//
// The passed `subscribers` map is typically a snapshot of the subscriber list
// taken during `update`. Only subscribers with IDs matching the snapshot
// (and thus the same version) are removed.
//
// Thread-safe.
func (a *Atom[T]) unsubscribe(subscribers map[string]*Subscriber[T]) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.subscribers == nil {
		return
	}
	for _, subscriber := range subscribers {
		delete(a.subscribers, subscriber.id)
	}
}

// value returns the current value stored in the Atom.
//
// This method acquires a read lock to ensure safe concurrent access,
// making it safe to call from multiple goroutines without risking
// data races.
//
// The returned value is a snapshot at the moment of the call;
// it is not guaranteed to reflect future changes.
func (a *Atom[T]) value() T {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Value
}

// update atomically updates the Atom's value and notifies all subscribers.
//
// The provided `updateFn` receives the current value and must return the new value.
//
// Implementation details:
//  1. The subscriber map is copied inside the lock. This prevents:
//     - Holding the lock while running callbacks (avoids deadlocks).
//     - Panics from concurrent map iteration/writes.
//     - Inconsistent subscriber sets if subscribe/unsubscribe happens mid-notification.
//  2. The `version` is incremented with every update. Subscriber IDs embed this version so
//     that only subscribers from the same render cycle are cleaned up.
//  3. Callbacks are invoked outside the lock and wrapped with `recover` so that a panic in
//     one subscriber does not prevent others from running or crash the system.
//  4. After notifying, stale subscribers from the old version are batch-removed.
//
// Thread-safe.
func (a *Atom[T]) update(updateFn func(oldValue T) T) {
	var newValue T
	subscribers := make(map[string]*Subscriber[T], len(a.subscribers))

	a.mu.Lock()
	oldValue := a.Value
	newValue = updateFn(oldValue)
	a.Value = newValue
	gomaps.Copy(subscribers, a.subscribers)
	a.version++
	a.mu.Unlock()

	for _, subscriber := range subscribers {
		func(sub *Subscriber[T]) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("subscriber %s panicked: %v\n", sub.id, r)
				}
			}()
			sub.cb(newValue)
		}(subscriber)
	}

	a.unsubscribe(subscribers)
}

// UseAtomState binds an Atom's value to a component's context and returns:
//  1. The current value
//  2. A setter function for updating the value
//
// When called:
//   - A subscription is registered for this context, with the current Atom version
//     embedded in the subscriber ID.
//   - On state change, the UI is re-rendered and stale subscriptions from
//     previous versions are automatically cleaned up.
//
// Thread-safe.
func UseAtomState[T any](ctx *Context, atom *Atom[T]) (T, func(func(T) T)) {
	id := fmt.Sprintf("%s/%d", ctx.id, atom.version)
	atom.subscribe(&Subscriber[T]{id: id, cb: func(value T) {
		ctx.channels.render <- struct{}{}
	}})
	atom.mu.RLock()
	defer atom.mu.RUnlock()
	return atom.Value, atom.update
}

// UseAtomValue binds an Atom's value to a component's context and returns
// only the current value.
//
// The component will re-render when the Atom's value changes.
// Subscription uses the current Atom version to ensure proper cleanup.
func UseAtomValue[T any](ctx *Context, atom *Atom[T]) T {
	id := fmt.Sprintf("%s/%d", ctx.id, atom.version)
	atom.subscribe(&Subscriber[T]{id: id, cb: func(value T) {
		ctx.channels.render <- struct{}{}
	}})
	atom.mu.RLock()
	defer atom.mu.RUnlock()
	return atom.Value
}

// UseAtomSetter binds only a setter function for updating the Atom's value.
//
// The setter applies the same update/notify/unsubscribe cycle as UseAtomState,
// but does not return the current value.
func UseAtomSetter[T any](ctx *Context, atom *Atom[T]) func(func(T) T) {
	id := fmt.Sprintf("%s/%d", ctx.id, atom.version)
	atom.subscribe(&Subscriber[T]{id: id, cb: func(value T) {
		ctx.channels.render <- struct{}{}
	}})
	return atom.update
}
