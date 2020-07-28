package helper

import (
	"github.com/Logiase/gomirai/bot"
	"github.com/Logiase/gomirai/message"
)

type EventHandler struct {
	privateMessageHandlers []func(bot *bot.Bot, chain message.Chain, sender message.Sender)
	groupMessageHandlers   []func(bt *bot.Bot, chain message.Chain, sender message.Sender)
}
