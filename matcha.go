package matcha

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gdamore/tcell/v2"
	"github.com/muesli/termenv"
)

type channels struct {
	event  chan tcell.Event
	tree   chan *node
	quit   chan struct{}
	render chan struct{}
}

type managers struct {
	focus *focusManager
	event *eventManager
}

type App struct {
	root     Component
	screen   tcell.Screen
	channels *channels
	managers *managers
}

func NewApp(component Component) *App {
	return &App{
		root: component,
		channels: &channels{
			event:  make(chan tcell.Event, 1),
			tree:   make(chan *node, 1),
			quit:   make(chan struct{}, 1),
			render: make(chan struct{}, 1),
		},
		managers: &managers{
			focus: newFocusManager(),
			event: newEventManager(),
		},
	}
}

func (a *App) Render() error {
	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	defer screen.Fini()

	a.screen = screen

	if err := screen.Init(); err != nil {
		return err
	}

	go screen.ChannelEvents(a.channels.event, a.channels.quit)

	go dispatch(a)

	a.channels.render <- struct{}{}

	go build(a)

	<-a.channels.quit

	return nil
}

func (a *App) newContext(id string) *Context {
	return &Context{
		id:       id,
		channels: a.channels,
		managers: a.managers,
	}
}
