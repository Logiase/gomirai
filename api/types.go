// Types passed and returned to and from the API

package api

import (
	"encoding/json"
	"strconv"
)

// <-- Contact -->
// Group 群
type Group struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

// GroupMember 群成员
type GroupMember struct {
	ID         int64  `json:"id"`
	MemberName string `json:"memberName"`
	Permission string `json:"permission"`
	Group      Group  `json:"group"`
}

// Friend 好友
type Friend struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickName"`
	Remark   string `json:"remark"`
}

// GroupConfig 群设置
type GroupConfig struct {
	Name              string `json:"name"`
	Announcement      string `json:"announcement"`
	ConfessTalk       bool   `json:"confessTalk"`
	AllowMemberInvite bool   `json:"allowMemberInvite"`
	AutoApprove       bool   `json:"autoApprove"`
	AnonymousChat     bool   `json:"anonymousChat"`
}

// GroupMemberInfo 群成员信息
type GroupMemberInfo struct {
	Name         string `json:"name"`
	SpecialTitle string `json:"specialTitle"`
}

// <-- Message -->
// Message 消息
type Message struct {
	Type string `json:"type"`

	ID       int64     `json:"id,omitempty" message:"Source|Quote"`
	Text     string    `json:"text,omitempty" message:"Plain"`
	Time     int64     `json:"time,omitempty" message:"Source"`
	GroupID  int64     `json:"groupId,omitempty" message:"Quote"`
	SenderID int64     `json:"senderId,omitempty" message:"Quote"`
	Origin   []Message `json:"origin,omitempty" message:"Quote"`
	Target   int64     `json:"target,omitempty" message:"At"`
	Display  string    `json:"display,omitempty" message:"At"`
	FaceID   int64     `json:"faceId,omitempty" message:"Face"`
	Name     string    `json:"name,omitempty" message:"Face"`
	ImageID  string    `json:"imageId,omitempty" message:"Image"`
	URL      string    `json:"url,omitempty" message:"Image"`
	Path     string    `json:"path,omitempty" message:"Image"`
	XML      string    `json:"xml,omitempty" message:"Xml"`
	JSON     string    `json:"json,omitempty" message:"Json"`
	Content  string    `json:"content,omitempty" message:"App"`
}

// <-- Response -->
// Response 通用响应
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`

	MessageID int64 `json:"messageId,omitempty"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Code    int    `json:"code"`
	Session string `json:"session"`
}

// SessionConfig Session设置
type SessionConfig struct {
	SessionKey      string `json:"sessionKey,omitempty"`
	CacheSize       int64  `json:"cacheSize"`
	EnableWebsocket bool   `json:"enableWebsocket"`
}

// <-- Event -->
type Event struct {
	Type string `json:"type"`
	QQ   int64  `json:"qq,omitempty"`

	// Message
	MessageChain []Message   `json:"messageChain,omitempty"`
	Sender       GroupMember `json:"sender,omitempty"`

	// Event
	IsByBot         bool               `json:"isByBot,omitempty"`
	AuthorID        int64              `json:"authorId,omitempty"`
	MessageID       int64              `json:"messageId,omitempty"`
	Time            int64              `json:"time,omitempty"`
	Group           Group              `json:"group,omitempty"`
	Origin          interface{}        `json:"origin,omitempty"`
	Current         interface{}        `json:"current,omitempty"`
	DurationSeconds int64              `json:"durationSeconds,omitempty"`
	Member          GroupMemberWrapper `json:"member,omitempty"`
	Operator        GroupMemberWrapper `json:"operator,omitempty"` // Operator 可能是GroupMember或int64
}

type GroupMemberWrapper struct {
	GroupMember
}

// https://medium.com/@nate510/dynamic-json-umarshalling-in-go-88095561d6a0
// In order to correctly unmarshal either an ID or an object, we can now override the default behavior:
func (m *GroupMemberWrapper) UnmarshalJSON(data []byte) error {
	if id, e := strconv.ParseInt(string(data), 10, 64); e == nil {
		m.ID = id
		return nil
	}
	return json.Unmarshal(data, &m.GroupMember)
}

// <-- Call -->
// CommonCall 通用请求
type CommonCall struct {
	SessionKey string `json:"sessionKey"`
}

// MessageCall 消息用
type MessageCall struct {
	CommonCall
	Target       int64     `json:"target"`
	MessageChain []Message `json:"messageChain"`
}
