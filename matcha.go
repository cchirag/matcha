package matcha

import (
	"context"
	"time"

	"github.com/cchirag/matcha/renderer"
	"github.com/gdamore/tcell/v2"
)

type App struct {
	root      Component
	store     *store
	throttler *throttler
}

func NewApp(component Component) *App {
	return &App{root: component, store: newStore(), throttler: newThrottler(30)}
}

func (a *App) Render() error {
	screen, err := renderer.NewScreen(context.Background())
	if err != nil {
		return err
	}
	defer screen.Fini()

	if err := screen.Initialize(); err != nil {
		return err
	}

	go a.handleRender(screen)

	go a.handleEvent()

	time.Sleep(time.Second * 10)

	return nil
}

func (a *App) handleRender(screen *renderer.Screen) {
	for range a.throttler.channel() {
		tree := a.build(a.root, "root")
		box := a.parcel(tree, 0, 0)
		for i, row := range box.grid {
			for j, char := range row {
				screen.WriteContent(j, i, char.ch, tcell.StyleDefault)
			}
		}
		screen.Show()
	}
}

func (a *App) handleEvent() {}
