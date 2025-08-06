package renderer

import (
	"context"
	"errors"
	"sync"

	"github.com/cchirag/matcha/core"
	tcell "github.com/gdamore/tcell/v2"
)

type MouseFlags int

const (
	MouseButtonEvents = MouseFlags(1)
	MouseDragEvents   = MouseFlags(2)
	MouseMotionEvents = MouseFlags(4)
)

func (m MouseFlags) ToTCellFlags() tcell.MouseFlags {
	return tcell.MouseFlags(m)
}

type ScreenMode string

var (
	ScreenModeLight = ScreenMode("light")
	ScreenModeDark  = ScreenMode("dark")
)

var ErrNoMouseSupport = errors.New("no mouse support")

type Screen struct {
	screen       tcell.Screen
	mu           sync.Mutex
	mode         ScreenMode
	mouse        bool
	focus        bool
	paste        bool
	ctx          context.Context
	eventChannel chan tcell.Event
}

func NewScreen(ctx context.Context) (*Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	return &Screen{screen: screen, ctx: ctx, eventChannel: make(chan tcell.Event, 10)}, nil
}

func (s *Screen) TCellScreen() tcell.Screen {
	return s.screen
}

func (s *Screen) Initialize() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.screen.Init()
}

func (s *Screen) Fini() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.Fini()
}

func (s *Screen) SetScreenMode(mode ScreenMode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mode = mode
}

func (s *Screen) EnableFocus() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.EnableFocus()
	s.focus = true
}

func (s *Screen) DisableFocus() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.DisableFocus()
	s.focus = false
}

func (s *Screen) EnablePaste() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.EnablePaste()
	s.paste = true
}

func (s *Screen) DisablePaste() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.DisablePaste()
	s.paste = false
}

func (s *Screen) EnableMouse(flags ...MouseFlags) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.screen.HasMouse() {
		return ErrNoMouseSupport
	}

	f := make([]tcell.MouseFlags, len(flags))
	for index, flag := range flags {
		f[index] = flag.ToTCellFlags()
	}

	s.screen.EnableMouse(f...)
	s.mouse = true

	return nil
}

func (s *Screen) DisableMouse() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.DisableMouse()
	s.mouse = false
}

func (s *Screen) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.Fini()
}

func (s *Screen) Show() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.Show()
}

func (s *Screen) Sync() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.Sync()
}

func (s *Screen) Size() (width, height int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.screen.Size()
}

func (s *Screen) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.Clear()
}

func (s *Screen) MessageListener(messageCh chan<- core.Message, quitCh <-chan struct{}) {
	eventCh := make(chan tcell.Event, cap(messageCh))
	go func(screen tcell.Screen, quitCh <-chan struct{}) {
		// matcha <-> tcell
		screen.ChannelEvents(eventCh, quitCh)
	}(s.screen, quitCh)

	go func(s *Screen, eventCh chan tcell.Event, messageCh chan<- core.Message, quitCh <-chan struct{}) {
		// user <-> matcha
		for {
			select {
			case event, ok := <-eventCh:
				if !ok {
					// NOTE: Means event channel has been closed, either through the quitCh or Fini
					return
				}
				messageCh <- s.mapTCellEventToMessage(event)
			case <-quitCh:
				return
			}
		}
	}(s, eventCh, messageCh, quitCh)
}

func (s *Screen) WriteContent(x, y int, ch rune, style tcell.Style) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screen.SetContent(x, y, ch, nil, style)
}

func (s *Screen) mapTCellEventToMessage(event tcell.Event) core.Message {
	switch e := event.(type) {
	case *tcell.EventKey:
		return &core.MessageKey{EventKey: *e}
	case *tcell.EventResize:
		return &core.MessageResize{EventResize: *e}
	case *tcell.EventClipboard:
		return &core.MessageClipboard{EventClipboard: *e}
	case *tcell.EventError:
		return &core.MessageError{EventError: *e}
	case *tcell.EventFocus:
		return &core.MessageFocus{EventFocus: *e}
	case *tcell.EventInterrupt:
		return &core.MessageInterrupt{EventInterrupt: *e}
	case *tcell.EventMouse:
		return &core.MessageMouse{EventMouse: *e}
	case *tcell.EventPaste:
		return &core.MessagePaste{EventPaste: *e}
	case *tcell.EventTime:
		return &core.MessageTime{EventTime: *e}
	}
	return nil
}
