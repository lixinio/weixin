package externalcontact_api

/**
	客户联系
	https://developer.work.weixin.qq.com/document/path/92109
	https://developer.work.weixin.qq.com/document/path/92238
**/

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	// 获取配置了客户联系功能的成员列表
	apiGetFollowUserList = "/cgi-bin/externalcontact/get_follow_user_list"
	// 获取客户列表
	apiGetExternalContactList = "/cgi-bin/externalcontact/list"
	// 获取客户详情
	apiGetExternalContact = "/cgi-bin/externalcontact/get"
	// 批量获取客户详情
	// apiGetExternalContactBatch = "/cgi-bin/externalcontact/batch/get_by_user"

	// https://developer.work.weixin.qq.com/document/path/92577
	// 配置客户联系「联系我」方式
	// apiUpdateTaskcard = "/cgi-bin/externalcontact/add_contact_way"
	// 获取企业已配置的「联系我」方式
	// apiAppchatCreate = "/cgi-bin/externalcontact/get_contact_way"
	// 获取企业已配置的「联系我」列表
	// apiAppchatUpdate = "/cgi-bin/externalcontact/list_contact_way"
	// 更新企业已配置的「联系我」方式
	// apiAppchatGet = "/cgi-bin/externalcontact/update_contact_way"
	// 删除企业已配置的「联系我」方式
	// apiAppchatSend = "/cgi-bin/externalcontact/del_contact_way"
	// 结束临时会话
	// apiLinkedcorpMessageSend = "/cgi-bin/externalcontact/close_temp_chat"
)

type ExternalContactApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *ExternalContactApi {
	return &ExternalContactApi{Client: client}
}

type GetFollowUserListResponse struct {
	utils.WeixinError
	FollowUsers []string `json:"follow_user"`
}

// https://developer.work.weixin.qq.com/document/path/92576
// 获取配置了客户联系功能的成员列表
func (api *ExternalContactApi) GetFollowUserList(
	ctx context.Context,
) (*GetFollowUserListResponse, error) {
	result := &GetFollowUserListResponse{}
	if err := api.Client.HTTPGet(ctx, apiGetFollowUserList, result); err != nil {
		return nil, err
	}
	return result, nil
}

type apiGetExternalContactListResponse struct {
	utils.WeixinError
	ExternalUserids []string `json:"external_userid"`
}

// https://developer.work.weixin.qq.com/document/path/92113
// 获取客户列表
func (api *ExternalContactApi) GetExternalContactList(
	ctx context.Context,
	userid string,
) (*apiGetExternalContactListResponse, error) {
	result := &apiGetExternalContactListResponse{}

	if err := api.Client.HTTPGetWithParams(ctx, apiGetExternalContactList, func(params url.Values) {
		params.Add("userid", userid)
	}, result); err != nil {
		return nil, err
	}

	return result, nil
}

type GetExternalContactResponse struct {
	utils.WeixinError
	ExternalContact map[string]interface{} `json:"external_contact"`
}

// https://developer.work.weixin.qq.com/document/path/92114
// 获取客户详情
func (api *ExternalContactApi) GetExternalContact(
	ctx context.Context,
	externalUserid string,
	cursor string,
) (*GetExternalContactResponse, error) {
	result := &GetExternalContactResponse{}

	if err := api.Client.HTTPGetWithParams(ctx, apiGetExternalContact, func(params url.Values) {
		if externalUserid != "" {
			params.Add("external_userid", externalUserid)
		}

		if cursor != "" {
			params.Add("cursor", cursor)
		}
	}, result); err != nil {
		return nil, err
	}

	return result, nil
}
