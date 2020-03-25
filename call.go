package gomirai

type CommonCall struct {
	SessionKey string `json:"sessionKey"`
}

type MessageCall struct {
	CommonCall
	Target       int64     `json:"target"`
	MessageChain []Message `json:"messageChain"`
}
