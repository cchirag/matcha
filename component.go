package matcha

type Component interface {
	Render(ctx *Context) Component
}

// Text
type text struct {
	content string
}

func (t *text) Render(ctx *Context) Component {
	return t
}

func Text(content string) Component {
	return &text{content: content}
}

// Column
type column struct {
	children []Component
}

func (c *column) Render(ctx *Context) Component {
	return c
}

func Column(children []Component) Component {
	return &column{children: children}
}

// Row

type row struct {
	children []Component
}

func (r *row) Render(ctx *Context) Component {
	return r
}

func Row(children []Component) Component {
	return &row{children: children}
}
