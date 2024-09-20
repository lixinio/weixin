package group_msg_api

import (
	"context"
	"github.com/lixinio/weixin/utils"
)

const (
	// 创建企业群发
	apiCreateGroupMSg = "/cgi-bin/externalcontact/add_msg_template"
	// 提醒成员群发
	apiRemindGroupMsg = "/cgi-bin/externalcontact/remind_groupmsg_send"
	// 停止企业群发
	apiCancelGroupMsg = "/cgi-bin/externalcontact/cancel_groupmsg_send"
	// 获取群成员发送任务列表
	apiGetGroupMsgTask = "/cgi-bin/externalcontact/get_groupmsg_task"
	// 获取群发成员执行结果
	apiGetGroupMsgSendResult = "/cgi-bin/externalcontact/get_groupmsg_send_result"
)

type GroupMsgApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *GroupMsgApi {
	return &GroupMsgApi{client}
}

// CreateGroupMsgRequest 创建企业群发请求
type CreateGroupMsgRequest struct {
	ChatType       string        `json:"chat_type"`
	ExternalUserID []string      `json:"external_userid,omitempty"`
	ChatIDList     []string      `json:"chat_id_list,omitempty"`
	TagFilter      TagFilter     `json:"tag_filter,omitempty"`
	Sender         string        `json:"sender,omitempty"`
	AllowSelect    bool          `json:"allow_select"`
	Text           Text          `json:"text"`
	Attachments    []*Attachment `json:"attachments,omitempty"`
}

type TagFilter struct {
	GroupList []Group `json:"group_list"`
}

type Group struct {
	TagList []string `json:"tag_list"`
}

type Text struct {
	Content string `json:"content"`
}

type Attachment struct {
	MsgType     string       `json:"msgtype"`
	Image       *Image       `json:"image,omitempty"`
	Link        *Link        `json:"link,omitempty"`
	Miniprogram *Miniprogram `json:"miniprogram,omitempty"`
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

type Miniprogram struct {
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

// CreateGroupMsgResponse 创建群发消息响应
type CreateGroupMsgResponse struct {
	utils.WeixinError
	MsgID    string   `json:"msgid"`
	FailList []string `json:"fail_list"`
}

func (api *GroupMsgApi) CreateGroupMsg(
	ctx context.Context,
	req *CreateGroupMsgRequest,
) (*CreateGroupMsgResponse, error) {
	result := &CreateGroupMsgResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiCreateGroupMSg,
		req,
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// RemindGroupMsgRequest 提醒成员群发请求
type RemindGroupMsgRequest struct {
	MsgID string `json:"msgid"`
}

// RemindGroupMsgResponse 提醒成员群发响应
type RemindGroupMsgResponse struct {
	utils.WeixinError
}

// RemindGroupMsg 提醒成员群发
func (api *GroupMsgApi) RemindGroupMsg(
	ctx context.Context,
	req *RemindGroupMsgRequest,
) (*RemindGroupMsgResponse, error) {
	result := &RemindGroupMsgResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiRemindGroupMsg,
		req,
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// CancelGroupMsgRequest 停止企业群发请求
type CancelGroupMsgRequest struct {
	MsgID string `json:"msgid"`
}

// CancelGroupMsgResponse 停止企业群发响应
type CancelGroupMsgResponse struct {
	utils.WeixinError
}

// CancelGroupMsg 停止企业群发
func (api *GroupMsgApi) CancelGroupMsg(
	ctx context.Context,
	req *CancelGroupMsgRequest,
) (*CancelGroupMsgResponse, error) {
	result := &CancelGroupMsgResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiCancelGroupMsg,
		req,
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// GetGroupMsgTaskRequest 获取群发成员发送任务列表请求
type GetGroupMsgTaskRequest struct {
	MsgID  string `json:"msgid"`
	Limit  int32  `json:"limit"`
	Cursor string `json:"cursor"`
}

// GetGroupMsgTaskResponse 获取群发成员发送任务列表响应
type GetGroupMsgTaskResponse struct {
	utils.WeixinError
	NextCursor string `json:"next_cursor"`
	TaskList   []*struct {
		UserID   string `json:"userid"`
		Status   int32  `json:"status"`
		SendTime int64  `json:"send_time"`
	} `json:"task_list"`
}

// GetGroupMsgTask 获取群发成员发送任务列表
func (api *GroupMsgApi) GetGroupMsgTask(
	ctx context.Context,
	req *GetGroupMsgTaskRequest,
) (*GetGroupMsgTaskResponse, error) {
	result := &GetGroupMsgTaskResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiGetGroupMsgTask,
		req,
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// GetGroupMsgSendResultRequest 获取群发成员执行结果请求
type GetGroupMsgSendResultRequest struct {
	MsgID  string `json:"msgid"`
	UserID string `json:"userid"`
	Limit  int32  `json:"limit"`
	Cursor string `json:"cursor"`
}

// GetGroupMsgSendResultResponse 获取群发成员执行结果响应
type GetGroupMsgSendResultResponse struct {
	utils.WeixinError
	NextCursor string `json:"next_cursor"`
	SendList   []*struct {
		ExternalUserID string `json:"external_userid"`
		ChatID         string `json:"chat_id"`
		UserID         string `json:"userid"`
		Status         int32  `json:"status"`
		SendTime       int64  `json:"send_time"`
	} `json:"send_list"`
}

// GetGroupMsgSendResult 获取群发成员执行结果
func (api *GroupMsgApi) GetGroupMsgSendResult(
	ctx context.Context,
	req *GetGroupMsgSendResultRequest,
) (*GetGroupMsgSendResultResponse, error) {
	result := &GetGroupMsgSendResultResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiGetGroupMsgSendResult,
		req,
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}
