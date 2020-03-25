package gomirai

// Message 消息
type Message struct {
	Type string `json:"type"`

	ID       int64     `json:"id,omitempty" message:"Source|Quote"`
	Text     string    `json:"text,omitempty" message:"Plain"`
	Time     int64     `json:"time,omitempty" message:"Source"`
	GroupID  int64     `json:"groupId,omitempty" message:"Quote"`
	SenderID int64     `json:"senderId,omitempty" message:"Quote"`
	Origin   []Message `json:"origin,omitempty" message:"Quote"`
	Target   int64     `json:"target,omitempty" message:"At"`
	Display  string    `json:"display,omitempty" message:"At"`
	FaceID   int64     `json:"faceId,omitempty" message:"Face"`
	Name     string    `json:"name,omitempty" message:"Face"`
	ImageID  string    `json:"imageId,omitempty" message:"Image"`
	URL      string    `json:"url,omitempty" message:"Image"`
	Path     string    `json:"path,omitempty" message:"Image"`
	XML      string    `json:"xml,omitempty" message:"Xml"`
	JSON     string    `json:"json,omitempty" message:"Json"`
	Content  string    `json:"content,omitempty" message:"App"`
}
