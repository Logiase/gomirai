package gomirai

import (
	"github.com/Logiase/gomirai/bot"
	"github.com/Logiase/gomirai/message"
)

func SendGroupMessageWithBot(b *bot.Bot, qq, quote uint, msg ...message.Message) (uint, error) {
	return b.SendGroupMessage(qq, quote, msg...)
}

func SendFriendMessageWithBot(b *bot.Bot, group, quote uint, msg ...message.Message) (uint, error) {
	return b.SendGroupMessage(group, quote, msg...)
}
