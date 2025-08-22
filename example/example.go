package main

import (
	"fmt"

	"github.com/cchirag/matcha"
	"github.com/charmbracelet/lipgloss"
	"github.com/gdamore/tcell/v2"
)

var (
	buttonStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#44624a"))

	columnStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#FFFFFF"))

	rowStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#ffffff")).
			Padding(1)
	textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFFFF"))
)

type app struct{}

func (a *app) Render(ctx *matcha.Context) matcha.Component {
	count, setCount := matcha.UseState(ctx, 0)

	return matcha.Row([]matcha.Component{
		matcha.Text(fmt.Sprintf("Count: %d", count), textStyle),
		matcha.Column([]matcha.Component{
			matcha.Button("Increment", func(event *tcell.EventMouse) bool {
				if event.Buttons() == tcell.Button1 {
					setCount(func(i int) int { return i + 1 })
					return true
				}
				return false
			}, buttonStyle),
			matcha.Button("Decrement", func(event *tcell.EventMouse) bool {
				if event.Buttons() == tcell.Button1 {
					setCount(func(i int) int { return i - 1 })
					return true
				}
				return false
			}, buttonStyle),
		}, 2, columnStyle),
	}, 2, rowStyle)
}

func main() {
	app := matcha.NewApp(&app{})

	if err := app.Render(); err != nil {
		fmt.Println(err.Error())
	}
}
