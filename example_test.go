package gomirai

import (
	"fmt"
	"github.com/Logiase/gomirai/api"
	"net/http"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// 自定义的链接地址
	bot := NewBot("http://127.0.0.1:1234")
	// 或使用自定义的Client
	bot = NewBotWithClient("http://127.0.0.1:1234", http.Client{
		// 自定义
	})

	// authKey
	ok, err := bot.Auth("12345678")
	if !ok {
		// handle error
		fmt.Println(err)
		return
	}

	ok, err = bot.Verify(123456789)
	if !ok {
		// handle error
		fmt.Println(err)
		return
	}

	//// 方法1
	// 在协程中开始获取消息，消息传输至Bot.MessageChan
	// 请先使用Bot.InitChannel初始化Channel
	bot.InitChannel(20, time.Second)
	// 忽略错误
	go bot.StartChannel()
	// 处理错误
	go func() {
		err = bot.StartChannel()
		if err != nil {
			// 处理错误
		}
	}()
	for {
		e := <-bot.MsgChan
		switch e.Type {
		case "GroupMessage":
			msg := api.Message{Type: "Plain", Text: "hello hello"}
			bot.SendGroupMessage()
		case "FriendMessage":
		//do something
		default:
			// do something
		}
	}

	//// 方法2
	// 或手动获取消息
	resp, err := bot.FetchMessage(10)
	if err != nil {
		// 处理错误
	}
}
