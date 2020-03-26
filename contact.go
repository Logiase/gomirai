package gomirai

// Group 群
type Group struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Permission string `json:"permission"`

	// 所属Bot
	Bot *Bot `json:"-"`
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

	Bot *Bot `json:"-"`
}
