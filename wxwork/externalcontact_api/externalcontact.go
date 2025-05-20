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
	apiListExternalContactDetails = "/cgi-bin/externalcontact/batch/get_by_user"

	// 获取客户群列表
	apiGetGroupChatList = "/cgi-bin/externalcontact/groupchat/list"
	// 获取客户群详情
	apiGetGroupChat = "/cgi-bin/externalcontact/groupchat/get"

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

type GetExternalContactListResponse struct {
	utils.WeixinError
	ExternalUserids []string `json:"external_userid"`
}

// https://developer.work.weixin.qq.com/document/path/92113
// 获取客户列表
func (api *ExternalContactApi) GetExternalContactList(
	ctx context.Context,
	userid string,
) (*GetExternalContactListResponse, error) {
	result := &GetExternalContactListResponse{}

	if err := api.Client.HTTPGetWithParams(ctx, apiGetExternalContactList, func(params url.Values) {
		params.Add("userid", userid)
	}, result); err != nil {
		return nil, err
	}

	return result, nil
}

type GetExternalContactResponse struct {
	utils.WeixinError
	// ExternalContact map[string]interface{} `json:"external_contact"`
	ExternalContact ExternalContact `json:"external_contact"`
	FollowUsers     []*FollowUser   `json:"follow_user"`
}

type ExternalContact struct {
	ExternalUserid  string           `json:"external_userid"`
	Name            string           `json:"name"`
	Position        string           `json:"position"`
	Avatar          string           `json:"avatar"`
	CorpName        string           `json:"corp_name"`
	CorpFullName    string           `json:"corp_full_name"`
	Type            uint8            `json:"type"`
	Gender          uint8            `json:"gender"`
	Unionid         string           `json:"unionid"`
	ExternalProfile *ExternalProfile `json:"external_profile"`
}

type ExternalProfile struct {
	ExternalCorpName string                    `json:"external_corp_name"`
	WechatChannels   *WechatChannel            `json:"wechat_channels"`
	ExternalAttr     []*map[string]interface{} `json:"external_attr"`
}

type WechatChannel struct {
	Nickname string `json:"nickname"`
	Status   uint8  `json:"status"`
}

type FollowUser struct {
	Userid         string         `json:"userid"`
	Remark         string         `json:"remark"`
	Description    string         `json:"description"`
	CreateTime     uint64         `json:"createtime"`
	Tags           []Tag          `json:"tags"`
	RemarkCorpName string         `json:"remark_corp_name"`
	RemarkMobiles  []string       `json:"remark_mobiles"`
	WechatChannels *WechatChannel `json:"wechat_channels"`
	OperUserid     string         `json:"oper_userid"`
	AddWay         uint8          `json:"add_way"`
	State          string         `json:"state"`
	NextCursor     string         `json:"next_cursor"`
}

type Tag struct {
	GroupName string `json:"group_name"`
	TagName   string `json:"tag_name"`
	Type      uint8  `json:"type"`
	TagID     string `json:"tag_id"`
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

// ListExternalContactRequest 批量获取客户详情请求
type ListExternalContactRequest struct {
	UseridList []string `json:"userid_list"`
	Cursor     string   `json:"cursor"`
	Limit      int32    `json:"limit"`
}

// ListExternalContactResponse 批量获取客户详情响应
type ListExternalContactResponse struct {
	utils.WeixinError
	ExternalContactList []*struct {
		// 外部联系人
		ExternalContact ExternalContact `json:"external_contact"`
		FollowUser      FollowUser      `json:"follow_info"`
	} `json:"external_contact_list"`
	NextCursor string `json:"next_cursor"`
	FailInfo   []*struct {
		UnlicensedUseridList []string `json:"unlicensed_userid_list"`
	}
}

// ListExternalContactDetails 批量获取客户详情
func (api *ExternalContactApi) ListExternalContactDetails(
	ctx context.Context,
	userIds []string,
	cursor string,
	limit int32,
) (*ListExternalContactResponse, error) {
	result := &ListExternalContactResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiListExternalContactDetails,
		&ListExternalContactRequest{
			UseridList: userIds,
			Cursor:     cursor,
			Limit:      limit,
		},
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

type GetGroupChatListResponse struct {
	utils.WeixinError
	GroupChatList []*struct {
		ChatId string `json:"chat_id"`
		Status uint8  `json:"status"`
	} `json:"group_chat_list"`
	NextCursor string `json:"next_cursor"`
}

func (api *ExternalContactApi) GetGroupChatList(
	ctx context.Context,
	status uint8,
	ownerUserIDs []string,
	cursor string,
	limit int32,
) (*GetGroupChatListResponse, error) {
	result := &GetGroupChatListResponse{}
	// {
	// 	"status_filter": 0,
	// 	"owner_filter": {
	// 		"userid_list": ["abel"]
	// 	},
	// 	"cursor" : "r9FqSqsI8fgNbHLHE5QoCP50UIg2cFQbfma3l2QsmwI",
	// 	"limit" : 10
	// }
	body := map[string]interface{}{
		"status_filter": status,
		"cursor":        cursor,
		"limit":         limit,
	}
	if len(ownerUserIDs) > 0 {
		body["owner_filter"] = map[string]interface{}{
			"userid_list": ownerUserIDs,
		}
	}
	if err := api.Client.HTTPPostJson(ctx, apiGetGroupChatList, body, result); err != nil {
		return nil, err
	}

	return result, nil
}

type GroupChatMember struct {
	Userid    string `json:"userid"`
	Type      uint8  `json:"type"`
	JoinTime  uint64 `json:"join_time"`
	JoinScene uint8  `json:"join_scene"`
	Invitor   *struct {
		Userid string `json:"userid"`
	} `json:"invitor"`
	GroupNickname string `json:"group_nickname"`
	Name          string `json:"name"`
	Unionid       string `json:"unionid"`
}

type GroupChat struct {
	ChatId     string             `json:"chat_id"`
	Name       string             `json:"name"`
	Owner      string             `json:"owner"`
	CreateTime uint64             `json:"create_time"`
	Notice     string             `json:"notice"`
	MemberList []*GroupChatMember `json:"member_list"`
	AdminList  []*struct {
		Userid string `json:"userid"`
	} `json:"admin_list"`
	MemberVersion string `json:"member_version"`
}

type GetGroupChatResponse struct {
	utils.WeixinError
	GroupChat *GroupChat `json:"group_chat"`
}

func (api *ExternalContactApi) GetGroupChat(
	ctx context.Context,
	chatId string,
	needName uint8,
) (*GetGroupChatResponse, error) {
	result := &GetGroupChatResponse{}

	if err := api.Client.HTTPPostJson(ctx, apiGetGroupChat, map[string]interface{}{
		"chat_id":   chatId,
		"need_name": needName,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}
