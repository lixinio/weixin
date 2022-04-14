package externalcontact_api

/**
	客户联系
	https://developer.work.weixin.qq.com/document/path/92109
	https://developer.work.weixin.qq.com/document/path/92238
**/

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiGetFollowUserList = "/cgi-bin/externalcontact/get_follow_user_list"

	// https://developer.work.weixin.qq.com/document/path/92577
	// 配置客户联系「联系我」方式
	apiUpdateTaskcard = "/cgi-bin/externalcontact/add_contact_way"
	// 获取企业已配置的「联系我」方式
	apiAppchatCreate = "/cgi-bin/externalcontact/get_contact_way"
	// 获取企业已配置的「联系我」列表
	apiAppchatUpdate = "/cgi-bin/externalcontact/list_contact_way"
	// 更新企业已配置的「联系我」方式
	apiAppchatGet = "/cgi-bin/externalcontact/update_contact_way"
	// 删除企业已配置的「联系我」方式
	apiAppchatSend = "/cgi-bin/externalcontact/del_contact_way"
	// 结束临时会话
	apiLinkedcorpMessageSend = "/cgi-bin/externalcontact/close_temp_chat"
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
