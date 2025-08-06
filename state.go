package matcha

type stateKey struct {
	componentID string
	hookIndex   int
}

type store struct {
	entries map[stateKey]any
}

func newStore() *store {
	store := new(store)
	store.entries = make(map[stateKey]any)
	return store
}

func UseState[T any](ctx *Context, initialValue T) (T, func(func(T) T)) {
	store := ctx.store
	throttler := ctx.throttler
	key := stateKey{componentID: ctx.id, hookIndex: ctx.hookIndex}
	ctx.hookIndex++

	raw, ok := store.entries[key]
	if !ok {
		store.entries[key] = initialValue
		raw = initialValue
	}

	value := raw.(T)

	setter := func(updateFn func(T) T) {
		prev := store.entries[key].(T)
		newValue := updateFn(prev)
		store.entries[key] = newValue
		throttler.trigger()
	}

	return value, setter
}
