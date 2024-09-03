package externalcontact_api

import (
	"context"
	"github.com/lixinio/weixin/utils"
)

const (
	apiAddMsgTemplate = "/cgi-bin/externalcontact/add_msg_template" // 创建企业群发
)

type ChatType string

const (
	ChatTypeSingle ChatType = "single" // 发送给客户
	ChatTypeGroup  ChatType = "group"  // 发送给客户群
)

type AddMsgTemplateRequest struct {
	ChatType       string       `json:"chat_type"`
	ExternalUserID []string     `json:"external_userid"`
	ChatIDList     []string     `json:"chat_id_list"`
	TagFilter      TagFilter    `json:"tag_filter"`
	Sender         string       `json:"sender"`
	AllowSelect    bool         `json:"allow_select"`
	Text           TextContent  `json:"text"`
	Attachments    []Attachment `json:"attachments"`
}

type TagFilter struct {
	GroupList []Group `json:"group_list"`
}

type Group struct {
	TagList []string `json:"tag_list"`
}

type TextContent struct {
	Content string `json:"content"`
}

type Attachment struct {
	MsgType     string       `json:"msgtype"`
	Image       *Image       `json:"image,omitempty"`
	Link        *Link        `json:"link,omitempty"`
	MiniProgram *MiniProgram `json:"miniprogram,omitempty"`
	Video       *Video       `json:"video,omitempty"`
	File        *File        `json:"file,omitempty"`
}

type Image struct {
	MediaID string `json:"media_id"`
	PicURL  string `json:"pic_url"`
}

type Link struct {
	Title  string `json:"title"`
	PicURL string `json:"picurl"`
	Desc   string `json:"desc"`
	URL    string `json:"url"`
}

type MiniProgram struct {
	Title      string `json:"title"`
	PicMediaID string `json:"pic_media_id"`
	AppID      string `json:"appid"`
	Page       string `json:"page"`
}

type Video struct {
	MediaID string `json:"media_id"`
}

type File struct {
	MediaID string `json:"media_id"`
}

type AddMsgTemplateResponse struct {
	utils.WeixinError
	FailList []string `json:"fail_list"`
	MsgID    string   `json:"msgid"`
}

func (api *ExternalContactApi) AddMsgTemplate(
	ctx context.Context,
	chatType ChatType,
	externalUserID []string,
	chatIDList []string,
	tagFilter TagFilter,
	sender string,
	allowSelect bool,
	text TextContent,
	attachments []Attachment,
) (*AddMsgTemplateResponse, error) {
	result := &AddMsgTemplateResponse{}
	if err := api.Client.HTTPPostJson(
		ctx,
		apiAddMsgTemplate,
		&AddMsgTemplateRequest{
			ChatType:       string(chatType),
			ExternalUserID: externalUserID,
			ChatIDList:     chatIDList,
			TagFilter:      tagFilter,
			Sender:         sender,
			AllowSelect:    allowSelect,
			Text:           text,
			Attachments:    attachments,
		},
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}
