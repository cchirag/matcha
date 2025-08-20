package matcha

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

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
	hooks *hooksManager
}

type App struct {
	root       Component
	screen     tcell.Screen
	channels   *channels
	managers   *managers
	dimensions *dimensions
	error      error
}

func NewApp(component Component) *App {
	return &App{
		root:       component,
		dimensions: &dimensions{},
		channels: &channels{
			event:  make(chan tcell.Event, 64),
			tree:   make(chan *node, 1),
			quit:   make(chan struct{}, 1),
			render: make(chan struct{}, 1),
		},
		managers: &managers{
			focus: newFocusManager(),
			event: newEventManager(),
			hooks: newHooksManager(),
		},
	}
}

func (a *App) Render() error {
	logFile, panicFile := a.setupLog()

	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	defer safe(func() {
		screen.Fini()
		panicFile.Close()
		logFile.Close()
	}, panicFile)

	a.screen = screen
	if err := screen.Init(); err != nil {
		return err
	}
	screen.EnableMouse(tcell.MouseButtonEvents)

	go screen.ChannelEvents(a.channels.event, a.channels.quit)

	go dispatch(a)

	go build(a)

	w, h := screen.Size()
	if err := screen.PostEvent(tcell.NewEventResize(w, h)); err != nil {
		panic(err.Error())
	}

	<-a.channels.quit

	return a.error
}

func (a *App) newContext(id string) *Context {
	return &Context{
		id:         id,
		channels:   a.channels,
		managers:   a.managers,
		dimensions: a.dimensions,
	}
}

func (a *App) setupLog() (*os.File, *os.File) {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	// Set log output to file
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.OpenFile("panic.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Could not open panic log file:", err)
		return nil, nil
	}
	return logFile, f
}

func safe(fn func(), panicFile *os.File) {
	if r := recover(); r != nil {
		panicFile.Write((debug.Stack()))
	}

	fn()
}
