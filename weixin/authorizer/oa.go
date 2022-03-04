package authorizer

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiWxaMpLinkGet = "/cgi-bin/wxopen/wxamplinkget"
	apiWxaMpLink    = "/cgi-bin/wxopen/wxamplink"
	apiWxaMpUnlink  = "/cgi-bin/wxopen/wxampunlink"
)

// 获取公众号关联的小程序
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Official__Accounts/Mini_Program_Management_Permission.html
type WxaMpLinkItem struct {
	Status              int        `json:"status"`   // 1：已关联； 2：等待小程序管理员确认中；3：小程序管理员拒绝关联 12：等待公众号管理员确认中；
	Username            string     `json:"username"` // 小程序 gh_id
	AppID               string     `json:"appid"`    // 小程序 appid
	Source              string     `json:"source"`
	NickName            string     `json:"nickname"`              // 小程序名称
	Selected            int        `json:"selected"`              // 是否在公众号管理页展示中
	NearbyDisplayStatus int        `json:"nearby_display_status"` // 是否展示在附近的小程序中
	Released            int        `json:"released"`              // 是否已经发布
	HeadimgUrl          string     `json:"headimg_url"`           // 头像 url
	CopyVerifyStatus    int        `json:"copy_verify_status"`
	Email               string     `json:"email"` // 小程序邮箱
	FuncInfos           []struct { // 微信认证及支付信息，0 表示未开通，1 表示开通
		Status int    `json:"status"`
		ID     int    `json:"id"`
		Name   string `json:"name"`
	} `json:"func_infos"`
}

func (api *Authorizer) WxaMpLinkGet(
	ctx context.Context,
) ([]WxaMpLinkItem, error) {
	result := struct {
		utils.WeixinError
		WxOpens struct {
			Items []WxaMpLinkItem `json:"items"`
		} `json:"wxopens"`
	}{}
	err := api.Client.HTTPPostJson(ctx, apiWxaMpLinkGet, map[int]int{}, &result)
	if err != nil {
		return nil, err
	}
	return result.WxOpens.Items, nil
}

// 关联小程序
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Official__Accounts/Mini_Program_Management_Permission.html
func (api *Authorizer) WxaMpLink(
	ctx context.Context, appid, notifyUsers, showProfile string,
) error {
	params := map[string]string{
		"appid":        appid,       // 小程序 appid
		"notify_users": notifyUsers, // 是否发送模板消息通知公众号粉丝
		"show_profile": showProfile, // 是否展示公众号主页中
	}

	return api.Client.HTTPPostJson(ctx, apiWxaMpLink, params, nil)
}

// 解除已关联的小程序
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Official__Accounts/Mini_Program_Management_Permission.html
func (api *Authorizer) WxaMpUnLink(
	ctx context.Context, appid string,
) error {
	params := map[string]string{
		"appid": appid, // 小程序 appid
	}

	return api.Client.HTTPPostJson(ctx, apiWxaMpUnlink, params, nil)
}
