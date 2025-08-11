package matcha

import "sync"

// focusableID represents a per-element unique ID within a component.
// componentID represents the globally unique and stable identifier for a component
// across rerenders. These are separated to clarify intent and avoid accidental mixing.
type (
	focusableID string
	componentID string
)

// focusManager manages global focus state for all components.
//
// It tracks:
//   - The currently focused component's ID (`focused`).
//   - A mapping from focusable element IDs (`focusableID`) to the component that owns them (`registered`).
//   - A reverse mapping from component IDs (`componentID`) to a representative focusable ID (`inverse`)
//     for quick lookup in the opposite direction.
//
// All access is synchronized with a mutex for concurrent safety.
type focusManager struct {
	focused    componentID
	registered map[focusableID]componentID
	inverse    map[componentID]focusableID
	mu         sync.Mutex
}

// newFocusManager creates and returns a new, empty focusManager.
func newFocusManager() *focusManager {
	return &focusManager{
		focused:    "",
		registered: make(map[focusableID]componentID),
		inverse:    make(map[componentID]focusableID),
	}
}

// clean removes all focusable elements and inverse mappings
// that do not belong to the currently focused component.
//
// Thread-safe.
func (f *focusManager) clean() {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Clean registered elements
	for key, value := range f.registered {
		if value != f.focused {
			delete(f.registered, key)
		}
	}

	// Clean inverse mappings
	for compID := range f.inverse {
		if compID != f.focused {
			delete(f.inverse, compID)
		}
	}
}

// register links a focusableID to a componentID in both maps.
func (f *focusManager) register(fid focusableID, cid componentID) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.registered[fid] = cid
	f.inverse[cid] = fid
}

// unregister removes a focusableID and its corresponding inverse entry.
func (f *focusManager) unregister(fid focusableID) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if cid, ok := f.registered[fid]; ok {
		delete(f.registered, fid)
		delete(f.inverse, cid)
	}
}

// UseFocus registers a focusable element for the current component and
// returns focus helpers.
//
// Parameters:
//   - ctx: The component's rendering context, which holds a globally stable
//     componentID for focus tracking.
//   - id:  A per-element focusable ID unique within the component.
//
// Returns:
//   - isFocused:    true if the current component is focused.
//   - setIsFocused: function to set focus to the component that owns the given focusableID.
//   - blur:         function to clear focus from all components.
//
// Behavior:
//   - Only one component can have focus at a time.
//   - The mapping from focusable IDs to component IDs enables fast lookup
//     during event bubbling.
//   - This function does not trigger a rerender on focus changes — that should
//     be handled by the caller if needed.
//
// Example:
//
//	isFocused, setFocus, blur := UseFocus(ctx, "input1")
//	if isFocused {
//	    // render input as highlighted
//	}
//	setFocus("input2") // move focus to another registered element
//	blur()             // clear focus completely
func UseFocus(ctx *Context, id string) (isFocused bool, setIsFocused func(id string), blur func()) {
	manager := ctx.managers.focus
	manager.mu.Lock()
	defer manager.mu.Unlock()

	isFocused = manager.focused == componentID(ctx.id)

	// Register the element's focusable ID to its owning component.
	manager.registered[focusableID(id)] = componentID(ctx.id)
	manager.inverse[componentID(ctx.id)] = focusableID(id)

	setIsFocused = func(newID string) {
		manager.mu.Lock()
		defer manager.mu.Unlock()
		if id, ok := manager.registered[focusableID(newID)]; ok && id != manager.focused {
			manager.focused = id
			ctx.channels.render <- struct{}{}
		}
	}

	blur = func() {
		manager.mu.Lock()
		defer manager.mu.Unlock()
		if manager.focused != "" {
			manager.focused = ""
			ctx.channels.render <- struct{}{}
		}
	}

	return isFocused, setIsFocused, blur
}
