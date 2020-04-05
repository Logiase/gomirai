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

// Bot qq机器人
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

// NewBot :)
func NewBot(addr string) *Bot {
	return &Bot{
		addr: addr,
	}
}

// NewBotWithClient :)
func NewBotWithClient(addr string, c http.Client) *Bot {
	b := NewBot(addr)
	b.client = c
	return b
}

// Auth :)
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

// Verify :)
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

// SendFriendMessage :)
func (b *Bot) SendFriendMessage(msg api.MessageCall) (resp *api.Response, e error) {
	resp = &api.Response{}
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	e = b.call("POST", "/sendFriendMessage", nil, buf, &resp)
	return
}

// SendGroupMessage :)
func (b *Bot) SendGroupMessage(msg api.MessageCall) (resp *api.Response, e error) {
	resp = &api.Response{}
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	e = b.call("POST", "/sendGroupMessage", nil, buf, &resp)
	return
}

// SendImageMessage :)
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

// Recall :)
func (b *Bot) Recall(msg api.MessageCall) (f bool, e error) {
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	resp := make(map[string]interface{})
	e = b.call("POST", "/recall", nil, buf, &resp)
	return checkUniformCodeResp(resp)
}

// FetchMessage :)
func (b *Bot) FetchMessage(count int) (resp []api.Event, e error) {
	resp = make([]api.Event, 0, count)
	e = b.call("GET", "/fetchMessage", url.Values{
		"sessionKey": []string{b.session},
		"count":      []string{strconv.Itoa(count)},
	}, nil, &resp)
	return
}

// MessageFromID :)
func (b *Bot) MessageFromID(id int64) (resp api.Event, e error) {
	e = b.call("GET", "/messageFromId", url.Values{
		"sessionKey": []string{b.session},
		"id":         []string{strconv.FormatInt(id, 10)},
	}, nil, &resp)
	return
}

// RefreshFriendList :)
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
	if !b.flagFriend {
		return b.RefreshFriendList()
	}
	return b.friendList, nil
}

// RefreshGroupList :)
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

// MemberList :)
func (b *Bot) MemberList(target int64) (list []api.GroupMember, e error) {
	list = make([]api.GroupMember, 0)
	e = b.call("GET", "/memberList", url.Values{
		"sessionKey": []string{b.session},
		"target":     []string{strconv.FormatInt(target, 10)},
	}, nil, &list)
	return
}

// MuteAll :)
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

// UnmuteAll :)
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

// Mute :)
func (b *Bot) Mute(msg api.ManageCall) (bool, error) {
	return b.manageCall("/mute", msg)
}

// Unmute :)
func (b *Bot) Unmute(msg api.ManageCall) (bool, error) {
	return b.manageCall("/unmute", msg)
}

// Kick :)
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

// GroupConfig :)
func (b *Bot) GroupConfig(msg api.ConfigCall) (f bool, e error) {
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	resp := make(map[string]interface{})
	e = b.call("POST", "/groupConfig", nil, buf, &resp)
	return checkUniformCodeResp(resp)
}

// GetGroupConfig :)
func (b *Bot) GetGroupConfig(target int64) (resp api.GroupConfig, e error) {
	e = b.call("GET", "/groupConfig", url.Values{
		"sessionKey": []string{b.session},
		"target":     []string{strconv.FormatInt(target, 10)},
	}, nil, &resp)
	return
}

// MemberInfo :)
func (b *Bot) MemberInfo(msg api.ConfigCall) (f bool, e error) {
	buf := bytes.NewBuffer([]byte{})
	if e = json.NewEncoder(buf).Encode(&msg); e != nil {
		return
	}
	resp := make(map[string]interface{})
	e = b.call("POST", "/memberInfo", nil, buf, &resp)
	return checkUniformCodeResp(resp)
}

// GetMemberInfo :)
func (b *Bot) GetMemberInfo(target int64) (resp api.GroupConfig, e error) {
	e = b.call("GET", "/memberInfo", url.Values{
		"sessionKey": []string{b.session},
		"target":     []string{strconv.FormatInt(target, 10)},
	}, nil, &resp)
	return
}

// QQ :)
func (b *Bot) QQ() int64 {
	return b.qq
}

// Session :)
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
		_, _ = sb.WriteString("/")
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
