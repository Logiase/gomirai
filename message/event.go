package message

const (
	EventReceiveFriendMessage            = "FriendMessage"
	EventReceiveGroupMessage             = "GroupMessage"
	EventReceiveTempMessage              = "TempMessage"
	EventBotOnline                       = "BotOnlineEvent"
	EventBotOfflineActive                = "BotOfflineEventActive"
	EventBotOfflineForce                 = "BotOfflineEventForce"
	EventBotOfflineDropped               = "BotOfflineEventDropped"
	EventBotRelogin                      = "BotReloginEvent"
	EventGroupRecall                     = "GroupRecallEvent"
	EventFriendRecall                    = "FriendRecallEvent"
	EventBotGroupPermissionChange        = "BotGroupPermissionChangeEvent"
	EventBotMute                         = "BotMuteEvent"
	EventBotUnmute                       = "BotUnmuteEvent"
	EventBotJoinGroup                    = "BotJoinGroupEvent"
	EventBotLeaveActive                  = "BotLeaveEventActive"
	EventBotLeaveKick                    = "BotLeaveEventKick"
	EventGroupNameChange                 = "GroupNameChangeEvent"
	EventGroupEntranceAnnouncementChange = "GroupEntranceAnnouncementChangeEvent"
	EventGroupMuteAll                    = "GroupMuteAllEvent"
	EventGroupAllowAnonymousChat         = "GroupAllowAnonymousChatEvent"
	EventGroupAllowConfessTalk           = "GroupAllowConfessTalkEvent"
	EventGroupAllowMemberInvite          = "GroupAllowMemberInviteEvent"
	EventMemberJoin                      = "MemberJoinEvent"
	EventMemberLeaveKick                 = "MemberLeaveEventKick"
	EventMemberLeaveQuit                 = "MemberLeaveEventQuit"
	EventMemberCardChange                = "MemberCardChangeEvent"
	EventMemberSpecialTitleChange        = "MemberSpecialTitleChangeEvent"
	EventMemberPermissionChange          = "MemberPermissionChangeEvent"
	EventMemberMute                      = "MemberMuteEvent"
	EventMemberUnmute                    = "MemberUnmuteEvent"
	EventNewFriendRequest                = "NewFriendRequestEvent"
	EventMemberJoinRequest               = "MemberJoinRequestEvent"
	EventBotInvitedJoinGroupRequest      = "BotInvitedJoinGroupRequestEvent"
)

type Event struct {
	Type         string    `json:"type"`         //事件类型
	MessageChain []Message `json:"messageChain"` //(ReceiveMessage)消息链
	Sender       Sender    `json:"sender"`       //(ReceiveMessage)发送者信息
	EventId      uint      `json:"eventId"`      //事件ID
	FromId       uint      `json:"fromId"`       //操作人
	GroupId      uint      `json:"groupId"`      //群号
}

type Group struct {
	Id        uint   `json:"id,omitempty"`        //消息来源群号
	Name      string `json:"name,omitempty"`      //消息来源群名
	Permisson string `json:"permisson,omitempty"` //bot在群中的角色
}

type Friend struct {
	Id       uint   `json:"id,omitempty"`         //QQ
	NickName string `json:"memberName,omitempty"` //昵称
	Remark   string `json:"remark,omitempty"`     //备注
}

type Sender struct {
	Friend
	MemberName string `json:"memberName,omitempty"` //(GroupMessage)发送者群昵称
	Permission string `json:"permission,omitempty"` //(GroupMessage)发送者在群中的角色
	Group      Group  `json:"group,omitempty"`      //(GroupMessage)消息来源群信息
}

type GroupConfig struct {
	Name              string
	Announcement      string
	ConfessTalk       bool
	AllowMemberInvite bool
	AutoApprove       bool
	AnonymousChat     bool
}

type MemberInfo struct {
	Name         string
	SpecialTitle string
}
