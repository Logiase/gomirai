package gomirai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Logiase/gomirai/api"
)

// Bot qq机器人结构体
type Bot struct {
	addr, authKey, session string
	qq                     int64

	flagFriend  bool
	flagGroup   bool
	chanCache   int
	currentSize int
	fetchTime   time.Duration

	MsgChan    chan api.Event
	friendList []api.Friend
	groupList  []api.Group

	client http.Client
}

// NewBot 新建一个以addr为api地址的Bot
func NewBot(addr string) *Bot {
	return &Bot{
		addr: addr,
	}
}

// NewBotWithClient 新建一个以addr为api地址、以c为默认http.Client的Bot
func NewBotWithClient(addr string, c http.Client) *Bot {
	b := NewBot(addr)
	b.client = c
	return b
}

// Auth 验证一个authKey，验证成功时将Bot的session设置为api的对应返回值
func (b *Bot) Auth(authKey string) (f bool, e error) {
	defer func() {
		if p := recover(); p != nil {
			e = p.(error)
		}
	}()
	var i = make(map[string]interface{})
	payload := `{"authKey": "` + authKey + `"}`
	e = b.call("POST", "/auth", nil, bytes.NewReader([]byte(payload)), &i)
	if e != nil {
		return false, e
	}

	f = i["code"].(int) == 0
	b.session = i["session"].(string)
	return
}

// Verify 校验并激活你的Session，同时将该Session与一个已登陆的Bot绑定
// 参数qq为希望绑定的Bot的QQ号
func (b *Bot) Verify(qq int64) (f bool, e error) {
	i := make(map[string]interface{})
	payload := `{"sessionKey": "` + b.session + `", "qq": ` + strconv.FormatInt(qq, 10) + `}`
	e = b.call("POST", "/verify", nil, bytes.NewReader([]byte(payload)), &i)
	if e != nil {
		return
	}
	if f, e = checkUniformCodeResp(i); f {
		b.qq = qq
	}
	return
}

// Release 释放一个session及其相关的资源，
// 不使用的session应当被释放，否则将导致backend的内存泄漏
func (b *Bot) Release(qq int64) (f bool, e error) {
	i := make(map[string]interface{})
	payload := `{"sessionKey": "` + b.session + `", "qq": ` + strconv.FormatInt(qq, 10) + `}`
	e = b.call("POST", "/release", nil, bytes.NewReader([]byte(payload)), &i)
	if e != nil {
		return
	}
	if f, e = checkUniformCodeResp(i); f {
		b.qq = 0
	}
	return
}

// SendFriendMessage 向指定好友发送消息
func (b *Bot) SendFriendMessage(msg api.MessageCall) (resp *api.Response, e error) {
	resp = &api.Response{}
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	e = b.call("POST", "/sendFriendMessage", nil, buf, &resp)
	return
}

// SendGroupMessage 向指定群发送消息
func (b *Bot) SendGroupMessage(msg api.MessageCall) (resp *api.Response, e error) {
	resp = &api.Response{}
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	e = b.call("POST", "/sendGroupMessage", nil, buf, &resp)
	return
}

// SendImageMessage 通过URL发送图片消息
func (b *Bot) SendImageMessage(msg api.MessageCall) (resp []string, e error) {
	resp = make([]string, 0)
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	e = b.call("POST", "/sendImageMessage", nil, buf, &resp)
	return
}

// TODO: multipart/form-data
// func (b *Bot) UploadImage() {}

// Recall 撤回指定消息
// 对于自己发送的消息，有2分钟的时间限制
// 对于群聊中的群员消息，需要有相应的权限
func (b *Bot) Recall(msg api.MessageCall) (f bool, e error) {
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	resp := make(map[string]interface{})
	e = b.call("POST", "/recall", nil, buf, &resp)
	return checkUniformCodeResp(resp)
}

// FetchMessage 获取Bot接收到的消息和各类事件
func (b *Bot) FetchMessage(count int) (resp []api.Event, e error) {
	resp = make([]api.Event, 0, count)
	e = b.call("GET", "/fetchMessage", url.Values{
		"sessionKey": []string{b.session},
		"count":      []string{strconv.Itoa(count)},
	}, nil, &resp)
	return
}

// MessageFromID 获取Bot接收到的、被缓存的消息和各类事件
func (b *Bot) MessageFromID(id int64) (resp api.Event, e error) {
	e = b.call("GET", "/messageFromId", url.Values{
		"sessionKey": []string{b.session},
		"id":         []string{strconv.FormatInt(id, 10)},
	}, nil, &resp)
	return
}

// RefreshFriendList 刷新Bot的好友列表
func (b *Bot) RefreshFriendList() (list []api.Friend, e error) {
	list = make([]api.Friend, 0)
	e = b.call("GET", "/friendList", url.Values{
		"sessionKey": []string{b.session},
	}, nil, &list)

	b.friendList = list
	b.flagFriend = true

	return
}

// FriendList 获取缓存的好友列表
func (b *Bot) FriendList() (list []api.Friend, e error) {
	// POTENTIAL RACE?
	if !b.flagFriend {
		return b.RefreshFriendList()
	}
	return b.friendList, nil
}

// RefreshGroupList 刷新Bot的群列表
func (b *Bot) RefreshGroupList() (list []api.Group, e error) {
	list = make([]api.Group, 0)
	e = b.call("GET", "/groupList", url.Values{
		"sessionKey": []string{b.session},
	}, nil, &list)

	b.groupList = list
	b.flagGroup = true

	return
}

// GroupList 获取缓存的群列表
func (b *Bot) GroupList() (list []api.Group, e error) {
	if !b.flagGroup {
		return b.RefreshGroupList()
	}
	return b.groupList, nil
}

// MemberList 获取指定群中的成员列表
func (b *Bot) MemberList(target int64) (list []api.GroupMember, e error) {
	list = make([]api.GroupMember, 0)
	e = b.call("GET", "/memberList", url.Values{
		"sessionKey": []string{b.session},
		"target":     []string{strconv.FormatInt(target, 10)},
	}, nil, &list)
	return
}

// MuteAll 对指定群进行全体禁言（需要有相关权限）
func (b *Bot) MuteAll(target int64) (f bool, e error) {
	sb := strings.Builder{}
	_, _ = sb.WriteString(`{"sessionKey": "`)
	_, _ = sb.WriteString(b.session)
	_, _ = sb.WriteString(`", "target": `)
	_, _ = sb.WriteString(strconv.FormatInt(target, 10))
	_, _ = sb.WriteString(`}`)

	resp := make(map[string]interface{})
	e = b.call("POST", "/muteAll", nil, bytes.NewReader([]byte(sb.String())), &resp)
	return checkUniformCodeResp(resp)
}

// UnmuteAll 对指定群解除全体禁言（需要有相关权限）
func (b *Bot) UnmuteAll(target int64) (f bool, e error) {
	sb := strings.Builder{}
	_, _ = sb.WriteString(`{"sessionKey": "`)
	_, _ = sb.WriteString(b.session)
	_, _ = sb.WriteString(`", "target": `)
	_, _ = sb.WriteString(strconv.FormatInt(target, 10))
	_, _ = sb.WriteString(`}`)

	resp := make(map[string]interface{})
	e = b.call("POST", "/unmuteAll", nil, bytes.NewReader([]byte(sb.String())), &resp)
	return checkUniformCodeResp(resp)
}

// Mute 对指定群中的指定群员进行禁言（需要有相关权限）
func (b *Bot) Mute(msg api.ManageCall) (bool, error) {
	return b.manageCall("/mute", msg)
}

// Unmute 对指定群中的指定群员解除禁言（需要有相关权限）
func (b *Bot) Unmute(msg api.ManageCall) (bool, error) {
	return b.manageCall("/unmute", msg)
}

// Kick 将指定群中的指定群员踢出（需要有相关权限）
func (b *Bot) Kick(msg api.ManageCall) (bool, error) {
	return b.manageCall("/kick", msg)
}

func (b *Bot) manageCall(endpoint string, msg api.ManageCall) (f bool, e error) {
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	resp := make(map[string]interface{})
	e = b.call("POST", endpoint, nil, buf, &resp)
	return checkUniformCodeResp(resp)
}

// GroupConfig 修改群设置（需要有相关权限）
func (b *Bot) GroupConfig(msg api.ConfigCall) (f bool, e error) {
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	resp := make(map[string]interface{})
	e = b.call("POST", "/groupConfig", nil, buf, &resp)
	return checkUniformCodeResp(resp)
}

// GetGroupConfig 获取群设置
func (b *Bot) GetGroupConfig(target int64) (resp api.GroupConfig, e error) {
	e = b.call("GET", "/groupConfig", url.Values{
		"sessionKey": []string{b.session},
		"target":     []string{strconv.FormatInt(target, 10)},
	}, nil, &resp)
	return
}

// MemberInfo 修改群员资料（需要有相关权限）
func (b *Bot) MemberInfo(msg api.ConfigCall) (f bool, e error) {
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	resp := make(map[string]interface{})
	e = b.call("POST", "/memberInfo", nil, buf, &resp)
	return checkUniformCodeResp(resp)
}

// GetMemberInfo 获取群员资料
func (b *Bot) GetMemberInfo(target int64) (resp api.GroupConfig, e error) {
	e = b.call("GET", "/memberInfo", url.Values{
		"sessionKey": []string{b.session},
		"target":     []string{strconv.FormatInt(target, 10)},
	}, nil, &resp)
	return
}

// QQ 返回Bot对应的QQ号
func (b *Bot) QQ() int64 {
	return b.qq
}

// Session 返回Bot对应的session
func (b *Bot) Session() string {
	return b.session
}

func checkUniformCodeResp(m map[string]interface{}) (f bool, e error) {
	defer func() {
		if p := recover(); p != nil {
			e = p.(error)
		}
	}()
	if f = m["code"].(int) == 0; !f {
		e = errors.New("code: " + strconv.Itoa(m["code"].(int)) + "| Msg: " + m["msg"].(string))
	}
	return
}

func (b *Bot) call(method, endpoint string, params url.Values, body io.Reader, response interface{}) (e error) {
	sb := strings.Builder{}
	// (strings.Builder).WriteString never returns non nil error.
	_, _ = sb.WriteString(b.addr)
	_, _ = sb.WriteString(endpoint)
	if params != nil {
		_, _ = sb.WriteString("?")
		_, _ = sb.WriteString(params.Encode())
	}

	req, e := http.NewRequest(method, sb.String(), body)
	if e != nil {
		return
	}
	req.Header.Add("Connection", "Keep-Alive")
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, e := b.client.Do(req)
	if e != nil {
		return
	}
	e = json.NewDecoder(resp.Body).Decode(response)
	resp.Body.Close()

	return
}
