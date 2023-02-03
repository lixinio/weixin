package message_api

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiCustomSend   = "/cgi-bin/message/custom/send"
	apiTemplateSend = "/cgi-bin/message/template/send"
)

type MessageApi struct{ *utils.Client }

func NewApi(client *utils.Client) *MessageApi {
	return &MessageApi{Client: client}
}

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

type TemplateMessageMp struct {
	AppID    string `json:"appid"`              // 所需跳转到的小程序appid（该小程序 appid 必须与发模板消息的公众号是绑定关联关系）
	PagePath string `json:"pagepath,omitempty"` // 所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar），要求该小程序已发布
}

type TemplateMessageData struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"` // 模板内容字体颜色，不填默认为黑色
}

type TemplateMessage struct {
	ToUser      string                          `json:"touser,omitempty"`        // 接收者openid
	TemplateID  string                          `json:"template_id"`             // 模板ID
	URL         string                          `json:"url,omitempty"`           // 模板跳转链接（海外帐号没有跳转能力）
	MiniProgram *TemplateMessageMp              `json:"miniprogram,omitempty"`   // 跳小程序所需数据，不需跳小程序可不用传该数据
	ClientMsgID string                          `json:"client_msg_id,omitempty"` // 防重入id。对于同一个openid + client_msg_id, 只发送一条消息,10分钟有效,超过10分钟不保证效果。若无防重入需求，可不填
	Datas       map[string]*TemplateMessageData `json:"data"`                    // 模板数据
}

type TemplateMessageResponse struct {
	utils.WeixinError
	MsgID int64 `json:"msgid"`
}

/*
发送模板消息
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
*/
func (api *MessageApi) SendTemplateMessage(
	ctx context.Context, msg *TemplateMessage,
) (int64, error) {
	resp := &TemplateMessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiTemplateSend, msg, resp); err != nil {
		return 0, err
	}
	return resp.MsgID, nil
}
