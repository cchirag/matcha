package components

import "github.com/cchirag/matcha/core"

type Column struct {
	Children []core.Component
}

func (c *Column) Render(ctx *core.Ctx) core.Component {
	return c
}
