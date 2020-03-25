package gomirai

type CommonResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type AuthResponse struct {
	Code    int    `json:"code"`
	Session string `json:"session"`
}

type VerifyResponse struct {
	CommonResponse
}

type MessageResponse struct {
	CommonResponse
	MessageID int64 `json:"messageId"`
}
