package id_transfer_api

import (
	"context"
	"github.com/lixinio/weixin/utils"
)

const (
	apiUnionID2ExternalUserID   = "/cgi-bin/idconvert/unionid_to_external_userid"
	apiExternalUserID2PendingID = "/cgi-bin/idconvert/batch/external_userid_to_pending_id"
	apiUserID2OpenUserID        = "/cgi-bin/batch/userid_to_openuserid"
)

type idTransferApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *idTransferApi {
	return &idTransferApi{Client: client}
}

type SubjectType int

const (
	// 0表示主体名称是企业的 (默认)
	SubjectTypeEnterprise = SubjectType(0)
	// 1表示主体名称是服务商的
	SubjectTypeProvider = SubjectType(1)
)

type UnionID2ExternalUserIDParam struct {
	UnionID     string      `json:"unionid"`
	OpenID      string      `json:"openid"`
	SubjectType SubjectType `json:"subject_type"` // 小程序或公众号的主体类型
}

type UnionID2ExternalUserIDResponse struct {
	utils.WeixinError
	ExternalUserID string `json:"external_userid"`
	PendingID      string `json:"pending_id"`
}

// https://developer.work.weixin.qq.com/document/path/95900
// UnionID转换成企业第三方的external_userid
func (api *idTransferApi) UnionID2ExternalUserID(
	ctx context.Context,
	unionID, openID string,
	subjectType SubjectType,
) (*UnionID2ExternalUserIDResponse, error) {
	result := &UnionID2ExternalUserIDResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiUnionID2ExternalUserID,
		&UnionID2ExternalUserIDParam{
			UnionID:     unionID,
			OpenID:      openID,
			SubjectType: subjectType,
		},
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

type ExternalUserID2PendingIDParam struct {
	// 群聊Id, 如果有传入该参数，则只检查群主是否在可见范围，同时会忽略在该群以外的external_userid。如果不传入该参数，则只检查客户跟进人是否在可见范围内。
	ChatID         string   `json:"chat_id,omitempty"`
	ExternalUserID []string `json:"external_userid"` //外部联系人ID，最多可同时查询100个外部联系人
}

type ExternalUserID2PendingIDResponse struct {
	utils.WeixinError
	Result []struct {
		ExternalUserID string `json:"external_userid"`
		PendingID      string `json:"pending_id"`
	} `json:"result"`
}

func (api *idTransferApi) ExternalUserID2PendingID(
	ctx context.Context,
	externalUserIDs []string,
	chatID string,
) (*ExternalUserID2PendingIDResponse, error) {
	result := &ExternalUserID2PendingIDResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiExternalUserID2PendingID,
		&ExternalUserID2PendingIDParam{
			ChatID:         chatID,
			ExternalUserID: externalUserIDs,
		},
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

type UserID2OpenUserIDParam struct {
	UserIDList []string `json:"userid_list"`
}

type UserID2OpenUserIDResponse struct {
	utils.WeixinError
	OpenUserIDList []struct {
		UserID     string `json:"userid"`
		OpenUserID string `json:"open_userid"`
	} `json:"open_userid_list"`
	InvalidUserIDList []string `json:"invalid_userid_list"`
}

func (api *idTransferApi) UserID2OpenUserID(
	ctx context.Context,
	userIDs []string,
) (*UserID2OpenUserIDResponse, error) {
	result := &UserID2OpenUserIDResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiUserID2OpenUserID,
		&UserID2OpenUserIDParam{
			UserIDList: userIDs,
		},
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}
