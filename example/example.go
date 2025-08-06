package main

import (
	"fmt"

	"github.com/cchirag/matcha"
)

type button struct {
	Text string
}

func Button(text string) *button {
	return &button{Text: text}
}

func (b *button) Render(ctx *matcha.Context) matcha.Component {
	return matcha.Text("Hello\nworld")
}

func main() {
	app := matcha.NewApp(Button("Hello"))

	if err := app.Render(); err != nil {
		fmt.Println(err.Error())
	}
}
