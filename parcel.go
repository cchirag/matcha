package matcha

type character struct {
	ch    rune
	style *Style
}

type box struct {
	x, y, height, width int
	grid                [][]character
}

func (a *App) parcel(tree *node, x, y int) *box {
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
