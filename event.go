package gomirai

import (
	"encoding/json"
	"reflect"
)

// InEvent 获取到的事件
type InEvent struct {
	Type string `json:"type"`
	QQ   int64  `json:"qq,omitempty"`

	// Event
	AuthorID        int64       `json:"authorId,omitempty"`
	MessageID       int64       `json:"messageId,omitempty"`
	Time            int64       `json:"time,omitempty"`
	Group           Group       `json:"group,omitempty"`
	Origin          interface{} `json:"origin,omitempty"`
	Current         interface{} `json:"current,omitempty"`
	DurationSeconds int64       `json:"durationSeconds,omitempty"`
	IsByBot         bool        `json:"isByBot,omitempty"`
	Member          GroupMember `json:"member,omitempty"`

	//Message
	MessageChain []Message   `json:"message_chain,omitempty"`
	Sender       interface{} `json:"sender,omitempty"`

	Operator       interface{} `json:"operator,omitempty"` // Operator 可能是GroupMember或int64
	OperatorGroup  GroupMember `json:"-"`
	OperatorFriend int64       `json:"-"`
}

// OperatorDetail 获取Operator的详细信息
func (e *InEvent) OperatorDetail() {
	// Operator 为GroupMember
	if reflect.TypeOf(e.Operator) == reflect.TypeOf(make(map[string]interface{})) {
		bytesData, _ := json.Marshal(e.Operator)
		_ = json.Unmarshal(bytesData, &e.OperatorGroup)
	} else {
		e.OperatorFriend = reflect.ValueOf(e.Operator).Int()
	}
}
