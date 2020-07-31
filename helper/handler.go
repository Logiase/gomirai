package helper

import (
	"github.com/Logiase/gomirai/message"
)

type EventHandler struct {
	privateMessageHandlers []func(bot interface{}, chain message.Chain, sender message.Sender)
	groupMessageHandlers   []func(bot interface{}, chain message.Chain, sender message.Sender)
}
