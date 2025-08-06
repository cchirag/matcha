package layout

import (
	"fmt"

	"github.com/cchirag/matcha/components"
	"github.com/cchirag/matcha/core"
)

type Node struct {
	ID        string
	Component core.Component
	Children  []*Node
}

func Build(component core.Component, id string) *Node {
	node := new(Node)
	switch component := component.(type) {
	case *components.Text:
		node.ID = id
		node.Component = component
	case *components.Column:
		node.ID = id
		node.Component = component
		for index, child := range component.Children {
			childID := fmt.Sprintf("%s/%d", id, index)
			node.Children = append(node.Children, Build(child, childID))
		}
	case *components.Row:
		node.ID = id
		node.Component = component
		for index, child := range component.Children {
			childID := fmt.Sprintf("%s/%d", id, index)
			node.Children = append(node.Children, Build(child, childID))
		}
	}
	return node
}
