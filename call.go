package gomirai

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
