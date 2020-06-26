package gomirai

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// 链接地址
	address := "http://127.0.0.1:8080"
	authKey := "12345678"
	// 用于进行网络操作的Client
	client := NewMiraiClient(address, authKey)

	// 可对Client做出自定义修改，该修改会应用于所有使用该client的网络请求
	// 如使用Proxy
	client.HTTPClient.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}

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
	// FetchMessage时间间隔 1s,每次获取的数量20,channel缓存容量20
	bot.InitChannel(20, time.Second)

	// 在协程中开始获取消息，消息传输至Bot.MessageChan
	// 忽略错误
	go bot.FetchMessage()
	// 检查错误
	go func() {
		err = bot.FetchMessage()
		if err != nil {
			//handle Error
		}
	}()

	// 从bot.MessageChan获取收到事件并处理
	for {
		e := <-bot.MessageChan
		switch e.Type {
		case "GroupMessage": // do something
			// some func(e)
		case "FriendMessage": // do something
		case "....省略": // do something
		default:
			// do something
		}
	}
}
