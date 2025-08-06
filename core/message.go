package core

import (
	"github.com/gdamore/tcell/v2"
)

type Message interface {
	tcell.Event
}

type MessageKey struct {
	tcell.EventKey
}

type MessageResize struct {
	tcell.EventResize
}

type MessageClipboard struct {
	tcell.EventClipboard
}

type MessageError struct {
	tcell.EventError
}

type MessageFocus struct {
	tcell.EventFocus
}

type MessageInterrupt struct {
	tcell.EventInterrupt
}

type MessageMouse struct {
	tcell.EventMouse
}

type MessagePaste struct {
	tcell.EventPaste
}

type MessageTime struct {
	tcell.EventTime
}
