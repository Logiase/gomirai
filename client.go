package gomirai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Client 用于服务的Client
type Client struct {
	Address string
	// 可以对HTTPClient进行自定义配置
	HTTPClient http.Client

	authKey string
	Bots    map[int64]*Bot
}

// NewMiraiClient 创建一个新的Client
func NewMiraiClient(address, authKey string) (client *Client) {
	client = &Client{}
	client.Address = address
	client.Bots = make(map[int64]*Bot)
	client.HTTPClient = http.Client{}
	client.authKey = authKey
	return
}

// Verify 使用此方法校验并激活你的Session，同时将Session与一个已登录的Bot绑定
func (client *Client) Verify(qq int64) (*Bot, error) {
	session, err := client.Auth()
	if err != nil {
		return nil, err
	}
	postBody := make(map[string]interface{}, 2)
	postBody["sessionKey"] = session
	postBody["qq"] = qq

	var respS Response
	err = client.httpPost("/verify", postBody, &respS)
	if err != nil {
		return nil, err
	}
	if respS.Code != 0 {
		return nil, errors.New(respS.Msg)
	}

	client.Bots[qq] = &Bot{
		Client:  client,
		QQ:      qq,
		Session: session,
	}
	return client.Bots[qq], nil
}

// Auth 使用此方法验证你的身份，并返回一个Session
func (client *Client) Auth() (session string, err error) {
	postBody := make(map[string]interface{}, 1)
	postBody["authKey"] = client.authKey
	var respS AuthResponse
	err = client.httpPost("/auth", postBody, &respS)
	if err != nil {
		return
	}
	if respS.Code != 0 {
		err = errors.New("错误的MIRAI API HTTP auth key")
		return
	}
	session = respS.Session
	return
}

// ReleaseAllSession 释放所有会话
func (client *Client) ReleaseAllSession() {
	for _, bot := range client.Bots {
		_ = bot.Release()
	}
}

// 用于内部的Post请求
func (client *Client) httpPost(path string, postBody interface{}, respS interface{}) error {
	bytesData, _ := json.Marshal(postBody)
	req, err := http.NewRequest("POST", client.Address+path, bytes.NewReader(bytesData))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "Keep-Alive")
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bytesData, err = ioutil.ReadAll(resp.Body)

	logrus.Debug(string(bytesData))
	return json.Unmarshal(bytesData, respS)
}

// 用于内部的Get请求
func (client *Client) httpGet(path string, respS interface{}) error {
	req, err := http.NewRequest("GET", client.Address+path, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Connection", "Keep-Alive")
	resp, err := client.HTTPClient.Get(client.Address + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytesData, err := ioutil.ReadAll(resp.Body)
	logrus.Debug(string(bytesData))
	return json.Unmarshal(bytesData, respS)
}

// GetSessionConfig 获取Session的设置
func (client *Client) GetSessionConfig(session string) (*SessionConfig, error) {
	var respS SessionConfig
	err := client.httpGet("/config?sessionKey="+session, &respS)
	if err != nil {
		return nil, err
	}
	respS.SessionKey = session
	return &respS, nil
}

// SetSessionConfig 设置Session
func (client *Client) SetSessionConfig(config *SessionConfig) error {
	var respS Response
	err := client.httpPost("/config", config, &respS)
	if err != nil {
		return err
	}
	if respS.Code != 0 {
		return errors.New(respS.Msg)
	}
	return nil
}

// SendCommand 发送指令
func (client *Client) SendCommand(commandName string, args ...string) string {
	postBody := make(map[string]interface{}, 3)
	postBody["authKey"] = client.authKey
	postBody["name"] = commandName
	postBody["args"] = args
	bytesData, _ := json.Marshal(postBody)
	resp, err := client.HTTPClient.Post("/command/send", "application/json", bytes.NewReader(bytesData))
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	bytesData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	return string(bytesData)
}

// GetManagers 获取指定bot的Managers
func (client *Client) GetManagers(qq int64) (*[]int64, error) {
	var respS []int64
	err := client.httpGet("/managers?qq="+strconv.FormatInt(qq, 10), &respS)
	if err != nil {
		return nil, err
	}
	return &respS, nil
}
