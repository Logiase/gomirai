package gomirai

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// 链接地址
	address := "http://127.0.0.1:8080"
	authKey := "12345678"
	// 用于进行网络操作的Client
	client := NewMiraiClient(address, authKey)
	// 获取Bot，Session信息保存在Bot中
	// 也可通过Client.Bots[]获取
	bot, err := client.Verify(123456789)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 释放资源
	defer bot.Release()

	// 初始化消息通道
	bot.InitChannel(20, time.Second)

	// 在协程中开始获取消息，消息传输至Bot.MessageChan
	go bot.FetchMessage()

	for {
		e := <-bot.MessageChan
		if e.Type == "FriendMessage" {
			// do something
		}
	}
}
