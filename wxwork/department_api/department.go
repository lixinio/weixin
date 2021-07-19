package department_api

// https://github.com/fastwego/wxwork/blob/master/corporation/apis/contact/department/department.go
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

// Package department 通讯录管理/部门管理

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/wxwork/agent"
)

const (
	apiCreate = "/cgi-bin/department/create"
	apiUpdate = "/cgi-bin/department/update"
	apiDelete = "/cgi-bin/department/delete"
	apiList   = "/cgi-bin/department/list"
)

type DepartmentApi struct {
	*utils.Client
}

func NewAgentApi(agent *agent.Agent) *DepartmentApi {
	return &DepartmentApi{
		Client: agent.Client,
	}
}

type CreateParam struct {
	Name     string `json:"name"`
	NameEn   string `json:"name_en,omitempty"`
	Parentid int    `json:"parentid"`
	Order    int    `json:"order,omitempty"`
	ID       int    `json:"id,omitempty"`
}

type UpdateParam struct {
	ID       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	NameEn   string `json:"name_en,omitempty"`
	Parentid int    `json:"parentid,omitempty"`
	Order    int    `json:"order,omitempty"`
}

type DepartmentID struct {
	utils.CommonError
	ID int `json:"id"`
}

type DepartmentItem struct {
	Name     string `json:"name"`
	NameEn   string `json:"name_en"`
	Parentid int    `json:"parentid"`
	Order    int    `json:"order"`
	ID       int    `json:"id"`
}

type DepartmentList struct {
	utils.CommonError
	Department []DepartmentItem `json:"department"`
}

/*
创建部门
See: https://work.weixin.qq.com/api/doc/90000/90135/90205
POST https://qyapi.weixin.qq.com/cgi-bin/department/create?access_token=ACCESS_TOKEN
*/
func (api *DepartmentApi) CreateRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiCreate, bytes.NewReader(payload), "application/json;charset=utf-8")
}
func (api *DepartmentApi) Create(params *CreateParam) (*DepartmentID, error) {
	var result DepartmentID
	err := utils.ApiPostWrapper(api.CreateRaw, params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/*
更新部门
See: https://work.weixin.qq.com/api/doc/90000/90135/90206
POST https://qyapi.weixin.qq.com/cgi-bin/department/update?access_token=ACCESS_TOKEN
*/
func (api *DepartmentApi) UpdateRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiUpdate, bytes.NewReader(payload), "application/json;charset=utf-8")
}
func (api *DepartmentApi) Update(params *UpdateParam) error {
	var result DepartmentID
	return utils.ApiPostWrapper(api.UpdateRaw, params, &result)
}

/*
删除部门
See: https://work.weixin.qq.com/api/doc/90000/90135/90207
GET https://qyapi.weixin.qq.com/cgi-bin/department/delete?access_token=ACCESS_TOKEN&id=ID
*/
func (api *DepartmentApi) DeleteRaw(params url.Values) (resp []byte, err error) {
	return api.Client.HTTPGet(apiDelete + "?" + params.Encode())
}
func (api *DepartmentApi) Delete(id int) error {
	return utils.ApiGetWrapper(api.DeleteRaw, func(params url.Values) {
		params.Add("id", fmt.Sprintf("%d", id))
	}, nil)
}

/*
获取部门列表
See: https://work.weixin.qq.com/api/doc/90000/90135/90208
GET https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=ACCESS_TOKEN&id=ID
*/
func (api *DepartmentApi) ListRaw(params url.Values) (resp []byte, err error) {
	return api.Client.HTTPGet(apiList + "?" + params.Encode())
}
func (api *DepartmentApi) List(id int) (*DepartmentList, error) {
	var result DepartmentList
	err := utils.ApiGetWrapper(api.ListRaw, func(params url.Values) {
		if id != 0 {
			params.Add("id", fmt.Sprintf("%d", id))
		}
	}, &result)
	if err == nil {
		return &result, nil
	}
	return nil, err
}
