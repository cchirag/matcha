package matcha

import (
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/uniseg"
)

type cell struct {
	primary rune
	comb    []rune
	style   tcell.Style
}

func (c *cell) string() string {
	var builder strings.Builder
	primary := ' '
	if c.primary != 0 {
		primary = c.primary
	}

	builder.WriteRune(primary)
	for _, r := range c.comb {
		builder.WriteRune(r)
	}
	return builder.String()
}

type box struct {
	// Absolute value
	x, y int
	grid [][]cell
}

func pack(node *node, x, y int) *box {
	if node == nil {
		return nil
	}
	b := &box{
		x:    x,
		y:    y,
		grid: make([][]cell, 0),
	}
	switch c := node.component.(type) {
	case *text:
		text := packText(c, b.x, b.y)
		copyInto(text, b)
	case *row:
		row := packRow(node, c, b.x, b.y)
		copyInto(row, b)
	case *column:
		column := packColumn(node, c, b.x, b.y)
		copyInto(column, b)
	default:
		log.Println("Hitting default")
		d := pack(node.children[0], x, y)
		copyInto(d, b)
	}

	// log.Println("-------------------")
	// log.Println("Type: ", reflect.TypeOf(node.component))
	// log.Println("Type: ", reflect.TypeOf(node.parent.component))
	node.box = b
	return b
}

func packRow(node *node, component *row, x, y int) *box {
	b := &box{
		x:    x,
		y:    y,
		grid: make([][]cell, 0),
	}

	lines := 0
	columns := 0
	boxes := make([]*box, len(node.children))
	offset := 0
	for i, child := range node.children {
		childBox := pack(child, x, y+offset)
		cw, ch := childBox.dimensions()
		offset += ch + component.gap
		lines += ch
		if i != len(node.children)-1 {
			lines += component.gap
		}
		if cw > columns {
			columns = cw
		}
		boxes[i] = childBox
	}

	b.expandTo(columns, lines)

	str := b.string()

	b.fromLipgloss(str, component.style)

	for _, cb := range boxes {
		borderLeft := 0
		if component.style.GetBorderLeft() {
			borderLeft = 1
		}
		cb.x = cb.x + component.style.GetMarginLeft() + borderLeft + component.style.GetPaddingLeft()
		borderTop := 0
		if component.style.GetBorderTop() {
			borderTop = 1
		}
		cb.y = cb.y + component.style.GetMarginTop() + borderTop + component.style.GetPaddingTop()
		copyInto(cb, b)
	}
	return b
}

func packColumn(node *node, component *column, x, y int) *box {
	b := &box{
		x:    x,
		y:    y,
		grid: make([][]cell, 0),
	}

	lines := 0
	columns := 0
	boxes := make([]*box, len(node.children))
	offset := 0
	for i, child := range node.children {
		childBox := pack(child, x+offset, y)

		cw, ch := childBox.dimensions()
		offset += cw + component.gap
		columns += cw
		if i != len(node.children)-1 {
			columns += component.gap
		}
		if ch > lines {
			lines = ch
		}
		boxes[i] = childBox
	}

	b.expandTo(columns, lines)

	str := b.string()

	b.fromLipgloss(str, lipgloss.NewStyle().Inherit(component.style))

	for _, cb := range boxes {
		borderLeft := 0
		if component.style.GetBorderLeft() {
			borderLeft = 1
		}
		cb.x = cb.x + component.style.GetMarginLeft() + borderLeft + component.style.GetPaddingLeft()
		borderTop := 0
		if component.style.GetBorderTop() {
			borderTop = 1
		}
		cb.y = cb.y + component.style.GetMarginTop() + borderTop + component.style.GetPaddingTop()
		copyInto(cb, b)
	}
	return b
}

func packText(t *text, x, y int) *box {
	b := &box{
		x:    x,
		y:    y,
		grid: make([][]cell, 0),
	}
	b.fromLipgloss(t.content, t.style)
	return b
}

func copyInto(from, to *box) {
	fw, fh := from.dimensions()

	relativeX := from.x - to.x
	relativeY := from.y - to.y
	to.expandTo(fw+relativeX, fh+relativeY)

	for row := range fh {
		for col := range fw {
			c := from.grid[row][col]
			to.grid[row+relativeY][col+relativeX] = c
		}
	}

	// to.ansiString = from.ansiString
}

// func (b *box) copyInto(child *box) {
// 	cw, ch := child.dimensions()
//
// 	b.expandTo(cw+child.relative.x, ch+child.relative.y)
//
// 	for row := range ch {
// 		for col := range cw {
// 			char := child.grid[row][col]
// 			b.grid[row+child.relative.y][col+child.relative.x] = char
// 		}
// 	}
// }

//
// func (b *box) copyIntoPreserveContent(child *box) {
// 	cw, ch := child.dimensions()
// 	log.Printf("DEBUG: child dimensions (%d,%d), child relative (%d,%d)\n", cw, ch, child.relative.x, child.relative.y)
//
// 	targetWidth := cw + child.relative.x
// 	targetHeight := ch + child.relative.y
// 	log.Printf("DEBUG: expanding to (%d,%d)\n", targetWidth, targetHeight)
//
// 	b.expandTo(targetWidth, targetHeight)
//
// 	parentW, parentH := b.dimensions()
// 	log.Printf("DEBUG: parent after expand (%d,%d)\n", parentW, parentH)
//
// 	for row := 0; row < ch; row++ {
// 		for col := 0; col < cw; col++ {
// 			parentRow := row + child.relative.y
// 			parentCol := col + child.relative.x
//
// 			log.Printf("DEBUG: accessing parent[%d][%d], parent size (%d,%d)\n", parentRow, parentCol, parentH, parentW)
//
// 			if parentRow >= parentH || parentCol >= parentW {
// 				panic(fmt.Sprintf("Index out of bounds: trying to access [%d][%d] in grid of size (%d,%d)", parentRow, parentCol, parentH, parentW))
// 			}
//
// 			childChar := child.grid[row][col]
// 			if childChar.primary != ' ' || childChar.style != tcell.StyleDefault {
// 				b.grid[parentRow][parentCol] = childChar
// 			}
// 		}
// 	}
// }

func (b *box) dimensions() (width, height int) {
	if len(b.grid) == 0 {
		return 0, 0
	}
	return len(b.grid[0]), len(b.grid)
}

func (b *box) expandTo(width, height int) {
	_, currentH := b.dimensions()

	if currentH < height {
		rowsToAdd := height - currentH
		for range rowsToAdd {
			b.grid = append(b.grid, make([]cell, width))
		}
	}

	for row := range b.grid {
		if len(b.grid[row]) < width {
			colsToAdd := width - len(b.grid[row])
			b.grid[row] = append(b.grid[row], make([]cell, colsToAdd)...)
		}
	}
}

func (b *box) inBound(x, y int) bool {
	if b == nil {
		return false
	}
	w, h := b.dimensions()
	return x >= b.x && x < b.x+w &&
		y >= b.y && y < b.y+h
}

func (b *box) string() string {
	var builder strings.Builder

	for row := range b.grid {
		for column := range b.grid[row] {
			builder.WriteString(b.grid[row][column].string())
		}
		if row != len(b.grid)-1 {
			builder.WriteRune('\n')
		}
	}

	return builder.String()
}

func (b *box) ansi() string {
	return ""
}

// fromLipgloss converts a string with a Lip Gloss style into a *box structure
// containing a 2D grid of characters (with styling) that can be rendered
// in a terminal using tcell.
//
// The function:
//  1. Renders the content with the given Lip Gloss style.
//  2. Strips ANSI escape codes from the rendered string.
//  3. Measures the height and width of the content in terms of grapheme clusters.
//  4. Initializes a grid of `character` cells with blank spaces.
//  5. Fills the grid with the graphemes from the content, assigning either
//     border styles or content styles to each cell based on position.
//
// Border detection uses the style's margin and border settings.
// Content styles are extracted via extractContentStyle().
func (b *box) fromLipgloss(content string, style lipgloss.Style) *box {
	rendered := style.Render(content)
	// b.ansiString = rendered
	rendered = ansi.Strip(rendered)

	lines := strings.Split(rendered, "\n")
	height := len(lines)
	width := 0
	for _, line := range lines {
		count := uniseg.GraphemeClusterCount(line)
		if count > width {
			width = count
		}
	}

	b.grid = make([][]cell, height)
	for i := range b.grid {
		b.grid[i] = make([]cell, width)
		for j := range b.grid[i] {
			b.grid[i][j] = cell{
				primary: ' ',
				style:   tcell.StyleDefault,
			}
		}
	}

	contentStyle := extractContentStyle(style)
	borderStyles := map[int]tcell.Style{
		1: getBorderStyle(&style, 1),
		2: getBorderStyle(&style, 2),
		3: getBorderStyle(&style, 3),
		4: getBorderStyle(&style, 4),
	}

	row := 0
	for _, line := range lines {
		if row >= height {
			break
		}
		column := 0
		gr := uniseg.NewGraphemes(line)
		count := uniseg.GraphemeClusterCount(line)
		for gr.Next() {
			if column >= width {
				break
			}

			runes := gr.Runes()
			cluster := string(runes)
			if cluster == "\n" {
				break
			}

			var cellStyle tcell.Style
			if ok, side := isBorder(row, column, height, count, style); ok {
				cellStyle = borderStyles[side]
			} else {
				cellStyle = contentStyle
			}

			primaryRune := runes[0]
			var comb []rune
			if len(runes) == 0 {
				primaryRune = ' '
			}
			if len(runes) > 1 {
				comb = runes[1:]
			}
			b.grid[row][column] = cell{
				primary: primaryRune,
				comb:    comb,
				style:   cellStyle,
			}
			column++
		}
		row++
	}

	return b
}

// extractContentStyle converts a Lip Gloss style into a tcell.Style that
// contains text attributes (bold, italic, underline, etc.) and foreground/
// background colors.
func extractContentStyle(style lipgloss.Style) tcell.Style {
	s := tcell.StyleDefault.Bold(style.GetBold()).
		Dim(style.GetFaint()).
		Italic(style.GetItalic()).
		Reverse(style.GetReverse()).
		StrikeThrough(style.GetStrikethrough()).
		Underline(style.GetUnderline()).
		Background(lipglossColorToTcell(style.GetBackground())).
		Foreground(lipglossColorToTcell(style.GetForeground()))
	return s
}

// getBorderStyle returns a tcell.Style for a given border side (1=top,
// 2=right, 3=bottom, 4=left) based on the corresponding Lip Gloss border
// foreground/background colors.
func getBorderStyle(style *lipgloss.Style, side int) tcell.Style {
	switch side {
	case 1:
		return tcell.StyleDefault.
			Background(lipglossColorToTcell(style.GetBorderTopBackground())).
			Foreground(lipglossColorToTcell(style.GetBorderTopForeground()))
	case 2:
		return tcell.StyleDefault.
			Background(lipglossColorToTcell(style.GetBorderRightBackground())).
			Foreground(lipglossColorToTcell(style.GetBorderRightForeground()))
	case 3:
		return tcell.StyleDefault.
			Background(lipglossColorToTcell(style.GetBorderBottomBackground())).
			Foreground(lipglossColorToTcell(style.GetBorderBottomForeground()))
	case 4:
		return tcell.StyleDefault.
			Background(lipglossColorToTcell(style.GetBorderLeftBackground())).
			Foreground(lipglossColorToTcell(style.GetBorderLeftForeground()))
	}
	return tcell.StyleDefault
}

// isBorder checks whether a given cell coordinate (row, column) lies on a
// styled border, returning true and the border side number if so.
// Border side codes: 1=top, 2=right, 3=bottom, 4=left.
//
// The logic respects Lip Gloss margins, so the border's position is offset
// inward by any top/left/right/bottom margins.
func isBorder(row, column, height, width int, style lipgloss.Style) (bool, int) {
	mt, mr, mb, ml := style.GetMarginTop(), style.GetMarginRight(), style.GetMarginBottom(), style.GetMarginLeft()

	if style.GetBorderTop() &&
		row == mt &&
		column >= ml &&
		column < width-mr {
		return true, 1
	}

	if style.GetBorderRight() &&
		column == width-mr-1 &&
		row > mt &&
		row < height-mb-1 {
		return true, 2
	}

	if style.GetBorderBottom() &&
		row == height-mb-1 &&
		column >= ml &&
		column < width-mr {
		return true, 3
	}

	if style.GetBorderLeft() &&
		column == ml &&
		row > mt &&
		row < height-mb-1 {
		return true, 4
	}
	return false, 0
}

// lipglossColorToTcell converts a Lip Gloss TerminalColor into a tcell.Color.
// If the alpha channel is zero, it returns ColorDefault (meaning "no color").
// Otherwise, it constructs a 24-bit RGB tcell color.
func lipglossColorToTcell(color lipgloss.TerminalColor) tcell.Color {
	r, g, b, a := color.RGBA()
	if a == 0 {
		return tcell.ColorDefault
	}
	return tcell.NewRGBColor(int32(r/257), int32(g/257), int32(b/257))
}
