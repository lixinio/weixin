package license_api

import (
	"context"
	"github.com/lixinio/weixin/utils"
)

const (
	// 获取企业的账号列表
	apiListActivedAccount = "/cgi-bin/license/list_actived_account"
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
