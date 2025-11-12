package message_api

// 群发消息

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiSendMassMessageByTag    = "/cgi-bin/message/mass/sendall" // 根据标签群发消息
	apiSendMassMessageByOpenid = "/cgi-bin/message/mass/send"    // 根据OpenID群发消息
	apiPreviewMassMessage      = "/cgi-bin/message/mass/preview" // 预览消息
)

type SendMassMessageResult struct {
	utils.WeixinError
	MsgID     int64 `json:"msg_id"`
	MsgDataID int64 `json:"msg_data_id"`
}

func buildSendMassMessageReq(
	sendAll bool, tagID int64, clientmsgid, key string, content any,
	keys ...string,
) any {
	v := map[string]any{
		"msgtype":     key,
		"clientmsgid": clientmsgid,
	}

	if sendAll {
		v["filter"] = map[string]bool{
			"is_to_all": true,
		}
	} else {
		v["filter"] = map[string]any{
			"is_to_all": false,
			"tag_id":    tagID,
		}
	}

	key2 := key
	if len(keys) > 0 {
		key2 = keys[0]
	}

	v[key2] = content

	return v
}

/*
根据标签群发消息
https://developers.weixin.qq.com/doc/service/api/notify/message/api_sendall.html
*/
func (api *MessageApi) SendMassTextMessage(
	ctx context.Context,
	sendAll bool, tagID int64, clientmsgid string,
	content string,
) (int64, int64, error) {
	resp := &SendMassMessageResult{}
	if err := api.Client.HTTPPostJson(
		ctx,
		apiSendMassMessageByTag,
		buildSendMassMessageReq(sendAll, tagID, clientmsgid, "text", map[string]any{
			"content": content,
		}),
		resp,
	); err != nil {
		return 0, 0, err
	}

	return resp.MsgID, resp.MsgDataID, nil
}

func (api *MessageApi) SendMassVoiceMessage(
	ctx context.Context,
	sendAll bool, tagID int64, clientmsgid string,
	mediaID string,
) (int64, int64, error) {
	resp := &SendMassMessageResult{}
	if err := api.Client.HTTPPostJson(
		ctx,
		apiSendMassMessageByTag,
		buildSendMassMessageReq(sendAll, tagID, clientmsgid, "voice", map[string]any{
			"media_id": mediaID,
		}),
		resp,
	); err != nil {
		return 0, 0, err
	}

	return resp.MsgID, resp.MsgDataID, nil
}

type MassImageRequest struct {
	MediaIds           []string `json:"media_ids"`             // 用于群发的图文消息的media_id列表
	Recommend          string   `json:"recommend"`             // 推荐语
	Title              string   `json:"title"`                 // 标题
	NeedOpenComment    int      `json:"need_open_comment"`     // 开启评论（1-开启，0-关闭）
	OnlyFansCanComment int      `json:"only_fans_can_comment"` // 只有粉丝能评论（1-开启，0-关闭）
}

func (api *MessageApi) SendMassImageMessage(
	ctx context.Context,
	sendAll bool, tagID int64, clientmsgid string,
	req *MassImageRequest,
) (int64, int64, error) {
	resp := &SendMassMessageResult{}
	if err := api.Client.HTTPPostJson(
		ctx,
		apiSendMassMessageByTag,
		buildSendMassMessageReq(sendAll, tagID, clientmsgid, "image", req, "images"),
		resp,
	); err != nil {
		return 0, 0, err
	}

	return resp.MsgID, resp.MsgDataID, nil
}

func (api *MessageApi) SendMassVideoMessage(
	ctx context.Context,
	sendAll bool, tagID int64, clientmsgid string,
	mediaID string,
) (int64, int64, error) {
	resp := &SendMassMessageResult{}
	if err := api.Client.HTTPPostJson(
		ctx,
		apiSendMassMessageByTag,
		buildSendMassMessageReq(sendAll, tagID, clientmsgid, "mpvideo", map[string]any{
			"media_id": mediaID,
		}),
		resp,
	); err != nil {
		return 0, 0, err
	}

	return resp.MsgID, resp.MsgDataID, nil
}

func (api *MessageApi) SendMassCardMessage(
	ctx context.Context,
	sendAll bool, tagID int64, clientmsgid string,
	cardID string,
) (int64, int64, error) {
	resp := &SendMassMessageResult{}
	if err := api.Client.HTTPPostJson(
		ctx,
		apiSendMassMessageByTag,
		buildSendMassMessageReq(sendAll, tagID, clientmsgid, "wxcard", map[string]any{
			"card_id": cardID,
		}),
		resp,
	); err != nil {
		return 0, 0, err
	}

	return resp.MsgID, resp.MsgDataID, nil
}

func (api *MessageApi) SendMassNewsMessage(
	ctx context.Context,
	sendAll bool, tagID int64, clientmsgid string,
	mediaID string,
) (int64, int64, error) {
	resp := &SendMassMessageResult{}
	if err := api.Client.HTTPPostJson(
		ctx,
		apiSendMassMessageByTag,

		buildSendMassMessageReq(sendAll, tagID, clientmsgid, "mpnews", map[string]any{
			"media_id": mediaID,
		}),
		resp,
	); err != nil {
		return 0, 0, err
	}

	return resp.MsgID, resp.MsgDataID, nil
}

// 预览消息
// https://developers.weixin.qq.com/doc/service/api/notify/message/api_preview.html
func buildPreviewMassMessageReq(
	touser, towxname, key string, content any,
) any {
	v := map[string]any{
		"msgtype": key,
	}

	if touser != "" {
		v["touser"] = touser
	} else {
		v["towxname"] = towxname
	}

	v[key] = content

	return v
}

func (api *MessageApi) PreviewMassTextMessage(
	ctx context.Context,
	touser, towxname string,
	content string,
) error {
	if err := api.Client.HTTPPostJson(
		ctx,
		apiPreviewMassMessage,
		buildPreviewMassMessageReq(touser, towxname, "text", map[string]any{
			"content": content,
		}),
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (api *MessageApi) PreviewMassVoiceMessage(
	ctx context.Context,
	touser, towxname string,
	mediaID string,
) error {
	if err := api.Client.HTTPPostJson(
		ctx,
		apiPreviewMassMessage,
		buildPreviewMassMessageReq(touser, towxname, "voice", map[string]any{
			"media_id": mediaID,
		}),
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (api *MessageApi) PreviewMassImageMessage(
	ctx context.Context,
	touser, towxname string,
	mediaID string,
) error {
	if err := api.Client.HTTPPostJson(
		ctx,
		apiPreviewMassMessage,
		buildPreviewMassMessageReq(touser, towxname, "image", map[string]any{
			"media_id": mediaID,
		}),
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (api *MessageApi) PreviewMassVideoMessage(
	ctx context.Context,
	touser, towxname string,
	mediaID string,
) error {
	if err := api.Client.HTTPPostJson(
		ctx,
		apiPreviewMassMessage,
		buildPreviewMassMessageReq(touser, towxname, "mpvideo", map[string]any{
			"media_id": mediaID,
		}),
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (api *MessageApi) PreviewMassCardMessage(
	ctx context.Context,
	touser, towxname string,
	cardID string,
) error {
	if err := api.Client.HTTPPostJson(
		ctx,
		apiPreviewMassMessage,
		buildPreviewMassMessageReq(touser, towxname, "wxcard", map[string]any{
			"card_id": cardID,
		}),
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (api *MessageApi) PreviewMassNewsMessage(
	ctx context.Context,
	touser, towxname string,
	mediaID string,
) error {
	if err := api.Client.HTTPPostJson(
		ctx,
		apiPreviewMassMessage,

		buildPreviewMassMessageReq(touser, towxname, "mpnews", map[string]any{
			"media_id": mediaID,
		}),
		nil,
	); err != nil {
		return err
	}

	return nil
}
