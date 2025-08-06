package state

import "github.com/cchirag/matcha/core"

type StateKey struct {
	ComponentID string
	HookIndex   int
}
type StateStore struct {
	Storage map[StateKey]any
}

func NewStateStore() *StateStore {
	return &StateStore{Storage: make(map[StateKey]any)}
}

func UseState[T any](ctx *core.Ctx, initialValue T) (T, func(T)) {
	return initialValue, func(t T) {}
}
