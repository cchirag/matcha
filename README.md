---

# Matcha

Matcha is an experimental terminal UI (TUI) framework written in Go.
It’s inspired by React and Flutter — focusing on **components**, **state**, and **composition** instead of the traditional Elm-style architecture.

⚠️ This is an early project. Expect breaking changes, incomplete features, and lots of iteration. Contributions, ideas, and feedback are welcome!

---

## Features (current + planned)

* ✅ Component-based API (similar to React/Flutter)
* ✅ Local component state with `UseState`
* ✅ Built-in components: `Row`, `Column`, `Text`, `Button`
* ✅ Style support with [Lip Gloss](https://github.com/charmbracelet/lipgloss)
* ✅ Event handling with [tcell](https://github.com/gdamore/tcell)
* 🚧 Layout improvements (flex, sizing)
* 🚧 More interactive components (input fields, lists, etc.)
* 🚧 Continuous rendering / animations

---

## Example

Below is a simple counter app built with Matcha:

```go
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

	textStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#FFFFFF"))
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
```

Run it:

```bash
go run example/example.go
```

---

## Why Matcha?

Most Go TUI frameworks (like [Bubble Tea](https://github.com/charmbracelet/bubbletea)) follow an Elm-style model/update/view loop.
Matcha explores a different approach:

* Components are **composable** (like React)
* Each component manages its own **state**
* You can pass **props** to children and render trees declaratively

This makes Matcha feel familiar to developers coming from UI frameworks like Flutter or React, but in the terminal.

---

## Getting Started

```bash
go get github.com/cchirag/matcha
```

Then start building by defining a root component with a `Render` method.

---

## Contributing

This project is experimental and under heavy iteration. If you have ideas, bug reports, or feature suggestions:

* Open an issue
* Start a discussion
* Or submit a PR

Let’s shape this together.

---

## License

[MIT](LICENSE)

---
