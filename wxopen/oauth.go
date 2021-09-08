package wxopen

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

var OauthAuthorizeServerUrl = "https://open.weixin.qq.com"

const (
	apiAuthorize    = "/connect/oauth2/authorize"
	apiAccessToken  = "/sns/oauth2/component/access_token"
	apiRefreshToken = "/sns/oauth2/component/refresh_token"
	apiUserInfo     = "/sns/userinfo"
)

// 代公众号发起网页授权
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Before_Develop/Official_Accounts/official_account_website_authorization.html

func (api *WxOpen) GetAuthorizeUrl(
	authorizerAppID string,
	redirectUri string,
	scope string, // 用逗号隔开
	state string,
) (authorizeUrl string) {
	params := url.Values{}
	params.Add("appid", authorizerAppID)
	params.Add("component_appid", api.Config.Appid)
	params.Add("redirect_uri", redirectUri)
	params.Add("response_type", "code")
	params.Add("scope", scope)
	params.Add("state", state)
	return OauthAuthorizeServerUrl + apiAuthorize + "?" + params.Encode()
}

type OauthAccessToken struct {
	utils.WeixinError
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

// 通过 code 换取 access_token
func (api *WxOpen) GetSnsAccessToken(
	ctx context.Context,
	authorizerAppID, code string,
) (*OauthAccessToken, error) {
	result := &OauthAccessToken{}
	if err := api.Client.HTTPGetWithParams(context.TODO(), apiAccessToken, func(params url.Values) {
		params.Add("appid", authorizerAppID)
		params.Add("component_appid", api.Config.Appid)
		params.Add("code", code)
		params.Add("grant_type", "authorization_code")
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

// 刷新 access_token（如果需要）
func (api *WxOpen) RefreshSnsToken(
	ctx context.Context,
	authorizerAppID, refreshToken string,
) (*OauthAccessToken, error) {
	result := &OauthAccessToken{}
	if err := api.Client.HTTPGetWithParams(context.TODO(), apiRefreshToken, func(params url.Values) {
		params.Add("appid", authorizerAppID)
		params.Add("component_appid", api.Config.Appid)
		params.Add("grant_type", "refresh_token")
		params.Add("refresh_token", refreshToken)
	}, result); err != nil {
		return nil, err
	}

	return result, nil
}

const (
	LANG_zh_CN = "zh_CN"
	LANG_zh_TW = "zh_TW"
	LANG_en    = "en"
)

type OauthUserInfo struct {
	utils.WeixinError
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int64    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

// 通过网页授权 access_token 获取用户基本信息（需授权作用域为 snsapi_userinfo）
func (api *WxOpen) GetUserInfo(
	ctx context.Context, accessToken string, openid string, lang string,
) (*OauthUserInfo, error) {
	result := &OauthUserInfo{}
	if err := api.Client.HTTPGetToken(context.TODO(), apiUserInfo, func(params url.Values) {
		params.Add("access_token", accessToken)
		params.Add("openid", openid)
		params.Add("lang", lang)
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}
