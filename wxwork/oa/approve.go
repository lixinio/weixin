package oa_api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

// https://work.weixin.qq.com/api/doc/90001/90143/92631
// https://work.weixin.qq.com/api/doc/90000/90135/91982

// https://github.com/fastwego/wxwork/blob/master/corporation/apis/oa/approve/approve.go
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

// Package approve OA/审批

const (
	apiGetTemplateDetail   = "/cgi-bin/oa/gettemplatedetail"
	apiGetApprovalInfo     = "/cgi-bin/oa/getapprovalinfo"
	apiGetApprovalDetail   = "/cgi-bin/oa/getapprovaldetail"
	apiGetOpenApprovalData = "/cgi-bin/corp/getopenapprovaldata"
	apiCopyTemplate        = "/cgi-bin/oa/approval/copytemplate"
)

type OaApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *OaApi {
	return &OaApi{Client: client}
}

type TmplItem struct {
	Text string `json:"text"`
	Lang string `json:"lang"`
}

type TmplDetail struct {
	utils.WeixinError
	TemplateNames   []TmplItem `json:"template_names"`
	TemplateContent struct {
		Controls []struct {
			Property struct {
				Control     string     `json:"control"`
				ID          string     `json:"id"`
				Require     int        `json:"require"`
				UnPrint     int        `json:"un_print"`
				Title       []TmplItem `json:"title"`
				Placeholder []TmplItem `json:"placeholder"`
			} `json:"property"`
			Config struct {
				Selector struct {
					Type    string `json:"type"`
					ExpType int    `json:"exp_type"`
					Options struct {
						Key   string     `json:"key"`
						Value []TmplItem `json:"value"`
					} `json:"options"`
				} `json:"selector"`
			} `json:"config"`
		} `json:"controls"`
	} `json:"template_content"`
}

/*
获取审批模板详情
企业可通过审批应用或自建应用Secret调用本接口，获取企业微信“审批应用”内指定审批模板的详情。
See: https://work.weixin.qq.com/api/doc/90000/90135/91982
POST https://qyapi.weixin.qq.com/cgi-bin/oa/gettemplatedetail?access_token=ACCESS_TOKEN
*/
func (api *OaApi) GetTemplateDetail(ctx context.Context, templateID string) (*TmplDetail, error) {
	result := &TmplDetail{}
	if err := api.Client.HTTPPostJson(ctx, apiGetApprovalDetail, map[string]string{
		"template_id": templateID,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
提交审批申请
企业可通过审批应用或自建应用Secret调用本接口，代应用可见范围内员工在企业微信“审批应用”内提交指定类型的审批申请。
See: https://work.weixin.qq.com/api/doc/90000/90135/91853
POST https://qyapi.weixin.qq.com/cgi-bin/oa/applyevent?access_token=ACCESS_TOKEN
*/
// func (api *OaApi) ApplyEvent(ctx context.Context, payload []byte) (resp []byte, err error) {
// 	return ctx.Client.HTTPPost(apiApplyEvent, bytes.NewReader(payload), "application/json;charset=utf-8")
// }

type TmplFilter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TmplList struct {
	utils.WeixinError
	SpNolist []string `json:"sp_no_list"`
}

/*
批量获取审批单号
See: https://work.weixin.qq.com/api/doc/90000/90135/91816
POST https://qyapi.weixin.qq.com/cgi-bin/oa/getapprovalinfo?access_token=ACCESS_TOKEN
*/
func (api *OaApi) GetApprovalInfo(
	ctx context.Context,
	starttime, endtime string,
	cursor, size int,
	filter []TmplFilter,
) (*TmplList, error) {
	result := &TmplList{}
	payload := &struct {
		StartTime string       `json:"starttime"`
		EndTime   string       `json:"endtime"`
		Cursor    int          `json:"cursor"`
		Size      int          `json:"size"`
		Filters   []TmplFilter `json:"filters,omitempty"`
	}{
		StartTime: starttime,
		EndTime:   endtime,
		Cursor:    cursor,
		Size:      size,
		Filters:   filter,
	}

	b, _ := json.Marshal(payload)
	fmt.Println(string(b))

	if err := api.Client.HTTPPostJson(
		ctx, apiGetApprovalInfo, payload, result,
	); err != nil {
		return nil, err
	}
	return result, nil
}

/*
https://work.weixin.qq.com/api/doc/90001/90143/93798
查询第三方应用审批申请当前状态
*/
type OpenApprovalDetail struct {
	ThirdNo        string   `json:"ThirdNo"`        // 审批单编号，由开发者在发起申请时自定义
	OpenTemplateID string   `json:"OpenTemplateId"` // 审批模板名称
	OpenSpName     string   `json:"OpenSpName"`     // 审批模板id
	OpenSpStatus   int      `json:"OpenSpstatus"`   // 申请单当前审批状态：1-审批中；2-已通过；3-已驳回；4-已取消
	ApplyTime      int64    `json:"ApplyTime"`      // 提交申请时间
	ApplyUserName  string   `json:"ApplyUsername"`  // 提交者姓名
	ApplyUserParty string   `json:"ApplyUserParty"` // 提交者userid
	ApplyUserImage string   `json:"ApplyUserImage"` // 提交者所在部门
	ApplyUserID    string   `json:"ApplyUserId"`    // 提交者头像
	ApproverStep   int      `json:"approverstep"`   // 当前审批节点：0-第一个审批节点；1-第二个审批节点…以此
	ApprovalNodes  struct { // 审批流程信息
		ApprovalNode []struct {
			NodeStatus int      `json:"NodeStatus"` // 节点审批操作状态：1-审批中；2-已同意；3-已驳回；4-已转审
			NodeAttr   int      `json:"NodeAttr"`   // 审批节点属性：1-或签；2-会签
			NodeType   int      `json:"NodeType"`   // 审批节点类型：1-固定成员；2-标签；3-上级
			Items      struct { // 审批节点信息，当节点为标签或上级时，一个节点可能有多个分支
				Item []struct {
					ItemName   string `json:"ItemName"`   // 分支审批人姓名
					ItemParty  string `json:"ItemParty"`  // 分支审批人userid
					ItemImage  string `json:"ItemImage"`  // 分支审批人所在部门
					ItemUserID string `json:"ItemUserId"` // 分支审批人头像
					ItemStatus int    `json:"ItemStatus"` // 分支审批审批操作状态：1-审批中；2-已同意；3-已驳回；4-已转审
					ItemSpeech string `json:"ItemSpeech"` // 分支审批人审批意见
					ItemOpTime int64  `json:"ItemOpTime"` // 分支审批人操作时间
				}
			} `json:"Items"`
		} `json:"ApprovalNode"`
	} `json:"ApprovalNodes"`
	NotifyNodes struct { // 抄送信息，可能有多个抄送人
		NotifyNode []struct {
			ItemName   string `json:"ItemName"`   // 抄送人姓名
			ItemParty  string `json:"ItemParty"`  // 抄送人userid
			ItemImage  string `json:"ItemImage"`  // 抄送人所在部门
			ItemUserId string `json:"ItemUserId"` // 抄送人头像
		} `json:"NotifyNode"`
	} `json:"NotifyNodes"`
}

type OpenApprovalData struct {
	utils.WeixinError
	Data *OpenApprovalDetail `json:"data"`
}

func (api *OaApi) GetOpenApprovalData(
	ctx context.Context,
	thirdNo string,
) (*OpenApprovalData, error) {
	result := &OpenApprovalData{}
	if err := api.Client.HTTPPostJson(ctx, apiGetOpenApprovalData, map[string]string{
		"thirdNo": thirdNo,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

// 复制/更新模板到企业
// https://work.weixin.qq.com/api/doc/90001/90143/92630
// 该接口仅限三方应用
// !!!! 实际测试发现改接口并不需要调用，直接用服务商的模板ID就行了， 调用也会报错
func (api *OaApi) CopyTemplate(ctx context.Context, openTemplateID string) (string, error) {
	result := &struct {
		utils.WeixinError
		TemplateID string `json:"template_id"`
	}{}
	if err := api.Client.HTTPPostJson(ctx, apiCopyTemplate, map[string]string{
		"open_template_id": openTemplateID,
	}, result); err != nil {
		return "", err
	}
	return result.TemplateID, nil
}
