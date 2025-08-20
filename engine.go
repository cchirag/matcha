package matcha

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type node struct {
	id        string
	component Component
	children  []*node
	parent    *node
	box       *box
}

func build(app *App) {
	buffer := 0
	ticker := time.NewTicker(time.Second / 24)
	for {
		select {
		case <-app.channels.render:
			buffer++
			if buffer >= 10 {
				tree := walk(app, app.root, "root", nil)
				app.channels.tree <- tree
				box := pack(tree, 0, 0)
				render(app.screen, box)
				buffer = 0
			}

		case <-ticker.C:
			if buffer != 0 {
				tree := walk(app, app.root, "root", nil)
				app.channels.tree <- tree
				box := pack(tree, 0, 0)
				render(app.screen, box)
				buffer = 0
			}
		case <-app.channels.quit:
			return
		}
	}
}

func walk(app *App, component Component, id string, parent *node) *node {
	nodeId := id
	if hk, ok := component.(HasKey); ok {
		keys := strings.Split(nodeId, "/")
		keys[len(keys)-1] = hk.Key()
		nodeId = strings.Join(keys, "/")
	}

	ctx := app.newContext(id)

	if component == nil {
		return nil
	}

	rendered := component.Render(ctx)
	if rendered == nil {
		return nil
	}

	node := &node{
		id:        nodeId,
		parent:    parent,
		component: rendered,
	}

	switch c := node.component.(type) {
	case *text:
		// node.component = c
	case *column:
		for i, child := range c.children {
			childID := fmt.Sprintf("%s/%d", id, i)
			if childNode := walk(app, child, childID, node); childNode != nil {
				node.children = append(node.children, childNode)
			}
		}

	case *row:
		node.component = c
		for i, child := range c.children {
			childID := fmt.Sprintf("%s/%d", id, i)
			if childNode := walk(app, child, childID, node); childNode != nil {
				node.children = append(node.children, childNode)
			}
		}
	default:
		node = walk(app, node.component, id, parent)
	}
	return node
}

func render(screen tcell.Screen, box *box) {
	for y, row := range box.grid {
		for x, column := range row {
			screen.SetContent(x, y, column.primary, column.comb, column.style)
		}
	}
	screen.Show()
}
