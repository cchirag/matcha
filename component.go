package matcha

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gdamore/tcell/v2"
)

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
	gap      int
}

func (c *column) Render(ctx *Context) Component {
	return c
}

func Column(children []Component, gap int, style lipgloss.Style) Component {
	return &column{children: children, gap: gap, style: style}
}

// Row

type row struct {
	children []Component
	style    lipgloss.Style
	gap      int
}

func (r *row) Render(ctx *Context) Component {
	return r
}

func Row(children []Component, gap int, style lipgloss.Style) Component {
	return &row{children: children, gap: gap, style: style}
}

type button struct {
	text    string
	style   lipgloss.Style
	handler func(event *tcell.EventMouse) bool
}

func (b *button) Render(ctx *Context) Component {
	UseEvent(ctx, func(event tcell.Event) bool {
		switch e := event.(type) {
		case *tcell.EventMouse:
			return b.handler(e)
		}
		return false
	})

	return Text(b.text, b.style)
}

func Button(text string, handler func(event *tcell.EventMouse) bool, style lipgloss.Style) Component {
	return &button{text: text, handler: handler, style: style}
}
