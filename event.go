package matcha

import (
	"maps"
	"sync"

	"github.com/gdamore/tcell/v2"
)

// eventManager manages event handlers for all components in the application.
// It stores a mapping from component IDs (string) to their event handling functions.
// Access is synchronized with a mutex to allow concurrent registration and dispatch.
type eventManager struct {
	handlers map[string]func(tcell.Event) bool
	mu       sync.Mutex
}

// newEventManager creates and returns a new, empty eventManager.
func newEventManager() *eventManager {
	return &eventManager{
		handlers: make(map[string]func(tcell.Event) bool),
	}
}

func dispatch(app *App) {
	var tree *node
	for {
		select {
		case t := <-app.channels.tree:
			tree = t
		default:
		}
		select {
		case event := <-app.channels.event:
			if tree == nil {
				continue
			}
			var startNode *node

			switch e := event.(type) {
			case *tcell.EventKey:
				startNode = getNodeWithFocusOrRoot(app, tree)
			case *tcell.EventMouse:
				x, y := e.Position()
				startNode = findDeepestNodeAtPosition(tree, x, y)

			default:
				startNode = tree
			}

			handlers := make(map[string]func(tcell.Event) bool)

			app.managers.event.mu.Lock()
			maps.Copy(handlers, app.managers.event.handlers)
			app.managers.event.mu.Unlock()
			// Bubble up from the starting node
			for n := startNode; n != nil; n = n.parent {
				if handler, ok := handlers[n.id]; !ok {
					continue
				} else if handled := handler(event); handled {
					app.channels.render <- struct{}{}
					break
				}
			}
		}
	}
}

func findDeepestNodeAtPosition(root *node, x, y int) *node {
	var found *node

	var visit func(*node)
	visit = func(n *node) {
		if pointInBounds(x, y, n.box) {
			found = n
			for _, child := range n.children {
				visit(child)
			}
		}
	}

	visit(root)
	return found
}

func pointInBounds(x, y int, bounds *box) bool {
	return x >= bounds.x && x < bounds.x+bounds.width &&
		y >= bounds.y && y < bounds.y+bounds.height
}

func getNodeWithFocusOrRoot(app *App, tree *node) *node {
	app.managers.focus.mu.Lock()
	defer app.managers.focus.mu.Unlock()

	focus := app.managers.focus.focused
	if focus == "" {
		return tree
	}

	if node := findNodeByID(tree, focus); node == nil {
		return tree
	} else {
		return node
	}
}

// findNodeByID searches the tree recursively for a node with the given componentID.
// Returns nil if no match is found.
// WARNING: This method is not thread-safe and should be used with appropriate locks
func findNodeByID(n *node, id componentID) *node {
	if n.id == string(id) {
		return n
	}
	for _, child := range n.children {
		if found := findNodeByID(child, id); found != nil {
			return found
		}
	}
	return nil
}

// UseEvent registers an event handler for the component associated with this Context.
//
// The `handler` function should return true if the event is handled and should not bubble
// further up the tree, or false if it should continue bubbling to parent components.
//
// Handlers are keyed by the component's ID (`ctx.id`) and are stored in the global event manager.
// Thread-safe.
func UseEvent(ctx *Context, handler func(event tcell.Event) bool) {
	manager := ctx.managers.event
	manager.mu.Lock()
	defer manager.mu.Unlock()

	manager.handlers[ctx.id] = handler
}
