package miniprogram

import (
	"context"
	"github.com/lixinio/weixin/utils"
)

const (
	apiSubscribeMsgSend = "/cgi-bin/message/subscribe/send"
)

type MessageApi struct{ *utils.Client }

func NewApi(client *utils.Client) *MessageApi {
	return &MessageApi{Client: client}
}

// SendSubscribeMessage 发送订阅消息
func (api MessageApi) SendSubscribeMessage(
	ctx context.Context, tmplID string, toUser string, data map[string]DataValue,
) (*utils.WeixinError, error) {
	respMsg := &utils.WeixinError{}
	reqMsg := defaultJsonMessage(tmplID, toUser, data)

	// 修改默认参数通过option进行
	err := api.Client.HTTPPostJson(
		ctx,
		apiSubscribeMsgSend,
		reqMsg,
		respMsg,
	)
	if err != nil {
		return nil, err
	}

	return respMsg, nil
}
