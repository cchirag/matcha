package matcha

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/uniseg"
)

// toBox converts a string with a Lip Gloss style into a *box structure
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
func toBox(content string, style lipgloss.Style) *box {
	box := new(box)

	rendered := style.Render(content)

	rendered = ansi.Strip(rendered)

	lines := strings.Split(rendered, "\n")
	box.height = len(lines)
	box.width = 0
	for _, line := range lines {
		width := uniseg.GraphemeClusterCount(line)
		if width > box.width {
			box.width = width
		}
	}

	box.grid = make([][]character, box.height)
	for i := range box.grid {
		box.grid[i] = make([]character, box.width)
		for j := range box.grid[i] {
			box.grid[i][j] = character{
				ch:    ' ',
				style: tcell.StyleDefault,
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
		if row >= box.height {
			break
		}
		column := 0
		gr := uniseg.NewGraphemes(line)
		count := uniseg.GraphemeClusterCount(line)
		for gr.Next() {
			if column >= box.width {
				break
			}

			runes := gr.Runes()
			cluster := string(runes)
			if cluster == "\n" {
				break
			}

			var cellStyle tcell.Style
			if ok, side := isBorder(row, column, box.height, count, style); ok {
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
			box.grid[row][column] = character{
				ch:    primaryRune,
				comb:  comb,
				style: cellStyle,
			}
			column++
		}
		row++
	}

	return box
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
