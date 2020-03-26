package gomirai

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`

	MessageID int64 `json:"messageId,omitempty"`
}

type AuthResponse struct {
	Code    int    `json:"code"`
	Session string `json:"session"`
}
