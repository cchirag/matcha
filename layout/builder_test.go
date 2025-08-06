package layout

import (
	"testing"

	"github.com/cchirag/matcha/components"
	"github.com/cchirag/matcha/core"
)

func TestBuilder(t *testing.T) {
	component := components.Column{Children: []core.Component{
		&components.Text{Content: "Hello"},
		&components.Text{Content: "World"},
	}}

	node := Build(&component, "root")

	t.Errorf("Tree: %+v", node.Children[0].Component)
}
