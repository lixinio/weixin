package tag_api

// https://github.com/fastwego/wxwork/blob/master/corporation/apis/contact/tag/tag.go
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

// Package tag 通讯录管理/标签管理

import (
	"bytes"
	"net/url"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/wxwork/agent"
)

const (
	apiCreate      = "/cgi-bin/tag/create"
	apiUpdate      = "/cgi-bin/tag/update"
	apiDelete      = "/cgi-bin/tag/delete"
	apiGet         = "/cgi-bin/tag/get"
	apiAddTagUsers = "/cgi-bin/tag/addtagusers"
	apiDelTagUsers = "/cgi-bin/tag/deltagusers"
	apiList        = "/cgi-bin/tag/list"
)

type TagApi struct {
	*utils.Client
}

func NewAgentApi(agent *agent.Agent) *TagApi {
	return &TagApi{
		Client: agent.Client,
	}
}

/*
创建标签
See: https://work.weixin.qq.com/api/doc/90000/90135/90210
POST https://qyapi.weixin.qq.com/cgi-bin/tag/create?access_token=ACCESS_TOKEN
*/
func (api *TagApi) Create(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiCreate, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
更新标签名字
See: https://work.weixin.qq.com/api/doc/90000/90135/90211
POST https://qyapi.weixin.qq.com/cgi-bin/tag/update?access_token=ACCESS_TOKEN
*/
func (api *TagApi) Update(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiUpdate, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
删除标签
See: https://work.weixin.qq.com/api/doc/90000/90135/90212
GET https://qyapi.weixin.qq.com/cgi-bin/tag/delete?access_token=ACCESS_TOKEN&tagid=TAGID
*/
func (api *TagApi) Delete(params url.Values) (resp []byte, err error) {
	return api.Client.HTTPGet(apiDelete + "?" + params.Encode())
}

/*
获取标签成员
See: https://work.weixin.qq.com/api/doc/90000/90135/90213
GET https://qyapi.weixin.qq.com/cgi-bin/tag/get?access_token=ACCESS_TOKEN&tagid=TAGID
*/
func (api *TagApi) Get(params url.Values) (resp []byte, err error) {
	return api.Client.HTTPGet(apiGet + "?" + params.Encode())
}

/*
增加标签成员
See: https://work.weixin.qq.com/api/doc/90000/90135/90214
POST https://qyapi.weixin.qq.com/cgi-bin/tag/addtagusers?access_token=ACCESS_TOKEN
*/
func (api *TagApi) AddTagUsers(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiAddTagUsers, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
删除标签成员
See: https://work.weixin.qq.com/api/doc/90000/90135/90215
POST https://qyapi.weixin.qq.com/cgi-bin/tag/deltagusers?access_token=ACCESS_TOKEN
*/
func (api *TagApi) DelTagUsers(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiDelTagUsers, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
获取标签列表
See: https://work.weixin.qq.com/api/doc/90000/90135/90216
GET https://qyapi.weixin.qq.com/cgi-bin/tag/list?access_token=ACCESS_TOKEN
*/
func (api *TagApi) List() (resp []byte, err error) {
	return api.Client.HTTPGet(apiList)
}
