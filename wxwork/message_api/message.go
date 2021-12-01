package message_api

// https://github.com/fastwego/wxwork/tree/master/corporation/apis/message
// Copyright 2021 FastWeGo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package message 消息推送

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiSend                  = "/cgi-bin/message/send"
	apiUpdateTaskcard        = "/cgi-bin/message/update_taskcard"
	apiAppchatCreate         = "/cgi-bin/appchat/create"
	apiAppchatUpdate         = "/cgi-bin/appchat/update"
	apiAppchatGet            = "/cgi-bin/appchat/get"
	apiAppchatSend           = "/cgi-bin/appchat/send"
	apiLinkedcorpMessageSend = "/cgi-bin/linkedcorp/message/send"
	apiGetStatistics         = "/cgi-bin/message/get_statistics"
)

type MessageApi struct {
	*utils.Client
	AgentID int
}

func NewApi(client *utils.Client, agentID int) *MessageApi {
	return &MessageApi{Client: client, AgentID: agentID}
}

type MessageResponse struct {
	utils.WeixinError
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
	MsgID        string `json:"msgid"`
	ResponseCode string `json:"response_code"`
}

/*
发送应用消息(文本)
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90236
POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
*/

type TextMessage struct {
	*MessageHeader
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func (api *MessageApi) SendTextMessage(
	ctx context.Context, header *MessageHeader, content string,
) (*MessageResponse, error) {
	result := &MessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiSend, &TextMessage{
		MessageHeader: header,
		AgentID:       api.AgentID,
		MsgType:       "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: content,
		},
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
发送应用消息(文本卡片)
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90236
POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
*/

type TextCardMessage struct {
	*MessageHeader
	MsgType  string `json:"msgtype"`
	AgentID  int    `json:"agentid"`
	TextCard struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		BtnTxt      string `json:"btntxt"`
	} `json:"textcard"`
}

func (api *MessageApi) SendTextCardMessage(
	ctx context.Context, header *MessageHeader,
	title, description, url, btntxt string,
) (*MessageResponse, error) {
	result := &MessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiSend, &TextCardMessage{
		MessageHeader: header,
		AgentID:       api.AgentID,
		MsgType:       "textcard",
		TextCard: struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			BtnTxt      string `json:"btntxt"`
		}{
			Title:       title,
			Description: description,
			URL:         url,
			BtnTxt:      btntxt,
		},
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
发送应用消息(图文)
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90236
POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
*/
type NewsMessageParam struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PicURL      string `json:"picurl"`
	AppID       string `json:"appid"`
	PagePath    string `json:"pagepath"`
}

type NewsMessage struct {
	*MessageHeader
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	News    struct {
		Articles []*NewsMessageParam `json:"articles"`
	} `json:"news"`
}

func (api *MessageApi) SendNewsMessage(
	ctx context.Context, header *MessageHeader,
	articles []*NewsMessageParam,
) (*MessageResponse, error) {
	result := &MessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiSend, &NewsMessage{
		MessageHeader: header,
		AgentID:       api.AgentID,
		MsgType:       "news",
		News: struct {
			Articles []*NewsMessageParam `json:"articles"`
		}{
			Articles: articles,
		},
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
发送应用消息(Markdown)
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90236
POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
*/
type MarkdownMessage struct {
	*MessageHeader
	MsgType  string `json:"msgtype"`
	AgentID  int    `json:"agentid"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

func (api *MessageApi) SendMarkdownMessage(
	ctx context.Context, header *MessageHeader, content string,
) (*MessageResponse, error) {
	result := &MessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiSend, &MarkdownMessage{
		MessageHeader: header,
		AgentID:       api.AgentID,
		MsgType:       "markdown",
		Markdown: struct {
			Content string `json:"content"`
		}{
			Content: content,
		},
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
发送应用消息(image)
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90236
POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
*/
type ImageMessage struct {
	*MessageHeader
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Image   struct {
		MediaID string `json:"media_id"`
	} `json:"image"`
}

func (api *MessageApi) SendImageMessage(
	ctx context.Context, header *MessageHeader, mediaID string,
) (*MessageResponse, error) {
	result := &MessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiSend, &ImageMessage{
		MessageHeader: header,
		AgentID:       api.AgentID,
		MsgType:       "image",
		Image: struct {
			MediaID string `json:"media_id"`
		}{
			MediaID: mediaID,
		},
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
发送应用消息(voice)
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90236
POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
*/
type VoiceMessage struct {
	*MessageHeader
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Voice   struct {
		MediaID string `json:"media_id"`
	} `json:"voice"`
}

func (api *MessageApi) SendVoiceMessage(
	ctx context.Context, header *MessageHeader, mediaID string,
) (*MessageResponse, error) {
	result := &MessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiSend, &VoiceMessage{
		MessageHeader: header,
		AgentID:       api.AgentID,
		MsgType:       "voice",
		Voice: struct {
			MediaID string `json:"media_id"`
		}{
			MediaID: mediaID,
		},
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
发送应用消息(video)
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90236
POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
*/
type VideoMessage struct {
	*MessageHeader
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Video   struct {
		MediaID string `json:"media_id"`
	} `json:"video"`
}

func (api *MessageApi) SendVideoMessage(
	ctx context.Context, header *MessageHeader, mediaID string,
) (*MessageResponse, error) {
	result := &MessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiSend, &VideoMessage{
		MessageHeader: header,
		AgentID:       api.AgentID,
		MsgType:       "video",
		Video: struct {
			MediaID string `json:"media_id"`
		}{
			MediaID: mediaID,
		},
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
发送应用消息(file)
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90236
POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
*/
type FileMessage struct {
	*MessageHeader
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	File    struct {
		MediaID string `json:"media_id"`
	} `json:"file"`
}

func (api *MessageApi) SendFileMessage(
	ctx context.Context, header *MessageHeader, mediaID string,
) (*MessageResponse, error) {
	result := &MessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiSend, &FileMessage{
		MessageHeader: header,
		AgentID:       api.AgentID,
		MsgType:       "file",
		File: struct {
			MediaID string `json:"media_id"`
		}{
			MediaID: mediaID,
		},
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
发送应用消息(图文)
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90236
POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
*/
type MpNewsMessageParam struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	Author           string `json:"author"`
	ContentSourceURL string `json:"content_source_url"`
	Content          string `json:"content"`
	Digest           string `json:"digest"`
}

type MpNewsMessage struct {
	*MessageHeader
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	MpNews  struct {
		Articles []*MpNewsMessageParam `json:"articles"`
	} `json:"mpnews"`
}

func (api *MessageApi) SendMpNewsMessage(
	ctx context.Context, header *MessageHeader,
	articles []*MpNewsMessageParam,
) (*MessageResponse, error) {
	result := &MessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiSend, &MpNewsMessage{
		MessageHeader: header,
		AgentID:       api.AgentID,
		MsgType:       "mpnews",
		MpNews: struct {
			Articles []*MpNewsMessageParam `json:"articles"`
		}{
			Articles: articles,
		},
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
发送应用消息(小程序)
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90236
POST https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN
*/
type MpNoticeItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MpNoticeMessageParam struct {
	AppID             string          `json:"appid"`
	Page              string          `json:"page"`
	Title             string          `json:"title"`
	Description       string          `json:"description"`
	EmphasisFirstItem bool            `json:"emphasis_first_item"`
	ContentItem       []*MpNoticeItem `json:"content_item"`
}

type MpNoticeMessage struct {
	*MessageHeader
	MsgType  string                `json:"msgtype"`
	AgentID  int                   `json:"agentid"`
	MpNotice *MpNoticeMessageParam `json:"miniprogram_notice"`
}

func (api *MessageApi) SendMpNoticeMessage(
	ctx context.Context, header *MessageHeader,
	msg *MpNoticeMessageParam,
) (*MessageResponse, error) {
	result := &MessageResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiSend, &MpNoticeMessage{
		MessageHeader: header,
		AgentID:       api.AgentID,
		MsgType:       "miniprogram_notice",
		MpNotice:      msg,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
更新任务卡片消息状态
应用可以发送任务卡片消息，发送之后可再通过接口更新用户任务卡片消息的选择状态。
See: https://work.weixin.qq.com/api/doc/90000/90135/91579
POST https://qyapi.weixin.qq.com/cgi-bin/message/update_taskcard?access_token=ACCESS_TOKEN
*/
// func (api *MessageApi) UpdateTaskcard(
// 	ctx context.Context,
// 	payload []byte,
// ) (resp []byte, err error) {
// 	return api.Client.HTTPPost(
// 		ctx,
// 		apiUpdateTaskcard,
// 		bytes.NewReader(payload),
// 		"application/json;charset=utf-8",
// 	)
// }

/*
创建群聊会话
See: https://work.weixin.qq.com/api/doc/90000/90135/90245
POST https://qyapi.weixin.qq.com/cgi-bin/appchat/create?access_token=ACCESS_TOKEN
*/
// func (api *MessageApi) AppchatCreate(ctx context.Context, payload []byte) (resp []byte, err error) {
// 	return api.Client.HTTPPost(
// 		ctx,
// 		apiAppchatCreate,
// 		bytes.NewReader(payload),
// 		"application/json;charset=utf-8",
// 	)
// }

/*
修改群聊会话
See: https://work.weixin.qq.com/api/doc/90000/90135/90246
POST https://qyapi.weixin.qq.com/cgi-bin/appchat/update?access_token=ACCESS_TOKEN
*/
// func (api *MessageApi) AppchatUpdate(ctx context.Context, payload []byte) (resp []byte, err error) {
// 	return api.Client.HTTPPost(
// 		ctx,
// 		apiAppchatUpdate,
// 		bytes.NewReader(payload),
// 		"application/json;charset=utf-8",
// 	)
// }

/*
获取群聊会话
See: https://work.weixin.qq.com/api/doc/90000/90135/90247
GET https://qyapi.weixin.qq.com/cgi-bin/appchat/get?access_token=ACCESS_TOKEN&chatid=CHATID
*/
// func (api *MessageApi) AppchatGet(ctx context.Context, params url.Values) (resp []byte, err error) {
// 	return api.Client.HTTPGet(ctx, apiAppchatGet+"?"+params.Encode())
// }

/*
应用推送消息
应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90248
POST https://qyapi.weixin.qq.com/cgi-bin/appchat/send?access_token=ACCESS_TOKEN
*/
// func (api *MessageApi) AppchatSend(ctx context.Context, payload []byte) (resp []byte, err error) {
// 	return api.Client.HTTPPost(
// 		ctx,
// 		apiAppchatSend,
// 		bytes.NewReader(payload),
// 		"application/json;charset=utf-8",
// 	)
// }

/*
互联企业消息推送
互联企业的应用支持推送文本、图片、视频、文件、图文等类型。
See: https://work.weixin.qq.com/api/doc/90000/90135/90250
POST https://qyapi.weixin.qq.com/cgi-bin/linkedcorp/message/send?access_token=ACCESS_TOKEN
*/
// func (api *MessageApi) LinkedcorpMessageSend(
// 	ctx context.Context,
// 	payload []byte,
// ) (resp []byte, err error) {
// 	return api.Client.HTTPPost(
// 		ctx,
// 		apiLinkedcorpMessageSend,
// 		bytes.NewReader(payload),
// 		"application/json;charset=utf-8",
// 	)
// }

/*
查询应用消息发送统计
See: https://work.weixin.qq.com/api/doc/90000/90135/92369
POST https://qyapi.weixin.qq.com/cgi-bin/message/get_statistics?access_token=ACCESS_TOKEN
*/
// func (api *MessageApi) GetStatistics(ctx context.Context, payload []byte) (resp []byte, err error) {
// 	return api.Client.HTTPPost(
// 		ctx,
// 		apiGetStatistics,
// 		bytes.NewReader(payload),
// 		"application/json;charset=utf-8",
// 	)
// }
