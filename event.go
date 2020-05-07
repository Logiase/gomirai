package gomirai

import (
	"encoding/json"
	"errors"
	"reflect"
)

// InEvent http-api新返回格式
type InEventAll struct {
	Code         int64     `json:"code"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
	Data         []InEvent `json:"data,omitempty"`
}

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
	MessageChain []Message `json:"messageChain,omitempty"`

	Sender       interface{} `json:"sender,omitempty"`
	SenderGroup  GroupMember `json:"-"`
	SenderFriend Friend      `json:"-"`

	Operator       interface{} `json:"operator,omitempty"` // Operator 可能是GroupMember或int64
	OperatorGroup  GroupMember `json:"-"`
	OperatorFriend int64       `json:"-"`
}

// OperatorDetail 获取Operator的详细信息
func (e *InEvent) OperatorDetail() error {
	// Operator 为GroupMember
	if reflect.TypeOf(e.Operator) == reflect.TypeOf(make(map[string]interface{})) {
		bytesData, err := json.Marshal(e.Operator)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytesData, &e.OperatorGroup)
		if err != nil {
			return err
		}
	} else {
		e.OperatorFriend = reflect.ValueOf(e.Operator).Int()
	}
	return nil
}

// SenderDetail 获取Sender的详细信息
func (e *InEvent) SenderDetail() error {
	if reflect.TypeOf(e.Sender).Kind() == reflect.Map {
		keys := reflect.ValueOf(e.Sender).MapKeys()
		for _, k := range keys {
			if k.String() == "memberName" {
				bytesData, err := json.Marshal(e.Sender)
				if err != nil {
					return err
				}
				err = json.Unmarshal(bytesData, &e.SenderGroup)
				if err != nil {
					return err
				}
				return nil
			} else if k.String() == "nickname" {
				bytesData, err := json.Marshal(e.Sender)
				if err != nil {
					return err
				}
				err = json.Unmarshal(bytesData, &e.SenderFriend)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	return errors.New("sender 序列化失败")
}
