package main

import (
	"fmt"
	"strconv"

	"github.com/cchirag/matcha"
	"github.com/gdamore/tcell/v2"
)

var countAtom = &matcha.Atom[int]{
	ID:    "count",
	Value: 10,
}

var nameAtom = &matcha.Atom[string]{
	ID:    "name",
	Value: "name",
}

type header struct {
	content string
}

func (h *header) Render(ctx *matcha.Context) matcha.Component {
	count, setCount := matcha.UseAtomState(ctx, countAtom)
	name, setName := matcha.UseAtomState(ctx, nameAtom)

	isFocused, setFocused, _ := matcha.UseFocus(ctx, "header")

	matcha.UseEvent(ctx, func(event tcell.Event) bool {
		switch event := event.(type) {
		case *tcell.EventKey:
			if event.Rune() == '+' {
				setCount(func(i int) int { return i + 1 })
				return true
			} else if event.Rune() == 'q' {
				ctx.Quit()
				return true
			} else if event.Key() == tcell.KeyCR && isFocused {
				setFocused("hero")
				return true
			} else {
				if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
					setName(func(newName string) string {
						if len(newName) == 0 {
							return ""
						}
						return newName[:len(newName)-1]
					})
					return true
				} else {
					setName(func(newName string) string { return fmt.Sprintf("%s%c", newName, event.Rune()) })
					return true
				}
			}
		}
		return false
	})

	text := strconv.Itoa(count)
	return matcha.Column([]matcha.Component{
		matcha.Text(text, matcha.NewStyle()),
		matcha.Text(name, matcha.NewStyle()),
	}, matcha.NewStyle())
}

type hero struct {
	content string
}

func (h *hero) Render(ctx *matcha.Context) matcha.Component {
	count, setCount := matcha.UseAtomState(ctx, countAtom)
	name, _ := matcha.UseAtomState(ctx, nameAtom)

	isFocused, setFocused, _ := matcha.UseFocus(ctx, "hero")
	matcha.UseEvent(ctx, func(event tcell.Event) bool {
		switch event := event.(type) {
		case *tcell.EventKey:
			if event.Rune() == '-' {
				setCount(func(i int) int { return i - 1 })
				return true
			} else if event.Rune() == 'q' {
				ctx.Quit()
				return true
			} else if event.Key() == tcell.KeyCR && isFocused {
				setFocused("header")
				return true
			}
			// } else if message.Rune() == '-' {
			// 	setCount(func(i int) int { return i - 1 })
			// }
		}
		return false
	})

	text := strconv.Itoa(count)
	return matcha.Column([]matcha.Component{
		matcha.Text(text, matcha.NewStyle()),
		matcha.Text(name, matcha.NewStyle()),
	}, matcha.NewStyle())
}

func main() {
	app := matcha.NewApp(matcha.Column([]matcha.Component{
		&header{},
		&hero{},
	}, matcha.NewStyle()))

	if err := app.Render(); err != nil {
		fmt.Println(err.Error())
	}
}
