package gomirai

const (
	MsgType_Source     = "Source"
	MsgType_Quote      = "Quote"
	MsgType_At         = "At"
	MsgType_AtAll      = "AtAll"
	MsgType_Face       = "Face"
	MsgType_Plain      = "Plain"
	MsgType_Image      = "Image"
	MsgType_FlashImage = "FlashImage"
	MsgType_Xml        = "Xml"
	MsgType_Json       = "Json"
	MsgType_App        = "App"
	MsgType_Poke       = "Poke"
)

type MessageChain struct {
	Msg []Message
}

type Message struct {
	Type      string    `json:"type,omitempty"`
	Id        uint      `json:"id,omitempty"`       //(Source,Quote)Source中表示消息id，Quote中表示被引用回复的原消息的id
	Time      int64     `json:"time,omitempty"`     //(Source) 发送时间
	GroupId   uint      `json:"groupId,omitempty"`  //(Quote)Quote中表示被引用回复的原消息的群号
	SenderId  uint      `json:"senderId,omitempty"` //(Quote)Quote中表示被引用回复的原消息的发送者QQ号
	TargetId  uint      `json:"targetId,omitempty"` //(Quote)Quote中表示被引用回复的原消息的接收者群号或QQ号
	Origin    []Message `json:"origin,omitempty"`   //(Quote)Quote中表示被引用回复的原消息的消息链对象
	Target    uint      `json:"target,omitempty"`   //(At)@的群员QQ号
	Display   string    `json:"display,omitempty"`  //(At)@的显示文本
	FaceId    uint      `json:"faceId,omitempty"`   //(Face)QQ表情的ID,发送时优先级比Name高
	Name      string    `json:"name,omitempty"`     //(Face,Poke)Face中为QQ表情的拼音,Poke中为戳一戳的类型
	Text      string    `json:"text,omitempty"`     //(Plain)纯文本
	ImageId   string    `json:"imageId,omitempty"`  //(Image,FlashImage)图片ID，注意消息类型，群图片和好友图片格式不一样，发送时优先级比ImageUrl高
	ImageUrl  string    `json:"url,omitempty"`      //(Image,FlashImage)图片url,发送时可使用网络图片的链接，优先级比ImagePath高；接收时为腾讯图片服务器的链接
	ImagePath string    `json:"path,omitempty"`     //(Image,FlashImage)图片的路径，发送本地图片，相对路径于plugins/MiraiAPIHTTP/images
	Xml       string    `json:"xml,omitempty"`      //(Xml) xml消息本体
	Json      string    `json:"json,omitempty"`     //(Json) json消息本体
	Content   string    `json:"content,omitempty"`  //(App) 不知道干嘛的，mirai也没有说明，估计是小程序连接？
}
