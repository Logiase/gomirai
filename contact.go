package gomirai

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

// GroupConfig
type GroupConfig struct {
	Name              string `json:"name"`
	Announcement      string `json:"announcement"`
	ConfessTalk       bool   `json:"confessTalk"`
	AllowMemberInvite bool   `json:"allowMemberInvite"`
	AutoApprove       bool   `json:"autoApprove"`
	AnonymousChat     bool   `json:"anonymousChat"`
}

// GroupMemberInfo
type GroupMemberInfo struct {
	Name         string `json:"name"`
	SpecialTitle string `json:"specialTitle"`
}
