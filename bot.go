package gomirai

import (
	"github.com/sirupsen/logrus"

	"github.com/logiase/gomirai/tools"
)

type Bot struct {
	QQ         int64
	SessionKey string

	Client *Client

	Logger *logrus.Entry
}

func (b *Bot) sendFriendMessage(qq, quote int64, msg MessageChain) (int64, error) {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "qq": qq, "messageChain": msg}
	if quote != 0 {
		data["quote"] = quote
	}
	res, err := b.Client.doPost("/sendFriendMessage", data)
	if err != nil {
		return 0, err
	}
	b.Logger.Info("Send FriendMessage to", qq)
	return tools.Json.Get([]byte(res), "messageId").ToInt64(), nil
}

func (b *Bot) sendGroupMessage(group, quote int64, msg MessageChain) (int64, error) {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "group": group, "messageChain": msg}
	if quote != 0 {
		data["quote"] = quote
	}
	res, err := b.Client.doPost("/sendGroupMessage", data)
	if err != nil {
		return 0, err
	}
	b.Logger.Info("Send FriendMessage to", group)
	return tools.Json.Get([]byte(res), "messageId").ToInt64(), nil
}
