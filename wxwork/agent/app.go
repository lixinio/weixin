package agent

// https://github.com/fastwego/wxwork/blob/master/corporation/apis/app/app.go
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

// Package app 应用管理

import (
	"context"
	"net/url"
	"strconv"
)

const (
	MenuTypeClick            = "click"              // 点击推事件
	MenuTypeView             = "view"               // 跳转URL
	MenuTypeScanCodePush     = "scancode_push"      // 扫码推事件
	MenuTypeScanCodeWaitmsg  = "scancode_waitmsg"   // 扫码推事件 且弹出“消息接收中”提示框
	MenuTypePicSysPhoto      = "pic_sysphoto"       // 弹出系统拍照发图
	MenuTypePicPhotoOrAlbum  = "pic_photo_or_album" // 弹出拍照或者相册发图
	MenuTypePicWeixin        = "pic_weixin"         // 弹出企业微信相册发图器
	MenuTypeLocationSelect   = "location_select"    // 弹出地理位置选择器
	MenuTypeViewMiniPrograme = "view_miniprogram"   // 跳转到小程序
)

const (
	apiAgentGet             = "/cgi-bin/agent/get"
	apiAgentList            = "/cgi-bin/agent/list"
	apiAgentSet             = "/cgi-bin/agent/set"
	apiMenuCreate           = "/cgi-bin/menu/create"
	apiMenuGet              = "/cgi-bin/menu/get"
	apiMenuDelete           = "/cgi-bin/menu/delete"
	apiSetWorkbenchTemplate = "/cgi-bin/agent/set_workbench_template"
	apiGetWorkbenchTemplate = "/cgi-bin/agent/get_workbench_template"
	apiSetWorkbenchData     = "/cgi-bin/agent/set_workbench_data"
)

type MenuEntryObj struct {
	Type      string          `json:"type"`
	Name      string          `json:"name"`
	Key       string          `json:"key,omitempty"`
	Url       string          `json:"url,omitempty"`
	AppID     string          `json:"appid,omitempty"`
	Pagepath  string          `json:"pagepath,omitempty"`
	SubButton []*MenuEntryObj `json:"sub_button,omitempty"`
}

/*
获取指定的应用详情
See: https://work.weixin.qq.com/api/doc/90000/90135/90227
GET https://qyapi.weixin.qq.com/cgi-bin/agent/get?access_token=ACCESS_TOKEN&agentid=AGENTID参数说明
*/
// func (agent *Agent) AgentGet(ctx context.Context, params url.Values) (resp []byte, err error) {
// 	return agent.Client.HTTPGet(ctx, apiAgentGet+"?"+params.Encode())
// }

/*
获取access_token对应的应用列表
See: https://work.weixin.qq.com/api/doc/90000/90135/90227
GET https://qyapi.weixin.qq.com/cgi-bin/agent/list?access_token=ACCESS_TOKEN
*/
// func (agent *Agent) AgentList(ctx context.Context) (resp []byte, err error) {
// 	return agent.Client.HTTPGet(ctx, apiAgentList)
// }

/*
设置应用
See: https://work.weixin.qq.com/api/doc/90000/90135/90228
POST https://qyapi.weixin.qq.com/cgi-bin/agent/set?access_token=ACCESS_TOKEN
*/
// func (agent *Agent) AgentSet(ctx context.Context, payload []byte) (resp []byte, err error) {
// 	return agent.Client.HTTPPost(
// 		ctx,
// 		apiAgentSet,
// 		bytes.NewReader(payload),
// 		"application/json;charset=utf-8",
// 	)
// }

/*
创建菜单
See: https://work.weixin.qq.com/api/doc/90000/90135/90231
POST https://qyapi.weixin.qq.com/cgi-bin/menu/create?access_token=ACCESS_TOKEN&agentid=AGENTID
*/
func (agent *Agent) MenuCreate(ctx context.Context, agentID int, menus []MenuEntryObj) error {
	payload := struct {
		Buttons []MenuEntryObj `json:"button,omitempty"`
	}{
		Buttons: menus,
	}
	return agent.Client.HTTPPost(ctx, apiMenuCreate, payload, func(params url.Values) {
		params.Add("agentid", strconv.Itoa(agentID))
	}, nil, "")
}

/*
获取菜单
See: https://work.weixin.qq.com/api/doc/90000/90135/90232
GET https://qyapi.weixin.qq.com/cgi-bin/menu/get?access_token=ACCESS_TOKEN&agentid=AGENTID
*/
// func (agent *Agent) MenuGet(ctx context.Context, params url.Values) (resp []byte, err error) {
// 	return agent.Client.HTTPGet(ctx, apiMenuGet+"?"+params.Encode())
// }

/*
删除菜单
See: https://work.weixin.qq.com/api/doc/90000/90135/90233
GET https://qyapi.weixin.qq.com/cgi-bin/menu/delete?access_token=ACCESS_TOKEN&agentid=AGENTID
*/
func (agent *Agent) MenuDelete(ctx context.Context, agentID int) error {
	return agent.Client.HTTPGetWithParams(ctx, apiMenuDelete, func(params url.Values) {
		params.Add("agentid", strconv.Itoa(agentID))
	}, nil)
}

/*
设置应用在工作台展示的模版
See: https://work.weixin.qq.com/api/doc/90000/90135/92535
POST https://qyapi.weixin.qq.com/cgi-bin/agent/set_workbench_template?access_token=ACCESS_TOKEN
*/
// func (agent *Agent) SetWorkbenchTemplate(
// 	ctx context.Context,
// 	payload []byte,
// ) (resp []byte, err error) {
// 	return agent.Client.HTTPPost(
// 		ctx,
// 		apiSetWorkbenchTemplate,
// 		bytes.NewReader(payload),
// 		"application/json;charset=utf-8",
// 	)
// }

/*
获取应用在工作台展示的模版
See: https://work.weixin.qq.com/api/doc/90000/90135/92535
POST https://qyapi.weixin.qq.com/cgi-bin/agent/get_workbench_template?access_token=ACCESS_TOKEN
*/
// func (agent *Agent) GetWorkbenchTemplate(
// 	ctx context.Context,
// 	payload []byte,
// ) (resp []byte, err error) {
// 	return agent.Client.HTTPPost(
// 		ctx,
// 		apiGetWorkbenchTemplate,
// 		bytes.NewReader(payload),
// 		"application/json;charset=utf-8",
// 	)
// }

/*
设置应用在用户工作台展示的数据
See: https://work.weixin.qq.com/api/doc/90000/90135/92535
POST https://qyapi.weixin.qq.com/cgi-bin/agent/set_workbench_data?access_token=ACCESS_TOKEN
*/
// func (agent *Agent) SetWorkbenchData(ctx context.Context, payload []byte) (resp []byte, err error) {
// 	return agent.Client.HTTPPost(
// 		ctx,
// 		apiSetWorkbenchData,
// 		bytes.NewReader(payload),
// 		"application/json;charset=utf-8",
// 	)
// }
