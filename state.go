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

func UseState[T any](ctx Context, initialValue T) (T, func(newValue T)) {
	store := ctx.store
	key := stateKey{componentID: ctx.id, hookIndex: ctx.hookIndex}
	ctx.hookIndex++
	raw, ok := store.entries[key]
	if !ok {
		store.entries[key] = initialValue
		raw = initialValue
	}

	value := raw.(T)

	setter := func(newValue T) {
		store.entries[key] = newValue
	}

	return value, setter
}
