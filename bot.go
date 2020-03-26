package gomirai

import (
	"errors"
	"fmt"
	"time"
)

type Bot struct {
	Client *Client

	QQ int64

	MessageChan chan InEvent
	chanCache   int
	currentSize int
	fetchTime   time.Duration

	Session string
}

// Release 释放
func (bot *Bot) Release() error {
	postBody := make(map[string]interface{}, 2)
	postBody["qq"] = bot.QQ
	postBody["session"] = bot.Session

	var respS CommonResponse
	err := bot.Client.httpPost("/release", postBody, &respS)
	if err != nil {
		return err
	}
	if respS.Code != 0 {
		return errors.New(respS.Msg)
	}
	return nil
}

// SendFriendQuoteMessage 发送好友消息
func (bot *Bot) SendFriendQuoteMessage(target, quote int64, msg []Message) (int64, error) {
	postBody := make(map[string]interface{}, 4)
	postBody["sessionKey"] = bot.Session
	postBody["target"] = target
	postBody["quote"] = quote
	postBody["messageChain"] = msg

	var respS MessageResponse

	err := bot.Client.httpPost("/sendFriendMessage", postBody, &respS)
	if err != nil {
		return 0, err
	}
	if respS.Code != 0 {
		return 0, errors.New(respS.Msg)
	}
	return respS.MessageID, nil
}

// SendFriendMessage 发送好友消息
func (bot *Bot) SendFriendMessage(target int64, msg []Message) (int64, error) {
	fmt.Println("SendFriendMessage")
	postBody := MessageCall{}
	postBody.MessageChain = msg
	postBody.SessionKey = bot.Session
	postBody.Target = target

	var respS MessageResponse

	err := bot.Client.httpPost("/sendFriendMessage", postBody, &respS)

	if err != nil {
		return 0, err
	}
	if respS.Code != 0 {

		return 0, errors.New(respS.Msg)
	}
	return respS.MessageID, nil
}

// InitChannel 初始化消息管道
// size 缓存数量 t 每次Fetch的时间间隔
func (bot *Bot) InitChannel(size int, t time.Duration) {
	bot.MessageChan = make(chan InEvent, size)
	bot.currentSize = 0
	bot.fetchTime = t
}

// FetchMessage 获取消息，会阻塞当前线程，消息保存在bot中的MessageChan
// 使用前请使用InitChannel初始化Channel
func (bot *Bot) FetchMessage() error {
	var respS []InEvent
	t := time.NewTicker(bot.fetchTime)
	for {
		err := bot.Client.httpGet("/fetchMessage?sessionKey="+bot.Session+"&count="+string(bot.chanCache), &respS)
		if err != nil {
			return err
		}

		for _, e := range respS {
			if len(bot.MessageChan) == bot.chanCache {
				<-bot.MessageChan
			}
			bot.MessageChan <- e
		}

		<-t.C
	}
}
