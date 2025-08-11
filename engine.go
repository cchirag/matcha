package matcha

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type character struct {
	ch    rune
	style *Style
}

type box struct {
	x, y, height, width int
	grid                [][]character
}

type node struct {
	id        string
	component Component
	children  []*node
	parent    *node
	box       *box
}

func (n *node) String() string {
	if n == nil {
		return "<nil>"
	}
	return n.stringWithPrefix("", true)
}

func (n *node) stringWithPrefix(prefix string, isLast bool) string {
	var result strings.Builder

	connector := "├── "
	if isLast {
		connector = "└── "
	}
	if prefix == "" {
		connector = ""
	}

	result.WriteString(fmt.Sprintf("%s%s%s (%d children)", prefix, connector, n.id, len(n.children)))

	childPrefix := prefix
	if prefix != "" {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	for i, child := range n.children {
		result.WriteString("\n")
		isLastChild := i == len(n.children)-1
		result.WriteString(child.stringWithPrefix(childPrefix, isLastChild))
	}

	return result.String()
}

func build(app *App, bufferSize int) {
	buffer := 0
	ticker := time.NewTicker(time.Second / 24)
	for {
		select {
		case <-app.channels.render:
			buffer++
			if buffer >= bufferSize {
				tree := walk(app, app.root, "root", nil)
				app.channels.tree <- tree
				box := pack(tree, 0, 0)
				render(app.screen, box)
				buffer = 0
			}

		case <-ticker.C:
			if buffer != 0 {
				tree := walk(app, app.root, "root", nil)
				app.channels.tree <- tree
				box := pack(tree, 0, 0)
				render(app.screen, box)
				buffer = 0
			}
		}
	}
}

func walk(app *App, component Component, id string, parent *node) *node {
	node := &node{
		id:     id,
		parent: parent,
	}
	ctx := app.newContext(id)
	switch c := component.(type) {

	case *text:
		node.component = c.Render(ctx)

	case *column:
		for i, child := range c.children {
			childID := fmt.Sprintf("%s/%d", id, i)
			childNode := walk(app, child, childID, node)
			node.children = append(node.children, childNode)
		}
		node.component = c.Render(ctx)

	case *row:
		node.component = c
		for i, child := range c.children {
			childID := fmt.Sprintf("%s/%d", id, i)
			childNode := walk(app, child, childID, node)
			node.children = append(node.children, childNode)
		}
		node.component = c.Render(ctx)
	default:
		rendered := c.Render(ctx)
		node.component = rendered
		childNode := walk(app, rendered, id+"/0", parent)
		node.children = append(node.children, childNode)
	}

	return node
}

func pack(tree *node, x, y int) *box {
	box := &box{
		x:    x,
		y:    y,
		grid: make([][]character, 50),
	}

	for i := range box.grid {
		box.grid[i] = make([]character, 100)
	}

	switch c := tree.component.(type) {

	case *text:
		row := 0
		column := 0
		maxWidth := 0

		for _, r := range c.content {
			if r == '\n' {
				column = 0
				row++
				continue
			}
			box.grid[row][column] = character{ch: r, style: c.style}
			column++
			if column > maxWidth {
				maxWidth = column
			}
		}

		box.width = maxWidth
		box.height = row + 1

	case *column:
		offsetY := 0
		maxWidth := 0

		for _, child := range tree.children {
			childBox := pack(child, x, y+offsetY)
			box.copyInto(childBox)
			offsetY += childBox.height
			if childBox.width > maxWidth {
				maxWidth = childBox.width
			}
		}

		box.width = maxWidth
		box.height = offsetY

	case *row:
		offsetX := 0
		maxHeight := 0

		for _, child := range tree.children {
			childBox := pack(child, x+offsetX, y)
			box.copyInto(childBox)
			offsetX += childBox.width
			if childBox.height > maxHeight {
				maxHeight = childBox.height
			}
		}

		box.width = offsetX
		box.height = maxHeight
	}

	tree.box = box
	return box
}

func (b *box) copyInto(child *box) {
	for row := 0; row < child.height; row++ {
		for col := 0; col < child.width; col++ {
			ch := child.grid[row][col]
			b.grid[row+child.y-b.y][col+child.x-b.x] = ch
		}
	}
}

func render(screen tcell.Screen, box *box) {
	for y, row := range box.grid {
		for x, column := range row {
			screen.SetContent(box.x+x, box.y+y, column.ch, nil, tcell.StyleDefault)
		}
	}
	screen.Show()
}
