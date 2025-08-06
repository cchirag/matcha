package components

import "github.com/cchirag/matcha/core"

type Text struct {
	Content string
}

func (t *Text) Render(ctx *core.Ctx) core.Component {
	return t
}
