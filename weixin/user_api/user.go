// Copyright 2020 FastWeGo
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

// Package user 用户管理
package user_api

import (
	"bytes"
	"net/url"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/weixin/official_account"
)

const (
	apiUpdateRemark     = "/cgi-bin/user/info/updateremark"
	apiGetUserInfo      = "/cgi-bin/user/info"
	apiBatchGetUserInfo = "/cgi-bin/user/info/batchget"
	apiGet              = "/cgi-bin/user/get"
	apiGetBlackList     = "/cgi-bin/tags/members/getblacklist"
	apiBatchBlackList   = "/cgi-bin/tags/members/batchblacklist"
	apiBatchUnBlackList = "/cgi-bin/tags/members/batchunblacklist"
)

type UserApi struct {
	*utils.Client
}

func NewOfficialAccountApi(officialAccount *official_account.OfficialAccount) *UserApi {
	return &UserApi{
		Client: officialAccount.Client,
	}
}

/*
设置用户备注名

开发者可以通过该接口对指定用户设置备注名，该接口暂时开放给微信认证的服务号

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/Configuring_user_notes.html

POST https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token=ACCESS_TOKEN
*/
func (api *UserApi) UpdateRemarkRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiUpdateRemark, bytes.NewReader(payload), "application/json;charset=utf-8")
}
func (api *UserApi) UpdateRemark(openID, remark string) error {
	var result utils.CommonError
	return utils.ApiPostWrapper(api.UpdateRemarkRaw, map[string]string{
		"openid": openID,
		"remark": remark,
	}, &result)
}

type User struct {
	Subscribe      int32   `json:"subscribe"`
	OpenID         string  `json:"openid"`
	Nickname       string  `json:"nickname"`
	Sex            int32   `json:"sex"`
	City           string  `json:"city"`
	Country        string  `json:"country"`
	Province       string  `json:"province"`
	Language       string  `json:"language"`
	Headimgurl     string  `json:"headimgurl"`
	SubscribeTime  int32   `json:"subscribe_time"`
	UnionID        string  `json:"unionid"`
	Remark         string  `json:"remark"`
	GroupID        int32   `json:"groupid"`
	TagIDList      []int32 `json:"tagid_list"`
	SubscribeScene string  `json:"subscribe_scene"`
	QrScene        int     `json:"qr_scene"`
	QrSceneStr     string  `json:"qr_scene_str"`
}

type UserInfo struct {
	User
}

/*
获取用户基本信息

在关注者与公众号产生消息交互后，公众号可获得关注者的OpenID（加密后的微信号，每个用户对每个公众号的OpenID是唯一的。对于不同公众号，同一用户的openid不同）。公众号可通过本接口来根据OpenID获取用户基本信息，包括昵称、头像、性别、所在城市、语言和关注时间

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId

GET https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
*/
func (api *UserApi) GetUserInfoRaw(params url.Values) (resp []byte, err error) {
	return api.Client.HTTPGet(apiGetUserInfo + "?" + params.Encode())
}
func (api *UserApi) GetUserInfo(openid, lang string) (*UserInfo, error) {
	var result UserInfo
	err := utils.ApiGetWrapper(api.GetUserInfoRaw, func(params url.Values) {
		params.Add("openid", openid)
		params.Add("lang", lang)
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type BatchGetUserParams struct {
	UserList []struct {
		OpenID string `json:"openid"`
		Lang   string `json:"lang"`
	} `json:"user_list"`
}
type UserInfoList struct {
	UserInfoList []User `json:"user_info_list"`
}

/*
批量获取用户基本信息

开发者可通过该接口来批量获取用户基本信息。最多支持一次拉取100条

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId

POST https://api.weixin.qq.com/cgi-bin/user/info/batchget?access_token=ACCESS_TOKEN
*/
func (api *UserApi) BatchGetUserInfoRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiBatchGetUserInfo, bytes.NewReader(payload), "application/json;charset=utf-8")
}
func (api *UserApi) BatchGetUserInfo(param *BatchGetUserParams) (*UserInfoList, error) {
	var result UserInfoList
	err := utils.ApiPostWrapper(api.BatchGetUserInfoRaw, param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// OpenidList 用户列表
type OpenidList struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		OpenIDs []string `json:"openid"`
	} `json:"data"`
	NextOpenID string `json:"next_openid"`
}

/*
获取用户列表

公众号可通过本接口来获取帐号的关注者列表，关注者列表由一串OpenID（加密后的微信号，每个用户对每个公众号的OpenID是唯一的）组成。一次拉取调用最多拉取10000个关注者的OpenID，可以通过多次拉取的方式来满足需求

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/Getting_a_User_List.html

GET https://api.weixin.qq.com/cgi-bin/user/get?access_token=ACCESS_TOKEN&next_openid=NEXT_OPENID
*/
func (api *UserApi) GetRaw(params url.Values) (resp []byte, err error) {
	return api.Client.HTTPGet(apiGet + "?" + params.Encode())
}
func (api *UserApi) Get(next_openid string) (*OpenidList, error) {
	var result OpenidList
	err := utils.ApiGetWrapper(api.GetRaw, func(params url.Values) {
		params.Add("next_openid", next_openid)
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type BlackList struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		OpenIDs []string `json:"openid"`
	} `json:"data"`
	NextOpenID string `json:"next_openid"`
}

/*
获取公众号的黑名单列表

公众号可通过该接口来获取帐号的黑名单列表，黑名单列表由一串 OpenID（加密后的微信号，每个用户对每个公众号的OpenID是唯一的）组成

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/Manage_blacklist.html

POST https://api.weixin.qq.com/cgi-bin/tags/members/getblacklist?access_token=ACCESS_TOKEN
*/
func (api *UserApi) GetBlackListRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiGetBlackList, bytes.NewReader(payload), "application/json;charset=utf-8")
}
func (api *UserApi) GetBlackList(beginOpenid string) (*BlackList, error) {
	var result BlackList
	err := utils.ApiPostWrapper(api.GetBlackListRaw, map[string]string{
		"begin_openid": beginOpenid,
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/*
拉黑用户

公众号可通过该接口来拉黑一批用户，黑名单列表由一串 OpenID （加密后的微信号，每个用户对每个公众号的OpenID是唯一的）组成

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/Manage_blacklist.html

POST https://api.weixin.qq.com/cgi-bin/tags/members/batchblacklist?access_token=ACCESS_TOKEN
*/
func (api *UserApi) BatchBlackListRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiBatchBlackList, bytes.NewReader(payload), "application/json;charset=utf-8")
}
func (api *UserApi) BatchBlackList(openidList []string) (err error) {
	return utils.ApiPostWrapper(api.BatchBlackListRaw, map[string][]string{
		"openid_list": openidList,
	}, nil)
}

/*
取消拉黑用户

公众号可通过该接口来取消拉黑一批用户，黑名单列表由一串OpenID（加密后的微信号，每个用户对每个公众号的OpenID是唯一的）组成

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/Manage_blacklist.html

POST https://api.weixin.qq.com/cgi-bin/tags/members/batchunblacklist?access_token=ACCESS_TOKEN
*/
func (api *UserApi) BatchUnBlackListRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiBatchUnBlackList, bytes.NewReader(payload), "application/json;charset=utf-8")
}
func (api *UserApi) BatchUnBlackList(openidList []string) (err error) {
	return utils.ApiPostWrapper(api.BatchUnBlackListRaw, map[string][]string{
		"openid_list": openidList,
	}, nil)
}
