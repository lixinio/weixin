package license_api

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	// 获取企业的账号列表
	apiListActivedAccount = "/cgi-bin/license/list_actived_account"
	// 获取成员的激活详情
	// get_active_info_by_user
	apiGetActiveInfoByUser = "/cgi-bin/license/get_active_info_by_user"
)

// ListAcitvedAccountRequest 获取企业的账号列表请求
type ListAcitvedAccountRequest struct {
	CorpID string `json:"corpid"`
	Limit  int32  `json:"limit"`
	Cursor string `json:"cursor"`
}

type ListActivedAccountResponse struct {
	utils.WeixinError
	NextCursor  string `json:"next_cursor"`
	HasMore     int32  `json:"has_more"`
	AccountList []*struct {
		UserID     string `json:"userid"`
		Type       int    `json:"type"`
		ExpireTime int64  `json:"expire_time"`
		ActiveTime int64  `json:"active_time"`
	} `json:"account_list"`
}

// ListActivedAccount 获取企业的账号列表
func (api *LicenseApi) ListActivedAccount(
	ctx context.Context,
	req *ListAcitvedAccountRequest,
) (*ListActivedAccountResponse, error) {
	result := &ListActivedAccountResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiListActivedAccount,
		req,
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

type GetActiveInfoByUserRequest struct {
	CorpID string `json:"corpid"`
	UserID string `json:"userid"`
}

type GetActiveInfoByUserResponse struct {
	utils.WeixinError
	ActiveStatus   int `json:"active_status"`
	ActiveInfoList []*struct {
		ActiveCode string `json:"active_code"`
		Type       int    `json:"type"`
		UserID     string `json:"userid"`
		ActiveTime int64  `json:"active_time"`
		ExpireTime int64  `json:"expire_time"`
	} `json:"active_info_list"`
}

// GetActiveInfoByUser 获取成员的激活详情
func (api *LicenseApi) GetActiveInfoByUser(
	ctx context.Context,
	req *GetActiveInfoByUserRequest,
) (*GetActiveInfoByUserResponse, error) {
	result := &GetActiveInfoByUserResponse{}

	if err := api.Client.HTTPPostJson(
		ctx,
		apiGetActiveInfoByUser,
		req,
		result,
	); err != nil {
		return nil, err
	}

	return result, nil
}
