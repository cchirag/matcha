package main

import (
	"fmt"

	"github.com/cchirag/matcha"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	app := matcha.NewApp(matcha.Text("Hello world",
		lipgloss.NewStyle().
			Height(30).
			Width(30).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Top).
			Foreground(lipgloss.Color("#C3D7EE")).
			Background(lipgloss.Color("#002B49")).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.AdaptiveColor{
				Light: "#3C3C3C",
				Dark:  "#04B575",
			}),
	))

	if err := app.Render(); err != nil {
		fmt.Println(err.Error())
	}
}
