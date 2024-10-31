package moment_api

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	// 创建客户朋友圈发表任务
	apiCreateMomentTask = "/cgi-bin/externalcontact/add_moment_task"
	// 获取任务创建结果
	apiGetMomentTaskResult = "/cgi-bin/externalcontact/get_moment_task_result"
	// 停止发表企业朋友圈
	apiCancelMomentTask = "/cgi-bin/externalcontact/cancel_moment_task"
	// 获取客户朋友圈企业发表的列表
	apiGetMomentTaskList = "/cgi-bin/externalcontact/get_moment_task"
)

type MomentApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *MomentApi {
	return &MomentApi{client}
}

type Text struct {
	Content string `json:"content"`
}

type Attachment struct {
	MsgType string `json:"msgtype"`
	Image   *Image `json:"image,omitempty"`
	Link    *Link  `json:"link,omitempty"`
	Video   *Video `json:"video,omitempty"`
}

type Image struct {
	MediaID string `json:"media_id"`
}

type Link struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	MediaID string `json:"media_id"`
}
type Video struct {
	MediaID string `json:"media_id"`
}

// CreateMomentTaskRequest 创建客户朋友圈发表任务请求
//
//	{
//		"text": {
//			"content": "文本消息内容"
//		},
//		"attachments": [
//			{
//				"msgtype": "image",
//				"image": {
//					"media_id": "MEDIA_ID"
//				}
//			},
//			{
//				"msgtype": "video",
//				"video": {
//					"media_id": "MEDIA_ID"
//				}
//			},
//			{
//				"msgtype": "link",
//				"link": {
//					"title": "消息标题",
//					"url": "https://example.link.com/path",
//					"media_id": "MEDIA_ID"
//				}
//			}
//		],
//	 	"visible_range":{
//			"sender_list":{
//				"user_list":["zhangshan","lisi"],
//				"department_list":[2,3]
//			},
//			"external_contact_list":{
//				"tag_list":[ "etXXXXXXXXXX", "etYYYYYYYYYY"]
//			}
//		}
//	}

type SenderList struct {
	UserList       []string `json:"user_list"`
	DepartmentList []int32  `json:"department_list"`
}

type ExternalContactList struct {
	TagList []string `json:"tag_list"`
}

type VisibleRange struct {
	SenderList          *SenderList          `json:"sender_list"`
	ExternalContactList *ExternalContactList `json:"external_contact_list"`
}

type CreateMomentTaskRequest struct {
	Text         Text          `json:"text"`
	Attachments  []*Attachment `json:"attachments,omitempty"`
	VisibleRange VisibleRange  `json:"visible_range"`
}

// CreateMomentTaskResponse 创建客户朋友圈发表任务响应
//
//	{
//		"errcode":0,
//		"errmsg":"ok",
//		"jobid":"xxxx"
//	}
type CreateMomentTaskResponse struct {
	utils.WeixinError
	JobID string `json:"jobid"`
}

// CreateMomentTask 创建客户朋友圈发表任务
func (api *MomentApi) CreateMomentTask(
	ctx context.Context,
	req *CreateMomentTaskRequest,
) (*CreateMomentTaskResponse, error) {
	result := &CreateMomentTaskResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiCreateMomentTask,
		req,
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// GetMomentTaskResultRequest 获取任务创建结果请求
type GetMomentTaskResultRequest struct {
	JobID string `json:"jobid"`
}

// GetMomentTaskResultResponse 获取任务创建结果响应
// {
//     "errcode": 0,
//     "errmsg": "ok",
//     "status": 1,
//     "type": "add_moment_task",
// 	"result": {
// 		"errcode":0,
// 		"errmsg":"ok",
// 		"moment_id":"xxxx",
// 		"invalid_sender_list":{
// 			"user_list":["zhangshan","lisi"],
// 			"department_list":[2,3]
// 		},
// 		"invalid_external_contact_list":{
// 			"tag_list":["xxx"]
// 		}
// 	}
// }

type MomentTaskResult struct {
	utils.WeixinError
	MomentID          string `json:"moment_id"`
	InvalidSenderList struct {
		UserList       []string `json:"user_list"`
		DepartmentList []int32  `json:"department_list"`
	} `json:"invalid_sender_list"`
	InvalidExternalContactList struct {
		TagList []string `json:"tag_list"`
	} `json:"invalid_external_contact_list"`
}

type GetMomentTaskResultResponse struct {
	utils.WeixinError
	Status int8              `json:"status"`
	Type   string            `json:"type"`
	Result *MomentTaskResult `json:"result"`
}

// GetMomentTaskResult 获取任务创建结果
func (api *MomentApi) GetMomentTaskResult(
	ctx context.Context,
	req *GetMomentTaskResultRequest,
) (*GetMomentTaskResultResponse, error) {
	result := &GetMomentTaskResultResponse{}

	if err := api.Client.HTTPGetWithParams(
		ctx,
		apiGetMomentTaskResult,
		func(v url.Values) {
			v.Add("jobid", req.JobID)
		},
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// CancelMomentTaskRequest 停止发表企业朋友圈请求
type CancelMomentTaskRequest struct {
	MomentID string `json:"moment_id"`
}

// CancelMomentTaskResponse 停止发表企业朋友圈响应
type CancelMomentTaskResponse struct {
	utils.WeixinError
}

// CancelMomentTask 停止发表企业朋友圈
func (api *MomentApi) CancelMomentTask(
	ctx context.Context,
	req *CancelMomentTaskRequest,
) (*CancelMomentTaskResponse, error) {
	result := &CancelMomentTaskResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiCancelMomentTask,
		req,
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// GetMomentTaskListRequest 获取客户朋友圈企业发表的列表请求
type GetMomentTaskListRequest struct {
	MomentID string `json:"moment_id"`
	Cursor   string `json:"cursor"`
	Limit    int32  `json:"limit"`
}

// GetMomentTaskListResponse 获取客户朋友圈企业发表的列表响应
//
//	{
//		"errcode":0,
//		"errmsg":"ok",
//		"next_cursor":"CURSOR",
//		"task_list":[
//			{
//				"userid":"zhangsan",
//				"publish_status":1
//			}
//		]
//	}
type GetMomentTaskListResponse struct {
	utils.WeixinError
	NextCursor string `json:"next_cursor"`
	TaskList   []*struct {
		UserID        string `json:"userid"`
		PublishStatus int8   `json:"publish_status"`
	} `json:"task_list"`
}

// GetMomentTaskList 获取客户朋友圈企业发表的列表
func (api *MomentApi) GetMomentTaskList(
	ctx context.Context,
	req *GetMomentTaskListRequest,
) (*GetMomentTaskListResponse, error) {
	result := &GetMomentTaskListResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiGetMomentTaskList,
		req,
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}
