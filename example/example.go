package main

import (
	"fmt"
	"strconv"

	"github.com/cchirag/matcha"
)

type header struct {
	content string
}

func (h *header) Render(ctx *matcha.Context) matcha.Component {
	count, setCount := matcha.UseState(ctx, 0)
	setCount(func(i int) int { return i + 1 })
	text := strconv.Itoa(count)
	return matcha.Text(text, matcha.NewStyle())
}

func main() {
	app := matcha.NewApp(&header{})

	if err := app.Render(); err != nil {
		fmt.Println(err.Error())
	}
}
