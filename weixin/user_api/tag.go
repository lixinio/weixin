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

// Package tags 用户标签管理
package user_api

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiTagCreate         = "/cgi-bin/tags/create"
	apiTagGet            = "/cgi-bin/tags/get"
	apiTagUpdate         = "/cgi-bin/tags/update"
	apiTagDelete         = "/cgi-bin/tags/delete"
	apiTagGetUsersByTag  = "/cgi-bin/user/tag/get"
	apiTagBatchTagging   = "/cgi-bin/tags/members/batchtagging"
	apiTagBatchUnTagging = "/cgi-bin/tags/members/batchuntagging"
	apiTagGetTagIdList   = "/cgi-bin/tags/getidlist"
)

type TagItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// TagInfo 标签信息
type TagInfo struct {
	utils.WeixinError
	Tag TagItem `json:"tag"`
}

type TagList struct {
	utils.WeixinError
	Tags []TagItem
}

// TagOpenIDList 标签用户列表
type TagOpenIDList struct {
	utils.WeixinError
	Count int `json:"count"`
	Data  struct {
		OpenIDs []string `json:"openid"`
	} `json:"data"`
	NextOpenID string `json:"next_openid"`
}

type UserTagList struct {
	utils.WeixinError
	TagIDList []int `json:"tagid_list"`
}

/*
创建标签

一个公众号，最多可以创建100个标签

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html

POST https://api.weixin.qq.com/cgi-bin/tags/create?access_token=ACCESS_TOKEN
*/
func (api *UserApi) CreateTag(ctx context.Context, tagName string) (*TagInfo, error) {
	var result TagInfo
	tag := &struct {
		Name string `json:"name"`
	}{Name: tagName}

	params := map[string]*struct {
		Name string `json:"name"`
	}{"tag": tag}

	if err := api.Client.HTTPPostJson(ctx, apiTagCreate, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

/*
获取公众号已创建的标签

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html

GET https://api.weixin.qq.com/cgi-bin/tags/get?access_token=ACCESS_TOKEN
*/
func (api *UserApi) GetTag(ctx context.Context) (*TagList, error) {
	var result TagList
	err := api.Client.HTTPGet(ctx, apiTagGet, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/*
编辑标签

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html

POST https://api.weixin.qq.com/cgi-bin/tags/update?access_token=ACCESS_TOKEN
*/
func (api *UserApi) UpdateTag(ctx context.Context, id int, name string) error {
	tag := &struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{ID: id, Name: name}

	params := map[string]*struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{"tag": tag}

	return api.Client.HTTPPostJson(ctx, apiTagUpdate, params, nil)
}

/*
删除标签

请注意，当某个标签下的粉丝超过10w时，后台不可直接删除标签。此时，开发者可以对该标签下的openid列表，先进行取消标签的操作，直到粉丝数不超过10w后，才可直接删除该标签

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html

POST https://api.weixin.qq.com/cgi-bin/tags/delete?access_token=ACCESS_TOKEN
*/
func (api *UserApi) DeleteTag(ctx context.Context, id int) error {
	tag := &struct {
		ID int `json:"id"`
	}{ID: id}
	params := map[string]*struct {
		ID int `json:"id"`
	}{"tag": tag}
	return api.Client.HTTPPostJson(ctx, apiTagDelete, params, nil)
}

/*
获取标签下粉丝列表

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html

POST https://api.weixin.qq.com/cgi-bin/user/tag/get?access_token=ACCESS_TOKEN
*/
func (api *UserApi) GetUsersByTag(
	ctx context.Context,
	tagID int,
	nextOpenid string,
) (*TagOpenIDList, error) {
	params := &struct {
		TagID      int    `json:"tagid"`
		NextOpenid string `json:"next_openid"`
	}{
		TagID:      tagID,
		NextOpenid: nextOpenid,
	}

	var result TagOpenIDList
	if err := api.Client.HTTPPostJson(ctx, apiTagGetUsersByTag, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

/*
批量为用户打标签

标签功能目前支持公众号为用户打上最多20个标签

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html

POST https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token=ACCESS_TOKEN
*/
func (api *UserApi) BatchTagging(ctx context.Context, tagID int, openIDList []string) error {
	params := &struct {
		TagID      int      `json:"tagid"`
		OpenIDList []string `json:"openid_list"`
	}{
		TagID:      tagID,
		OpenIDList: openIDList,
	}
	return api.Client.HTTPPostJson(ctx, apiTagBatchTagging, params, nil)
}

/*
批量为用户取消标签

标签功能目前支持公众号为用户打上最多20个标签

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html

POST https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token=ACCESS_TOKEN
*/
func (api *UserApi) BatchUnTagging(ctx context.Context, tagID int, openIDList []string) error {
	params := &struct {
		TagID      int      `json:"tagid"`
		OpenIDList []string `json:"openid_list"`
	}{
		TagID:      tagID,
		OpenIDList: openIDList,
	}
	return api.Client.HTTPPostJson(ctx, apiTagBatchUnTagging, params, nil)
}

/*
获取用户身上的标签列表

标签功能目前支持公众号为用户打上最多20个标签

See: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html

POST https://api.weixin.qq.com/cgi-bin/tags/getidlist?access_token=ACCESS_TOKEN
*/
func (api *UserApi) GetTagIdList(ctx context.Context, openID string) (*UserTagList, error) {
	var result UserTagList
	if err := api.Client.HTTPPostJson(ctx, apiTagGetTagIdList, map[string]string{
		"openid": openID,
	}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
