package api

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Message&Event response
var (
	fetchMessage = `[
	{
		"type":"GroupMessage",
		"messageChain":[
			{
				"type":"Source",
				"id":123456,
				"time":123456789
			},
			{
				"type":"Plain",
				"text":"Miral牛逼"
			}
		],
		"sender":{
			"id":123456789,
			"memberName":"化腾",
			"permission":"MEMBER",
			"group":{
				"id":1234567890,
				"name":"Miral Technology",
				"permission":"MEMBER"
			}
		}
	},
	{
		"type":"FriendMessage",
		"messageChain":[
			{
				"type":"Source",
				"id":123456,
				"time":123456789
			},
			{
				"type":"Plain",
				"text":"Miral牛逼"
			}
		],
		"sender":{
			"id":1234567890,
			"nickname":"",
			"remark":""
		}
	},
	{
		"type":"MemberMuteEvent",
		"durationSeconds":600,
		"member":{
			"id":123456789,
			"memberName":"禁言对象",
			"permission":"MEMBER",
			"group":{
				"id":123456789,
				"name":"Miral Technology",
				"permission":"MEMBER"
			}
		},
		"operator":{
			"id":987654321,
			"memberName":"群主大人",
			"permission":"OWNER",
			"group":{
				"id":123456789,
				"name":"Miral Technology",
				"permission":"MEMBER"
			}
		}
	}
]`
	fullResult       = `[]api.Event{api.Event{Type:"GroupMessage", QQ:0, MessageChain:[]api.Message{api.Message{Type:"Source", ID:123456, Text:"", Time:123456789, GroupID:0, SenderID:0, Origin:[]api.Message(nil), Target:0, Display:"", FaceID:0, Name:"", ImageID:"", URL:"", Path:"", XML:"", JSON:"", Content:""}, api.Message{Type:"Plain", ID:0, Text:"Miral牛逼", Time:0, GroupID:0, SenderID:0, Origin:[]api.Message(nil), Target:0, Display:"", FaceID:0, Name:"", ImageID:"", URL:"", Path:"", XML:"", JSON:"", Content:""}}, Sender:api.GroupMember{ID:123456789, MemberName:"化腾", Permission:"MEMBER", Group:api.Group{ID:1234567890, Name:"Miral Technology", Permission:"MEMBER"}}, IsByBot:false, AuthorID:0, MessageID:0, Time:0, Group:api.Group{ID:0, Name:"", Permission:""}, Origin:interface {}(nil), Current:interface {}(nil), DurationSeconds:0, Member:api.GroupMemberWrapper{GroupMember:api.GroupMember{ID:0, MemberName:"", Permission:"", Group:api.Group{ID:0, Name:"", Permission:""}}}, Operator:api.GroupMemberWrapper{GroupMember:api.GroupMember{ID:0, MemberName:"", Permission:"", Group:api.Group{ID:0, Name:"", Permission:""}}}}, api.Event{Type:"FriendMessage", QQ:0, MessageChain:[]api.Message{api.Message{Type:"Source", ID:123456, Text:"", Time:123456789, GroupID:0, SenderID:0, Origin:[]api.Message(nil), Target:0, Display:"", FaceID:0, Name:"", ImageID:"", URL:"", Path:"", XML:"", JSON:"", Content:""}, api.Message{Type:"Plain", ID:0, Text:"Miral牛逼", Time:0, GroupID:0, SenderID:0, Origin:[]api.Message(nil), Target:0, Display:"", FaceID:0, Name:"", ImageID:"", URL:"", Path:"", XML:"", JSON:"", Content:""}}, Sender:api.GroupMember{ID:1234567890, MemberName:"", Permission:"", Group:api.Group{ID:0, Name:"", Permission:""}}, IsByBot:false, AuthorID:0, MessageID:0, Time:0, Group:api.Group{ID:0, Name:"", Permission:""}, Origin:interface {}(nil), Current:interface {}(nil), DurationSeconds:0, Member:api.GroupMemberWrapper{GroupMember:api.GroupMember{ID:0, MemberName:"", Permission:"", Group:api.Group{ID:0, Name:"", Permission:""}}}, Operator:api.GroupMemberWrapper{GroupMember:api.GroupMember{ID:0, MemberName:"", Permission:"", Group:api.Group{ID:0, Name:"", Permission:""}}}}, api.Event{Type:"MemberMuteEvent", QQ:0, MessageChain:[]api.Message(nil), Sender:api.GroupMember{ID:0, MemberName:"", Permission:"", Group:api.Group{ID:0, Name:"", Permission:""}}, IsByBot:false, AuthorID:0, MessageID:0, Time:0, Group:api.Group{ID:0, Name:"", Permission:""}, Origin:interface {}(nil), Current:interface {}(nil), DurationSeconds:600, Member:api.GroupMemberWrapper{GroupMember:api.GroupMember{ID:123456789, MemberName:"禁言对象", Permission:"MEMBER", Group:api.Group{ID:123456789, Name:"Miral Technology", Permission:"MEMBER"}}}, Operator:api.GroupMemberWrapper{GroupMember:api.GroupMember{ID:987654321, MemberName:"群主大人", Permission:"OWNER", Group:api.Group{ID:123456789, Name:"Miral Technology", Permission:"MEMBER"}}}}}`
	fetchPartMessage = `[
	{
		"type":"MemberMuteEvent",
		"durationSeconds":600,
		"member":{
			"id":123456789,
			"memberName":"禁言对象",
			"permission":"MEMBER",
			"group":{
				"id":123456789,
				"name":"Miral Technology",
				"permission":"MEMBER"
			}
		},
		"operator": 987654321
	}
]`
	partResult = `[]api.Event{api.Event{Type:"MemberMuteEvent", QQ:0, MessageChain:[]api.Message(nil), Sender:api.GroupMember{ID:0, MemberName:"", Permission:"", Group:api.Group{ID:0, Name:"", Permission:""}}, IsByBot:false, AuthorID:0, MessageID:0, Time:0, Group:api.Group{ID:0, Name:"", Permission:""}, Origin:interface {}(nil), Current:interface {}(nil), DurationSeconds:600, Member:api.GroupMemberWrapper{GroupMember:api.GroupMember{ID:123456789, MemberName:"禁言对象", Permission:"MEMBER", Group:api.Group{ID:123456789, Name:"Miral Technology", Permission:"MEMBER"}}}, Operator:api.GroupMemberWrapper{GroupMember:api.GroupMember{ID:987654321, MemberName:"", Permission:"", Group:api.Group{ID:0, Name:"", Permission:""}}}}}`
)

func TestUnmarshalMessageEvent(t *testing.T) {
	var i []Event
	if e := json.Unmarshal([]byte(fetchMessage), &i); e == nil {
		if fmt.Sprintf("%#v", i) != fullResult {
			fmt.Printf("%#v", i)
			t.Fail()
		}
	} else {
		fmt.Println(e)
		t.Fail()
	}
}

func TestUnmarshalPartMessageEvent(t *testing.T) {
	var i []Event
	if e := json.Unmarshal([]byte(fetchPartMessage), &i); e == nil {
		if fmt.Sprintf("%#v", i) != partResult {
			fmt.Printf("%#v", i)
			t.Fail()
		}
	} else {
		fmt.Println(e)
		t.Fail()
	}
}
