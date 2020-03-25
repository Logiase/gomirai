package gomirai

import (
	"errors"
	"fmt"
)

type Bot struct {
	Client *Client

	QQ int64

	Session string
}

// Release 释放
func (bot Bot) Release() error {
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

// SendFriendQuoteMessage
func (bot Bot) SendFriendQuoteMessage(target, quote int64, msg []Message) (int64, error) {
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

func (bot Bot) SendFriendMessage(target int64, msg []Message) (int64, error) {
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
