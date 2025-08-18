package matcha

import (
	"fmt"
	"reflect"
	"sync"
)

type hooksManager struct {
	states   map[string]any
	statesMu sync.Mutex

	refs   map[string]any
	refsMu sync.Mutex

	effects string

	memos   map[string][]any
	memosMu sync.Mutex

	callbacks string
}

func newHooksManager() *hooksManager {
	return &hooksManager{
		states: make(map[string]any),
		memos:  make(map[string][]any),
	}
}

// UseState is a state management hook inspired by React's useState.
//
// It associates a persistent value with the current component and provides
// a setter function to update it. State is identified by the position of
// the hook call within a component's Render function, so the order of hooks
// must remain stable across renders.
//
// Example usage:
//
//	count, setCount := UseState(ctx, 0)
//
//	Button("Increment", func() {
//	    setCount(func(prev int) int { return prev + 1 })
//	})
//
// Parameters:
//
//	ctx    - the component's render context
//	initial - the initial value to use if no prior state exists
//
// Returns:
//
//	value  - the current state value
//	setter - a function that updates state. It accepts a function that
//	         transforms the previous state into a new state.
func UseState[T any](ctx *Context, initial T) (T, func(func(T) T)) {
	key := fmt.Sprintf("%s/%d", ctx.id, ctx.hookIdx)
	ctx.hookIdx++

	manager := ctx.managers.hooks

	manager.statesMu.Lock()
	var value T
	if v, ok := manager.states[key]; !ok {
		manager.states[key] = initial
		value = initial
	} else {
		value = v.(T)
	}
	manager.statesMu.Unlock()

	setter := func(fn func(prev T) T) {
		manager.statesMu.Lock()
		old := manager.states[key]
		manager.states[key] = fn(old.(T))
		manager.statesMu.Unlock()
		ctx.channels.render <- struct{}{}
	}

	return value, setter
}

func UseEffect(ctx *Context, effect func() func(), deps []any) {
	// On mount, run the effect and store the return
	// On every rerender, check if the dependency has changed, if yes, we rerun the effect and store the return
	// On unmount, we run the returned function (check for nil)
}

// Ref is a container type that holds a mutable value.
//
// The value is stored in the Current field and can be freely read or modified.
// Unlike state, updating Current does not cause the component to re-render.
//
// Refs are useful for persisting values across renders or for accessing values
// inside closures without capturing stale data.
type Ref[T any] struct {
	Current T
}

// UseRef creates and returns a new reference object that persists across
// component re-renders.
//
// The initial value is only applied on the first render. Subsequent renders
// return the same Ref. Updating ref.Current will not trigger a re-render.
//
// Typical use cases include:
//
//   - Storing values that should persist without causing re-renders
//   - Holding onto resources like timers, connections, or UI handles
//   - Accessing mutable values inside goroutines or event handlers
//
// Example:
//
//	func Render(ctx *matcha.Context) matcha.Component {
//	    countRef := matcha.UseRef(ctx, 0)
//
//	    // Increment without causing a re-render
//	    countRef.Current++
//
//	    return matcha.Text(fmt.Sprintf("Count is %d", countRef.Current))
//	}
func UseRef[T any](ctx *Context, initial T) *Ref[T] {
	key := fmt.Sprintf("%s/%d", ctx.id, ctx.hookIdx)
	ctx.hookIdx++
	manager := ctx.managers.hooks
	ref := new(Ref[T])
	manager.refsMu.Lock()
	if v, ok := manager.refs[key]; !ok {
		ref.Current = initial
		manager.refs[key] = ref
	} else {
		ref = v.(*Ref[T])
	}
	manager.refsMu.Unlock()
	return ref
}

// UseMemo memoizes the result of a computation between renders.
// It accepts a function `fn` that computes a value of type `T`,
// and a list of dependencies `deps`.
//
// On the first call, UseMemo executes `fn`, stores the result along
// with the dependencies, and returns the computed value. On subsequent
// calls, if the dependencies are deeply equal to the previously stored
// dependencies, the cached value is returned without re-executing `fn`.
// If the dependencies differ, `fn` is executed again, and the new value
// is cached.
//
// Usage example:
//
//	count := UseMemo(ctx, func() int {
//		return expensiveCalculation()
//	}, []any{dep1, dep2})
func UseMemo[T any](ctx *Context, fn func() T, deps []any) T {
	key := fmt.Sprintf("%s/%d", ctx.id, ctx.hookIdx)
	ctx.hookIdx++

	manager := ctx.managers.hooks

	// Note: To keep the internal storage simple, this implementation
	// assumes that the first element in the stored slice is always
	// the computed value, followed by the dependencies in order.
	manager.memosMu.Lock()
	defer manager.memosMu.Unlock()

	if v, ok := manager.memos[key]; !ok {
		manager.memos[key] = make([]any, len(deps)+1)
		value := fn()
		manager.memos[key][0] = value
		manager.memos[key] = append(manager.memos[key], deps...)

		return value
	} else {
		if reflect.DeepEqual(deps, v[1:]) {
			return v[0].(T)
		} else {
			return fn()
		}
	}
}

func UseCallback[T any](ctx *Context, fn func() T, deps []any) func() T {
	return fn
}
