package gomirai

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Bot qq机器人
type Bot struct {
	Client *Client

	QQ int64

	MessageChan chan InEvent
	chanCache   int
	currentSize int
	fetchTime   time.Duration

	Session string

	flagFriend bool
	friendList []*Friend

	flagGroup bool
	groupList []*Group
}

// Release 使用此方式释放session及其相关资源（Bot不会被释放）
// 不使用的Session应当被释放，长时间（30分钟）未使用的Session将自动释放，否则Session持续保存Bot收到的消息，将会导致内存泄露
func (bot *Bot) Release() error {
	if bot.Session == "" {
		return errors.New("bot未实例化")
	}

	postBody := make(map[string]interface{}, 2)
	postBody["qq"] = bot.QQ
	postBody["sessionKey"] = bot.Session

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

// SendFriendMessage 使用此方法向指定好友发送消息
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
	err := bot.Client.httpPost("/sendImageMessage", postBody, &respS)
	if err != nil {
		return nil, err
	}
	return respS, nil
}

// Recall 撤回一条消息
func (bot *Bot) Recall(target int64) error {
	postBody := make(map[string]interface{}, 2)
	postBody["sessionKey"] = bot.Session
	postBody["target"] = target

	var respS Response
	err := bot.Client.httpPost("/recall", postBody, &respS)
	if err != nil {
		return err
	}
	if respS.Code != 0 {
		return errors.New(respS.Msg)
	}
	return nil
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
	var respS InEventAll
	t := time.NewTicker(bot.fetchTime)
	for {
		err := bot.Client.httpGet("/fetchMessage?sessionKey="+bot.Session+"&count="+strconv.Itoa(bot.chanCache), &respS)
		if err != nil {
			return err
		}

		for _, e := range respS.Data {
			if len(bot.MessageChan) == bot.chanCache {
				<-bot.MessageChan
			}
			bot.MessageChan <- e
		}

		<-t.C
	}
}

// MessageFromID 通过ID获取一条缓存的消息
func (bot *Bot) MessageFromID(id int64) (*InEvent, error) {
	var respS InEvent
	err := bot.Client.httpGet("/messageFromId?sessionKey="+bot.Session+"&id="+strconv.FormatInt(id, 10), &respS)
	if err != nil {
		return nil, err
	}
	return &respS, nil
}

// FriendList 获取Bot的好友列表
// 会获取本地缓存的好友列表，如需刷新请使用RefreshFriendList
// 没有缓存时会自动刷新
func (bot *Bot) FriendList() ([]*Friend, error) {
	if bot.flagFriend {
		return bot.friendList, nil
	}
	return bot.RefreshFriendList()
}

// RefreshFriendList 刷新好友列表
func (bot *Bot) RefreshFriendList() ([]*Friend, error) {
	var respS []*Friend
	err := bot.Client.httpGet("/friendList?sessionKey="+bot.Session, &respS)
	if err != nil {
		return nil, err
	}
	bot.friendList = respS
	return respS, nil
}

// GroupList 获取Bot的群列表
// 会获取本地缓存的群列表，如需刷新请使用RefreshGroupList
// 没有缓存时会自动刷新
func (bot *Bot) GroupList() ([]*Group, error) {
	if bot.flagGroup {
		return bot.groupList, nil
	}
	return bot.RefreshGroupList()
}

// RefreshGroupList 刷新群列表s
func (bot *Bot) RefreshGroupList() ([]*Group, error) {
	var respS []*Group
	err := bot.Client.httpGet("/groupList?sessionKey="+bot.Session, &respS)
	if err != nil {
		return nil, err
	}
	bot.groupList = respS
	return respS, nil
}

// MemberList 指定群内的群成员
func (bot *Bot) MemberList(group int64) ([]*GroupMember, error) {
	var respS []*GroupMember
	err := bot.Client.httpGet("/memberList?sessionKey="+bot.Session+"&target="+strconv.FormatInt(group, 10), &respS)
	if err != nil {
		return nil, err
	}
	return respS, nil
}

// MuteAll 全体禁言
func (bot *Bot) MuteAll(group int64) error {
	postBody := make(map[string]interface{}, 2)
	postBody["sessionKey"] = bot.Session
	postBody["target"] = group

	var respS Response
	err := bot.Client.httpPost("/muteAll", postBody, &respS)
	if err != nil {
		return err
	}
	if respS.Code != 0 {
		return errors.New(respS.Msg)
	}
	return nil
}

// UnmuteAll 接触全体禁言
func (bot *Bot) UnmuteAll(group int64) error {
	postBody := make(map[string]interface{}, 2)
	postBody["sessionKey"] = bot.Session
	postBody["target"] = group

	var respS Response
	err := bot.Client.httpPost("/unmuteAll", postBody, &respS)
	if err != nil {
		return err
	}
	if respS.Code != 0 {
		return errors.New(respS.Msg)
	}
	return nil
}

// Mute 禁言 second为0 解除禁言
func (bot *Bot) Mute(group int64, member int64, second int64) error {
	postBody := make(map[string]interface{})
	postBody["sessionKey"] = bot.Session
	postBody["target"] = group
	postBody["memberId"] = member

	var respS Response
	var err error

	if second == 0 {
		err = bot.Client.httpPost("/unmute", postBody, &respS)
	} else {
		postBody["time"] = second
		err = bot.Client.httpPost("/mute", postBody, &respS)
	}

	if err != nil {
		return err
	}
	if respS.Code != 0 {
		return errors.New(respS.Msg)
	}

	return nil
}

// Kick 踢出群聊
func (bot *Bot) Kick(group int64, member int64, msg string) error {
	postBody := make(map[string]interface{}, 4)
	postBody["sessionKey"] = bot.Session
	postBody["target"] = group
	postBody["memberId"] = member
	postBody["msg"] = msg
	var respS Response
	err := bot.Client.httpPost("/mute", postBody, &respS)
	if err != nil {
		return err
	}
	if respS.Code != 0 {
		return errors.New(respS.Msg)
	}
	return nil
}

// SetGroupConfig 设置群设置
func (bot *Bot) SetGroupConfig(config GroupConfig, group int64) error {
	postBody := make(map[string]interface{}, 3)
	postBody["sessionKey"] = bot.Session
	postBody["target"] = group
	postBody["config"] = config
	var respS Response
	err := bot.Client.httpPost("/groupConfig", postBody, &respS)
	if err != nil {
		return err
	}
	if respS.Code != 0 {
		return errors.New(respS.Msg)
	}
	return nil
}

// GetGroupConfig 获取群设置
func (bot *Bot) GetGroupConfig(group int64) (*GroupConfig, error) {
	var respS *GroupConfig
	err := bot.Client.httpGet("/groupConfig?sessionKey="+bot.Session+"&target="+strconv.FormatInt(group, 10), &respS)
	if err != nil {
		return nil, err
	}
	return respS, nil
}

// SetGroupMemberInfo 设置群成员信息
func (bot *Bot) SetGroupMemberInfo(info *GroupMemberInfo, group, member int64) error {
	postBody := make(map[string]interface{}, 3)
	postBody["sessionKey"] = bot.Session
	postBody["target"] = group
	postBody["memberId"] = member
	postBody["info"] = info
	var respS Response
	err := bot.Client.httpPost("/memberInfo", postBody, &respS)
	if err != nil {
		return err
	}
	if respS.Code != 0 {
		return errors.New(respS.Msg)
	}
	return nil
}

// GetGroupMemberInfo 获取群成员信息
func (bot *Bot) GetGroupMemberInfo(group, member int64) (*GroupMemberInfo, error) {
	var respS *GroupMemberInfo
	err := bot.Client.httpGet("/memberInfo?sessionKey="+bot.Session+"&target="+strconv.FormatInt(group, 10)+"&memberId="+strconv.FormatInt(member, 10), &respS)
	if err != nil {
		return nil, err
	}
	return respS, nil
}
