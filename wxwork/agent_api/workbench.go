package agent_api

// Package 设置工作台自定义展示

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiSetWorkbenchTemplate = "/cgi-bin/agent/set_workbench_template"
	apiGetWorkbenchTemplate = "/cgi-bin/agent/get_workbench_template"
	apiSetWorkbenchData     = "/cgi-bin/agent/set_workbench_data"
)

const (
	KeyTypeKeyData = "keydata"
	KeyTypeImage   = "image"
	KeyTypeList    = "list"
	KeyTypeWebview = "webview"
)

type AgentApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *AgentApi {
	return &AgentApi{client}
}

// 关键数据型
type WorkBenchKeyDataItemSection struct {
	Key      string `json:"key"`
	Data     string `json:"data"`
	JumpURL  string `json:"jump_url"`
	PagePath string `json:"pagepath,omitempty"`
}

type WorkBenchKeyDataItem struct {
	Items []*WorkBenchKeyDataItemSection `json:"items"`
}

// 图片型
type WorkBenchImageItem struct {
	URL      string `json:"url"`
	JumpURL  string `json:"jump_url"`
	PagePath string `json:"pagepath,omitempty"`
}

// 列表型
type WorkBenchListItemSection struct {
	Title    string `json:"title"`
	JumpURL  string `json:"jump_url"`
	PagePath string `json:"pagepath,omitempty"`
}

type WorkBenchListItem struct {
	Items []*WorkBenchListItemSection `json:"items"`
}

// webview型
type WorkBenchWebviewItem WorkBenchImageItem

/*
获取应用在工作台展示的模版
See: https://work.weixin.qq.com/api/doc/90000/90135/92535
	https://work.weixin.qq.com/api/doc/90001/90143/94620
POST https://qyapi.weixin.qq.com/cgi-bin/agent/get_workbench_template?access_token=ACCESS_TOKEN
*/
type WorkbenchTemplate struct {
	utils.WeixinError
	Type             string                `json:"type"`
	ReplaceUserData  bool                  `json:"replace_user_data"`
	WorkBenchKeyData *WorkBenchKeyDataItem `json:"keydata,omitempty"`
	WorkBenchImage   *WorkBenchImageItem   `json:"image,omitempty"`
	WorkBenchList    *WorkBenchListItem    `json:"list,omitempty"`
	WorkBenchWebview *WorkBenchWebviewItem `json:"webview,omitempty"`
}

type WorkbenchTemplateResp struct {
	utils.WeixinError
	WorkbenchTemplate
}

func (api *AgentApi) GetWorkbenchTemplate(
	ctx context.Context, agentID int,
) (*WorkbenchTemplateResp, error) {
	result := &WorkbenchTemplateResp{}
	if err := api.Client.HTTPPostJson(ctx, apiGetWorkbenchTemplate, map[string]int{
		"agentid": agentID,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
设置应用在工作台展示的模版
See: https://work.weixin.qq.com/api/doc/90000/90135/92535
	https://work.weixin.qq.com/api/doc/90001/90143/94620
POST https://qyapi.weixin.qq.com/cgi-bin/agent/set_workbench_template?access_token=ACCESS_TOKEN
*/
type WorkbenchTemplateParam struct {
	AgentID int `json:"agentid"`
	WorkbenchTemplate
}

func (api *AgentApi) SetWorkbenchTemplate(
	ctx context.Context, param *WorkbenchTemplateParam,
) error {
	return api.Client.HTTPPostJson(ctx, apiSetWorkbenchTemplate, param, nil)
}

/*
设置应用在用户工作台展示的数据
See: https://work.weixin.qq.com/api/doc/90000/90135/92535
	https://work.weixin.qq.com/api/doc/90001/90143/94620
POST https://qyapi.weixin.qq.com/cgi-bin/agent/set_workbench_data?access_token=ACCESS_TOKEN
*/
type WorkbenchDataParam struct {
	AgentID int    `json:"agentid"`
	UserID  string `json:"userid"`
	WorkbenchTemplate
}

func (api *AgentApi) SetWorkbenchData(
	ctx context.Context, param *WorkbenchDataParam,
) error {
	return api.Client.HTTPPostJson(ctx, apiSetWorkbenchData, param, nil)
}
