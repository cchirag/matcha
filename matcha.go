package matcha

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
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
			quit:   make(chan struct{}),
			render: make(chan struct{}, 10),
		},
		managers: &managers{
			focus: newFocusManager(),
			event: newEventManager(),
		},
	}
}

func (a *App) Render() error {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Set log output to file
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

	// Handle messages
	go dispatch(a)

	go build(a, 10)

	a.channels.render <- struct{}{}

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
