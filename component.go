package matcha

import "github.com/charmbracelet/lipgloss"

type Component interface {
	Render(ctx *Context) Component
}

type HasKey interface {
	Key() string
}

// Text
type text struct {
	key     string
	content string
	style   lipgloss.Style
}

func (t *text) Render(ctx *Context) Component {
	return t
}

func (t *text) Key() string {
	return t.key
}

func Text(content string, style lipgloss.Style, key ...string) Component {
	k := ""
	if len(key) > 0 {
		k = key[0]
	}
	return &text{content: content, style: style, key: k}
}

// Column
type column struct {
	children []Component
	style    lipgloss.Style
}

func (c *column) Render(ctx *Context) Component {
	return c
}

func Column(children []Component, style lipgloss.Style) Component {
	return &column{children: children, style: style}
}

// Row

type row struct {
	children []Component
	style    lipgloss.Style
}

func (r *row) Render(ctx *Context) Component {
	return r
}

func Row(children []Component, style lipgloss.Style) Component {
	return &row{children: children, style: style}
}
