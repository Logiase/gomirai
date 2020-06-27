package gomirai

import "github.com/Logiase/gomirai/message"

func SendGroupMessageWithBot(b *Bot, qq, quote uint, msg ...message.Message) (uint, error) {
	return b.SendGroupMessage(qq, quote, msg...)
}

func SendFriendMessageWithBot(b *Bot, group, quote uint, msg ...message.Message) (uint, error) {
	return b.SendGroupMessage(group, quote, msg...)
}
