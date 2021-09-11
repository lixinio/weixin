package wxwork_suite

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiAuthorize     = "https://open.weixin.qq.com/connect/oauth2/authorize"
	apiGetUserInfo   = "/cgi-bin/service/getuserinfo3rd"
	apiGetUserDetail = "/cgi-bin/service/getuserdetail3rd"
)

// /cgi-bin/service/getuserinfo3rd

// https://open.work.weixin.qq.com/api/doc/90001/90143/91120
/*
应用授权作用域。
snsapi_base：静默授权，可获取成员的基础信息（UserId与DeviceId）；
snsapi_userinfo：静默授权，可获取成员的详细信息，但不包含手机、邮箱等敏感信息；
snsapi_privateinfo：手动授权，可获取成员的详细信息，包含手机、邮箱等敏感信息（已不再支持获取手机号/邮箱）。
*/
func (suite *WxWorkSuite) GetAuthorizeUrl(redirectUri, scope, state string) (authorizeUrl string) {
	params := url.Values{}
	params.Add("appid", suite.Config.SuiteID)
	params.Add("redirect_uri", redirectUri)
	params.Add("response_type", "code")
	params.Add("scope", scope)
	params.Add("state", state)
	return apiAuthorize + "?" + params.Encode() + "#wechat_redirect"
}

type UserInfo3rd struct {
	utils.WeixinError
	CorpID     string `json:"CorpId"`
	UserID     string `json:"UserId"`
	DeviceID   string `json:"DeviceId"`
	UserTicket string `json:"user_ticket"`
	ExpiresIn  int    `json:"expires_in"`
	OpenUserID string `json:"open_userid"`
}

// 获取访问用户身份
// https://work.weixin.qq.com/api/doc/90001/90143/91121
func (suite *WxWorkSuite) GetUserInfo3rd(ctx context.Context, code string) (*UserInfo3rd, error) {
	result := &UserInfo3rd{}
	if err := suite.Client.HTTPGetWithParams(ctx, apiGetUserInfo, func(params url.Values) {
		params.Add("code", code)
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

type UserDetail3rd struct {
	utils.WeixinError
	CorpID string `json:"corpid"`
	UserID string `json:"userid"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Avatar string `json:"avatar"`
	QrCode string `json:"qr_code"`
}

// 获取访问用户敏感信息
// https://work.weixin.qq.com/api/doc/90001/90143/91122
func (suite *WxWorkSuite) GetUserDetail3rd(
	ctx context.Context,
	userTicket string,
) (*UserDetail3rd, error) {
	result := &UserDetail3rd{}
	if err := suite.Client.HTTPPostJson(ctx, apiGetUserDetail, map[string]string{
		"user_ticket": userTicket,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}
