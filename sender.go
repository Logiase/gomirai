package gomirai

// Sender 发送者
type Sender struct {
	ID int64	`json:"id,omitempty"`

	MemberName	string	`json:"memberName,omitempty"`
	PerMission	string	`json:"permission,omitempty"`
	Group	Group	`json:"group,omitempty"`
	NickName	string	`json:"nickname,omitempty"`
	Remark	string	`json:"remark,omitempty"`
}