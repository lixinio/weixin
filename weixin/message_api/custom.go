package message_api

// 客服消息

import "context"

const (
	apiCustomSend = "/cgi-bin/message/custom/send"
)

type MessageHeader struct {
	ToUser  string `json:"touser,omitempty"`
	MsgType string `json:"msgtype"`
}

type TextMessage struct {
	*MessageHeader
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

/*
发送客服消息（文本）
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
func (api *MessageApi) SendCustomTextMessage(
	ctx context.Context, openID, content string,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &TextMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "text",
		},
		Text: struct {
			Content string `json:"content"`
		}{
			Content: content,
		},
	}, nil)
}
