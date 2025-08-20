package matcha

import "github.com/gdamore/tcell/v2"

type handler struct {
	ch   chan tcell.Event
	quit chan struct{}
}

func newHandler() *handler {
	return &handler{ch: make(chan tcell.Event)}
}

func (h *handler) getChannel() chan tcell.Event {
	return h.ch
}

func (h *handler) getQuitChannel() chan struct{} {
	return h.quit
}
