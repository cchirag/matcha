package matcha

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

type character struct {
	ch    rune
	comb  []rune
	style tcell.Style
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

func build(app *App) {
	buffer := 0
	ticker := time.NewTicker(time.Second / 24)
	for {
		select {
		case <-app.channels.render:
			buffer++
			if buffer >= 10 {
				tree := walk(app, app.root, "root", nil)
				app.channels.tree <- tree
				box := pack(tree, 0, 0)
				render(app.screen, box, app)
				buffer = 0
			}

		case <-ticker.C:
			if buffer != 0 {
				tree := walk(app, app.root, "root", nil)
				app.channels.tree <- tree
				box := pack(tree, 0, 0)
				render(app.screen, box, app)
				buffer = 0
			}
		case <-app.channels.quit:
			return
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
	var b *box
	switch c := tree.component.(type) {
	case *text:
		b = toBox(c.content, c.style)
	}

	tree.box = b

	return b
}

func (b *box) copyInto(child *box) {
	for row := 0; row < child.height; row++ {
		for col := 0; col < child.width; col++ {
			ch := child.grid[row][col]
			b.grid[row+child.y-b.y][col+child.x-b.x] = ch
		}
	}
}

func render(screen tcell.Screen, box *box, app *App) {
	for y, row := range box.grid {
		for x, column := range row {
			screen.SetContent(box.x+x, box.y+y, column.ch, column.comb, column.style)
		}
	}
	screen.Show()
	time.Sleep(time.Second * 2)
	close(app.channels.quit)
}
