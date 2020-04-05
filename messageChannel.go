package gomirai

import (
	"github.com/Logiase/gomirai/api"
	"net/url"
	"strconv"
	"time"
)

// InitChannel 初始化一个接收消息的通道
func (b *Bot) InitChannel(size int, t time.Duration) {
	b.MsgChan = make(chan api.Event, size)
	b.chanCache = size
	b.currentSize = 0
	b.fetchTime = t
}

// StartChannel 开始获取消息，会导致阻塞
func (b *Bot) StartChannel() error {
	var resp []api.Event
	t := time.NewTicker(b.fetchTime)

	for {
		err := b.call("GET", "/fetchMessage", url.Values{
			"sessionKey": []string{b.session},
			"count":      []string{strconv.Itoa(b.chanCache)},
		}, nil, &resp)
		if err != nil {
			return err
		}

		for _, e := range resp {
			if len(b.MsgChan) == b.chanCache {
				<-b.MsgChan
			}
			b.MsgChan <- e
		}

		<-t.C
	}
}
