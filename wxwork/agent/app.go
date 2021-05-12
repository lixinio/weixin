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
	"bytes"
	"net/url"
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

/*
获取指定的应用详情
See: https://work.weixin.qq.com/api/doc/90000/90135/90227
GET https://qyapi.weixin.qq.com/cgi-bin/agent/get?access_token=ACCESS_TOKEN&agentid=AGENTID参数说明
*/
func (agent *Agent) AgentGet(params url.Values) (resp []byte, err error) {
	return agent.Client.HTTPGet(apiAgentGet + "?" + params.Encode())
}

/*
获取access_token对应的应用列表
See: https://work.weixin.qq.com/api/doc/90000/90135/90227
GET https://qyapi.weixin.qq.com/cgi-bin/agent/list?access_token=ACCESS_TOKEN
*/
func (agent *Agent) AgentList() (resp []byte, err error) {
	return agent.Client.HTTPGet(apiAgentList)
}

/*
设置应用
See: https://work.weixin.qq.com/api/doc/90000/90135/90228
POST https://qyapi.weixin.qq.com/cgi-bin/agent/set?access_token=ACCESS_TOKEN
*/
func (agent *Agent) AgentSet(payload []byte) (resp []byte, err error) {
	return agent.Client.HTTPPost(apiAgentSet, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
创建菜单
See: https://work.weixin.qq.com/api/doc/90000/90135/90231
POST https://qyapi.weixin.qq.com/cgi-bin/menu/create?access_token=ACCESS_TOKEN&agentid=AGENTID
*/
func (agent *Agent) MenuCreate(payload []byte, params url.Values) (resp []byte, err error) {
	return agent.Client.HTTPPost(apiMenuCreate+"?"+params.Encode(), bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
获取菜单
See: https://work.weixin.qq.com/api/doc/90000/90135/90232
GET https://qyapi.weixin.qq.com/cgi-bin/menu/get?access_token=ACCESS_TOKEN&agentid=AGENTID
*/
func (agent *Agent) MenuGet(params url.Values) (resp []byte, err error) {
	return agent.Client.HTTPGet(apiMenuGet + "?" + params.Encode())
}

/*
删除菜单
See: https://work.weixin.qq.com/api/doc/90000/90135/90233
GET https://qyapi.weixin.qq.com/cgi-bin/menu/delete?access_token=ACCESS_TOKEN&agentid=AGENTID
*/
func (agent *Agent) MenuDelete(params url.Values) (resp []byte, err error) {
	return agent.Client.HTTPGet(apiMenuDelete + "?" + params.Encode())
}

/*
设置应用在工作台展示的模版
See: https://work.weixin.qq.com/api/doc/90000/90135/92535
POST https://qyapi.weixin.qq.com/cgi-bin/agent/set_workbench_template?access_token=ACCESS_TOKEN
*/
func (agent *Agent) SetWorkbenchTemplate(payload []byte) (resp []byte, err error) {
	return agent.Client.HTTPPost(apiSetWorkbenchTemplate, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
获取应用在工作台展示的模版
See: https://work.weixin.qq.com/api/doc/90000/90135/92535
POST https://qyapi.weixin.qq.com/cgi-bin/agent/get_workbench_template?access_token=ACCESS_TOKEN
*/
func (agent *Agent) GetWorkbenchTemplate(payload []byte) (resp []byte, err error) {
	return agent.Client.HTTPPost(apiGetWorkbenchTemplate, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
设置应用在用户工作台展示的数据
See: https://work.weixin.qq.com/api/doc/90000/90135/92535
POST https://qyapi.weixin.qq.com/cgi-bin/agent/set_workbench_data?access_token=ACCESS_TOKEN
*/
func (agent *Agent) SetWorkbenchData(payload []byte) (resp []byte, err error) {
	return agent.Client.HTTPPost(apiSetWorkbenchData, bytes.NewReader(payload), "application/json;charset=utf-8")
}
