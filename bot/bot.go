package bot

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Logiase/gomirai/message"
	"github.com/Logiase/gomirai/tools"
)

// Bot 对应一个机器人账号
// 进行所有对账号相关操作
type Bot struct {
	QQ         uint
	SessionKey string

	Client *Client

	Logger *logrus.Entry

	fetchTime   time.Duration
	size        int
	currentSize int
	Chan        chan message.Event

	Friends []message.Friend
	Groups  []message.Group
}

// --- Bot 设置 ---

// SetChannel Channel相关设置
func (b *Bot) SetChannel(time time.Duration, size int) {
	b.Chan = make(chan message.Event, size)
	b.size = size
	b.currentSize = 0
	b.fetchTime = time
}

// --- 消息相关 ---

// SendFriendMessage 使用此方法向指定好友发送消息
// qq 好友qq
// quote 引用消息id 0为不引用
// msg 消息内容
func (b *Bot) SendFriendMessage(qq, quote uint, msg ...message.Message) (uint, error) {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "qq": qq, "messageChain": msg}
	if quote != 0 {
		data["quote"] = quote
	}
	res, err := b.Client.doPost("/sendFriendMessage", data)
	if err != nil {
		return 0, err
	}
	b.Logger.Info("Send FriendMessage to ", qq)
	return tools.Json.Get([]byte(res), "messageId").ToUint(), nil
}

// SendTempMessage 使用此方法向临时会话对象发送消息
// qq 好友qq
// group 群qq
// msg 消息内容
func (b *Bot) SendTempMessage(group, qq uint, msg ...message.Message) (uint, error) {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "qq": qq, "group": group, "messageChain": msg}
	res, err := b.Client.doPost("/sendTempMessage", data)
	if err != nil {
		return 0, err
	}
	b.Logger.Info("Send TempMessage to ", qq)
	return tools.Json.Get([]byte(res), "messageId").ToUint(), nil
}

// SendGroupMessage 使用此方法向指定群发送消息
// group 群qq
// quote 引用消息id 0为不引用
// msg 消息内容
func (b *Bot) SendGroupMessage(group, quote uint, msg ...message.Message) (uint, error) {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "group": group, "messageChain": msg}
	if quote != 0 {
		data["quote"] = quote
	}
	res, err := b.Client.doPost("/sendGroupMessage", data)
	if err != nil {
		return 0, err
	}
	b.Logger.Info("Send FriendMessage to", group)
	return tools.Json.Get([]byte(res), "messageId").ToUint(), nil
}

// SendImageMessage 使用此方法向指定对象（群或好友）发送图片消息
// 除非需要通过此手段获取imageId，否则不推荐使用该接口
// 请保证 qq group 不同时有值
func (b *Bot) SendImageMessage(qq, group int64, urls ...string) (imageIds []string, err error) {
	if qq*group == 0 {
		return nil, errors.New("非法参数")
	}
	data := map[string]interface{}{"sessionKey": b.SessionKey, "urls": urls}
	if qq == 0 {
		data["group"] = group
	} else {
		data["qq"] = qq
	}
	res, err := b.Client.doPost("sendImageMessage", data)
	if err != nil {
		return nil, err
	}
	b.Logger.Info("Send Images")
	err = tools.Json.UnmarshalFromString(res, &imageIds)
	return
}

// UploadImage 使用此方法上传图片文件至服务器并返回ImageId
func (b *Bot) UploadImage(t string, imgFilepath string) (string, error) {
	imgReader, err := os.Open(imgFilepath)
	if err != nil {
		return "", err
	}
	defer imgReader.Close()

	data := map[string]interface{}{"sessionKey": b.SessionKey, "type": t, "img": imgReader}
	res, err := b.Client.doPostWithFormData("/uploadImage", data)
	if err != nil {
		return "", err
	}
	b.Logger.Info("UploadFriendImage ", imgFilepath)
	return tools.Json.Get([]byte(res), "imageId").ToString(), nil
}

// Recall 使用此方法撤回指定消息
// 对于bot发送的消息，有2分钟时间限制。对于撤回群聊中群员的消息，需要有相应权限
// target 消息id
func (b *Bot) Recall(target int64) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "target": target}
	_, err := b.Client.doPost("/recall", data)
	return err
}

// FetchMessages 获取消息
func (b *Bot) FetchMessages() error {
	t := time.NewTicker(b.fetchTime)

	for {
		res, err := b.Client.doGet("/fetchMessage", map[string]string{
			"sessionKey": b.SessionKey,
			"count":      strconv.Itoa(b.size),
		})
		if err != nil {
			return err
		}
		var tc []message.Event
		tools.Json.Get([]byte(res), "data").ToVal(&tc)
		for _, v := range tc {
			if len(b.Chan) == b.size {
				<-b.Chan
			}
			b.Chan <- v
		}

		<-t.C
	}
}

// --- 管理相关 ---

// FriendList 使用此方法获取bot的好友列表
func (b *Bot) FriendList() error {
	data := map[string]string{"sessionKey": b.SessionKey}
	res, err := b.Client.doGet("/friendList", data)
	if err != nil {
		return err
	}
	return tools.Json.UnmarshalFromString(res, &b.Friends)
}

// GroupList 使用此方法获取bot的群列表
func (b *Bot) GroupList() error {
	data := map[string]string{"sessionKey": b.SessionKey}
	res, err := b.Client.doGet("/groupList", data)
	if err != nil {
		return err
	}
	return tools.Json.UnmarshalFromString(res, &b.Groups)
}

// MemberList 使用此方法获取bot指定群种的成员列表
func (b *Bot) MemberList(target int64) ([]message.Sender, error) {
	data := map[string]string{"sessionKey": b.SessionKey, "target": strconv.FormatInt(target, 10)}
	res, err := b.Client.doGet("/memberList", data)
	if err != nil {
		return nil, err
	}
	var list []message.Sender
	err = tools.Json.UnmarshalFromString(res, &list)
	return list, err
}

// MuteAll 使用此方法令指定群进行全体禁言（需要有相关限权）
func (b *Bot) MuteAll(target int64) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "target": target}
	_, err := b.Client.doPost("/muteAll", data)
	return err
}

// UnMuteAll 使用此方法令指定群解除全体禁言（需要有相关限权）
func (b *Bot) UnMuteAll(target int64) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "target": target}
	_, err := b.Client.doPost("/unmuteAll", data)
	return err
}

// Mute 使用此方法指定群禁言指定群员（需要有相关限权）
func (b *Bot) Mute(target, memberID, time int64) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "target": target, "memberId": memberID, "time": time}
	_, err := b.Client.doPost("/mute", data)
	return err
}

// UnMute 使用此方法指定群解除群成员禁言（需要有相关限权）
func (b *Bot) UnMute(target, memberID int64) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "target": target, "memberId": memberID}
	_, err := b.Client.doPost("/unmute", data)
	return err
}

// Kick 使用此方法移除指定群成员（需要有相关限权）
func (b *Bot) Kick(target, memberID int64, msg string) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "target": target, "memberId": memberID, "msg": msg}
	_, err := b.Client.doPost("/kick", data)
	return err
}

// Quit 使用此方法使Bot退出群聊
func (b *Bot) Quit(target int64) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "target": target}
	_, err := b.Client.doPost("/quit", data)
	return err
}

// GroupConfig 使用此方法修改群设置（需要有相关限权）
func (b *Bot) GroupConfig(target int64, config message.GroupConfig) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "target": target, "config": config}
	_, err := b.Client.doPost("/groupConfig", data)
	return err
}

// GetGroupConfig 使用此方法获取群设置
func (b *Bot) GetGroupConfig(target int64) (message.GroupConfig, error) {
	r := message.GroupConfig{}
	data := map[string]string{"sessionKey": b.SessionKey, "target": strconv.FormatInt(target, 10)}
	res, err := b.Client.doGet("/groupConfig", data)
	if err != nil {
		return r, err
	}
	err = tools.Json.UnmarshalFromString(res, &r)
	return r, err
}

// MemberInfo 使用此方法修改群员资料（需要有相关限权）
func (b *Bot) MemberInfo(target, memberID int64, info message.MemberInfo) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "target": target, "memberId": memberID, "info": info}
	_, err := b.Client.doPost("/memberInfo", data)
	return err
}

// GetMemberInfo 使用此方法获取群员资料
func (b *Bot) GetMemberInfo(target, memberID int64) (message.MemberInfo, error) {
	r := message.MemberInfo{}
	data := map[string]string{"sessionKey": b.SessionKey, "target": strconv.FormatInt(target, 10), "memberId": strconv.FormatInt(memberID, 10)}
	res, err := b.Client.doGet("/groupConfig", data)
	if err != nil {
		return r, err
	}
	err = tools.Json.UnmarshalFromString(res, &r)
	return r, err
}

// --- 相应 ---

// RespondMemberJoinRequest 响应用户加群请求
// operate	说明
// 0	同意入群
// 1	拒绝入群
// 2	忽略请求
// 3	拒绝入群并添加黑名单，不再接收该用户的入群申请
// 4	忽略入群并添加黑名单，不再接收该用户的入群申请
func (b *Bot) RespondMemberJoinRequest(eventID, fromID, groupID int64, operate int, message string) error {
	data := map[string]interface{}{"sessionKey": b.SessionKey, "eventId": eventID, "fromId": fromID, "groupId": groupID, "operate": operate, "message": message}
	_, err := b.Client.doPost("/resp/memberJoinRequestEvent", data)
	if err != nil {
		return err
	}
	b.Logger.Info("Respond Member Join Request ", fromID, " join ", groupID, " operate: ", operate)
	return nil
}

// --- Handler ---

// UseHandler 使用选定的 EventHandler 进行事件响应
// 未实装
//func (b *Bot) UseHandler(handler helper.EventHandler) {
//	b.handlers = handler
//}

// Run 使用 eventHandler 进行事件相应
// 与使用 FetchMessage 方法有所冲突
// 未实装
func (b *Bot) run() {
	go func() {
		err := b.FetchMessages()
		if err != nil {
			b.Client.Logger.Panic(err)
		}
	}()

	for {
		e := <-b.Chan
		switch e.Type {
		case message.EventReceiveGroupMessage:

		}
	}
}
