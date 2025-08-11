package matcha

import "github.com/cchirag/matcha/core"

type (
	Message          = core.Message
	MessageKey       = core.MessageKey
	MessageResize    = core.MessageResize
	MessageClipboard = core.MessageClipboard
	MessageError     = core.MessageError
	MessageFocus     = core.MessageFocus
	MessageInterrupt = core.MessageInterrupt
	MessageMouse     = core.MessageMouse
	MessagePaste     = core.MessagePaste
	MessageTime      = core.MessageTime
)

type messageEntry struct {
	action func(message Message)
}
