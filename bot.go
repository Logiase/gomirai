package gomirai

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/logiase/gomirai/message"
	"github.com/logiase/gomirai/tools"
)

type Bot struct {
	QQ         uint
	SessionKey string

	Client *Client

	Logger *logrus.Entry

	fetchTime   time.Duration
	size        int
	currentSize int
	Chan        chan message.Event
}

func (b *Bot) SetChannel(time time.Duration, size int) {
	b.Chan = make(chan message.Event, size)
	b.size = size
	b.currentSize = 0
	b.fetchTime = time
}

func (b *Bot) SendFriendMessage(qq, quote uint, msg ...message.Message) (uint, error) {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "qq": qq, "messageChain": msg}
	if quote != 0 {
		data["quote"] = quote
	}
	res, err := b.Client.doPost("/sendFriendMessage", data)
	if err != nil {
		return 0, err
	}
	b.Logger.Info("Send FriendMessage to", qq)
	return tools.Json.Get([]byte(res), "messageId").ToUint(), nil
}

func (b *Bot) SendGroupMessage(group, quote uint, msg ...message.Message) (uint, error) {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "group": group, "messageChain": msg}
	if quote != 0 {
		data["quote"] = quote
	}
	res, err := b.Client.doPost("/sendGroupMessage", data)
	if err != nil {
		return 0, err
	}
	b.Logger.Info("Send FriendMessage to", group)
	return tools.Json.Get([]byte(res), "messageId").ToUint(), nil
}

func (b *Bot) FetchMessages() error {
	t := time.NewTicker(b.fetchTime)

	for {
		res, err := b.Client.doGet("/fetchMessage", map[string]string{
			"sessionKey": b.SessionKey,
			"count":      strconv.Itoa(b.size),
		})
		if err != nil {
			return err
		}
		var tc []message.Event
		tools.Json.Get([]byte(res), "data").ToVal(&tc)
		for _, v := range tc {
			if len(b.Chan) == b.size {
				<-b.Chan
			}
			b.Chan <- v
		}

		<-t.C
	}
}
