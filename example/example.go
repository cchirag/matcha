package main

import (
	"fmt"

	"github.com/cchirag/matcha"
	"github.com/charmbracelet/lipgloss"
	"github.com/gdamore/tcell/v2"
)

type Banner struct{}

func (b *Banner) Render(ctx *matcha.Context) matcha.Component {
	count, setCount := matcha.UseState(ctx, 0)
	_, _, _ = matcha.UseFocus(ctx, "banner")
	matcha.UseEvent(ctx, func(event tcell.Event) bool {
		switch e := event.(type) {
		case *tcell.EventKey:
			switch e.Rune() {
			case '+':
				setCount(func(i int) int {
					return i + 1
				})
				return true
			case '-':
				setCount(func(i int) int {
					return i - 1
				})
				return true
			}
		case *tcell.EventMouse:
			switch e.Buttons() {
			case tcell.Button1:
				setCount(func(i int) int {
					return i + 1
				})
				return true
			case tcell.Button2:
				setCount(func(i int) int {
					return i - 1
				})
				return true
			}
		}
		return false
	})
	return matcha.Text(fmt.Sprintf("Count: %d", count),
		lipgloss.NewStyle().
			Height(30).
			Width(30).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Bottom).
			Foreground(lipgloss.Color("#C3D7EE")).
			Background(lipgloss.Color("#002B49")).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.AdaptiveColor{
				Light: "#3C3C3C",
				Dark:  "#04B575",
			}),
	)
}

type app struct{}

func (a *app) Render(ctx *matcha.Context) matcha.Component {
	count, setCount := matcha.UseState(ctx, 0)
	return matcha.Row([]matcha.Component{
		matcha.Text(fmt.Sprintf("Count: %d", count), lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFFFF"))),

		matcha.Column([]matcha.Component{
			matcha.Button("Increment", func(event *tcell.EventMouse) bool {
				if event.Buttons() == tcell.Button1 {
					setCount(func(i int) int { return i + 1 })
					return true
				}
				return false
			}, lipgloss.NewStyle().
				Padding(1, 2).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#44624a")),
			),

			matcha.Button("Decrement", func(event *tcell.EventMouse) bool {
				if event.Buttons() == tcell.Button1 {
					setCount(func(i int) int { return i - 1 })
					return true
				}
				return false
			}, lipgloss.NewStyle().
				Padding(1, 2).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#44624a"))),
		}, 2, lipgloss.NewStyle().
			Background(lipgloss.Color("#FFFFFF"))),
	}, 2,
		lipgloss.NewStyle().
			Background(lipgloss.Color("#ffffff")).
			Padding(1),
	)
}

func main() {
	app := matcha.NewApp(&app{})

	if err := app.Render(); err != nil {
		fmt.Println(err.Error())
	}
}
