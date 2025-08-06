package matcha

import (
	"context"
	"fmt"
	"time"

	"github.com/cchirag/matcha/renderer"
	"github.com/gdamore/tcell/v2"
)

type App struct {
	root  Component
	store *store
}

func NewApp(component Component) *App {
	return &App{root: component, store: newStore()}
}

func (a *App) Render() error {
	tree := a.build(a.root, "root")

	box := a.parcel(tree, 10, 10)

	screen, err := renderer.NewScreen(context.Background())
	if err != nil {
		return err
	}

	if err := screen.Initialize(); err != nil {
		return err
	}

	for i, row := range box.grid {
		for j, char := range row {
			screen.WriteContent(box.x+j, box.y+i, char.ch, tcell.StyleDefault)
		}
	}

	screen.Show()
	time.Sleep(time.Second * 10)
	return nil
}

type node struct {
	id        string
	component Component
	children  []*node
}

func (a *App) build(component Component, id string) *node {
	node := &node{
		id: id,
	}

	switch c := component.(type) {

	case *text:
		node.component = c

	case *column:
		node.component = c
		for i, child := range c.children {
			childID := fmt.Sprintf("%s/%d", id, i)
			childNode := a.build(child, childID)
			node.children = append(node.children, childNode)
		}

	case *row:
		node.component = c
		for i, child := range c.children {
			childID := fmt.Sprintf("%s/%d", id, i)
			childNode := a.build(child, childID)
			node.children = append(node.children, childNode)
		}

	default:
		ctx := &Context{
			id:        id,
			hookIndex: 0,
		}
		rendered := component.Render(ctx)
		node.component = rendered
		childNode := a.build(rendered, id+"/0")
		node.children = append(node.children, childNode)
	}

	return node
}

type character struct {
	ch rune
}

type box struct {
	x, y, height, width int
	grid                [][]character
}

func (a *App) parcel(tree *node, x, y int) *box {
	box := &box{
		x:    x,
		y:    y,
		grid: make([][]character, 50), // for demo, adjust later
	}

	for i := range box.grid {
		box.grid[i] = make([]character, 100) // terminal width
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
			box.grid[row][column] = character{ch: r}
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
			childBox := a.parcel(child, x, y+offsetY)
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
			childBox := a.parcel(child, x+offsetX, y)
			box.copyInto(childBox)
			offsetX += childBox.width
			if childBox.height > maxHeight {
				maxHeight = childBox.height
			}
		}

		box.width = offsetX
		box.height = maxHeight
	}

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
