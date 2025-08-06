package components

import "github.com/cchirag/matcha/core"

type Row struct {
	Children []core.Component
}

func (r *Row) Render(ctx *core.Ctx) core.Component {
	return r
}
