package gomirai

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
