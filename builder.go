package matcha

import (
	"fmt"
)

type node struct {
	id        string
	component Component
	children  []*node
}

func (a *App) build(component Component, id string) *node {
	node := &node{
		id: id,
	}

	ctx := &Context{
		id:        id,
		hookIndex: 0,
		store:     a.store,
		throttler: a.throttler,
	}

	switch c := component.(type) {

	case *text:
		node.component = c.Render(ctx)

	case *column:
		for i, child := range c.children {
			childID := fmt.Sprintf("%s/%d", id, i)
			childNode := a.build(child, childID)
			node.children = append(node.children, childNode)
		}
		node.component = c.Render(ctx)

	case *row:
		node.component = c
		for i, child := range c.children {
			childID := fmt.Sprintf("%s/%d", id, i)
			childNode := a.build(child, childID)
			node.children = append(node.children, childNode)
		}
		node.component = c.Render(ctx)
	default:
		rendered := c.Render(ctx)
		node.component = rendered
		childNode := a.build(rendered, id+"/0")
		node.children = append(node.children, childNode)
	}

	return node
}
