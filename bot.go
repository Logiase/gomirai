package gomirai

import (
	"errors"
	"strconv"
	"strings"
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

// Release 使用此方式释放session及其相关资源（Bot不会被释放）
// 不使用的Session应当被释放，长时间（30分钟）未使用的Session将自动释放，否则Session持续保存Bot收到的消息，将会导致内存泄露
func (bot *Bot) Release() error {
	if bot.Session == "" {
		return errors.New("bot未实例化")
	}

	postBody := make(map[string]interface{}, 2)
	postBody["qq"] = bot.QQ
	postBody["session"] = bot.Session

	var respS Response
	err := bot.Client.httpPost("/release", postBody, &respS)
	if err != nil {
		return err
	}
	if respS.Code != 0 {
		return errors.New(respS.Msg)
	}
	return nil
}

// SendFriendQuoteMessage 使用此方法向指定好友发送消息
// 如果不需要引用回复，quote设0
func (bot *Bot) SendFriendMessage(target, quote int64, msg []Message) (int64, error) {
	postBody := make(map[string]interface{})
	postBody["sessionKey"] = bot.Session
	postBody["target"] = target
	if quote != 0 {
		postBody["quote"] = quote
	}
	postBody["messageChain"] = msg

	var respS Response
	err := bot.Client.httpPost("/sendFriendMessage", postBody, &respS)
	if err != nil {
		return 0, err
	}
	if respS.Code != 0 {
		return 0, errors.New(respS.Msg)
	}
	return respS.MessageID, nil
}

// SendGroupMessage 使用此方法向指定群发送消息
func (bot *Bot) SendGroupMessage(target, quote int64, msg []Message) (int64, error) {
	postBody := make(map[string]interface{})
	postBody["sessionKey"] = bot.Session
	postBody["target"] = target
	if quote != 0 {
		postBody["quote"] = quote
	}
	postBody["messageChain"] = msg

	var respS Response
	err := bot.Client.httpPost("/sendGroupMessage", postBody, &respS)
	if err != nil {
		return 0, err
	}
	if respS.Code != 0 {
		return 0, errors.New(respS.Msg)
	}
	return respS.MessageID, nil
}

// SendImageMessage 使用此方法向指定对象（群或好友）发送图片消息 除非需要通过此手段获取imageId，否则不推荐使用该接口
func (bot *Bot) SendImageMessage(target int64, targetType string, urls []string) ([]string, error) {
	postBody := make(map[string]interface{})
	postBody["sessionKey"] = bot.Session
	switch strings.ToLower(targetType) {
	case "group":
		postBody["group"] = target
	case "qq":
		postBody["qq"] = target
	default:
		return nil, errors.New("target Type错误 应为 qq 或 group")
	}
	postBody["urls"] = urls

	var respS []string
	err := bot.Client.httpPost("/sendGroupMessage", postBody, &respS)
	if err != nil {
		return nil, err
	}
	return respS, nil
}

// InitChannel 初始化消息管道
// size 缓存数量 t 每次Fetch的时间间隔
func (bot *Bot) InitChannel(size int, t time.Duration) {
	bot.MessageChan = make(chan InEvent, size)
	bot.chanCache = size
	bot.currentSize = 0
	bot.fetchTime = t
}

// FetchMessage 获取消息，会阻塞当前线程，消息保存在bot中的MessageChan
// 使用前请使用InitChannel初始化Channel
func (bot *Bot) FetchMessage() error {
	var respS []InEvent
	t := time.NewTicker(bot.fetchTime)
	for {
		err := bot.Client.httpGet("/fetchMessage?sessionKey="+bot.Session+"&count="+strconv.Itoa(bot.chanCache), &respS)
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
