package matcha

type Component interface {
	Render(ctx *Context) Component
}

// Text
type text struct {
	content string
	style   *Style
}

func (t *text) Render(ctx *Context) Component {
	return t
}

func Text(content string, style *Style) Component {
	return &text{content: content, style: style}
}

// Column
type column struct {
	children []Component
	style    *Style
}

func (c *column) Render(ctx *Context) Component {
	return c
}

func Column(children []Component, style *Style) Component {
	return &column{children: children, style: style}
}

// Row

type row struct {
	children []Component
	style    *Style
}

func (r *row) Render(ctx *Context) Component {
	return r
}

func Row(children []Component, style *Style) Component {
	return &row{children: children, style: style}
}
