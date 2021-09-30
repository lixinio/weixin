package wxopen

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiCreatePreAuthCode            = "/cgi-bin/component/api_create_preauthcode"
	apiGetAuthorizationRedirectUri  = "/cgi-bin/componentloginpage"
	apiGetAuthorizationRedirectUri2 = "/safe/bindcomponent"
	apiApiQueryAuth                 = "/cgi-bin/component/api_query_auth"
	apiApiAuthorizerToken           = "/cgi-bin/component/api_authorizer_token"
	apiApiGetAuthorizerInfo         = "/cgi-bin/component/api_get_authorizer_info"
	apiApiGetAuthorizerOption       = "/cgi-bin/component/api_get_authorizer_option"
	apiApiSetAuthorizerOption       = "/cgi-bin/component/api_set_authorizer_option"
	apiApiGetAuthorizerList         = "/cgi-bin/component/api_get_authorizer_list"
)

// 获取预授权码
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/pre_auth_code.html
func (api *WxOpen) CreatePreAuthCode(ctx context.Context) (string, int, error) {
	payload := map[string]string{
		"component_appid": api.Config.Appid,
	}

	result := struct {
		utils.WeixinError
		PreAuthCode string `json:"pre_auth_code"`
		ExpiresIn   int    `json:"expires_in"`
	}{}

	err := api.Client.HTTPPostJson(ctx, apiCreatePreAuthCode, payload, &result)
	if err == nil {
		return result.PreAuthCode, result.ExpiresIn, nil
	} else {
		return "", 0, err
	}
}

// 要授权的帐号类型：1 则商户点击链接后，手机端仅展示公众号、2 表示仅展示小程序，3 表示公众号和小程序都展示。
// 如果为未指定，则默认小程序和公众号都展示。第三方平台开发者可以使用本字段来控制授权的帐号类型。
const (
	AuthTypeOaOnly = "1" // 仅显示公众号
	AuthTypeMpOnly = "2" // 仅显示小程序
	AuthTypeAll    = "3" // 所有
)

/*
方式一：授权注册页面扫码授权
第三方平台方可以在自己的网站中放置“微信公众号授权”或者“小程序授权”的入口，或生成授权链接放置在移动网页中，引导公众号和小程序管理员进入授权页。
See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Authorization_Process_Technical_Description.html
GET https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=xxxx&pre_auth_code=xxxxx&redirect_uri=xxxx&auth_type=xxx
*/
func (api *WxOpen) GetComponentLoginPage(
	preAuthCode, redirectUri, authType, bizAppid string,
) string {
	params := url.Values{}
	params.Add("component_appid", api.Config.Appid)
	params.Add("pre_auth_code", preAuthCode)
	params.Add("redirect_uri", redirectUri)
	if authType != "" {
		params.Add("auth_type", authType)
	}
	if bizAppid != "" {
		params.Add("biz_appid", bizAppid)
	}
	return "https://mp.weixin.qq.com/cgi-bin/componentloginpage?" + params.Encode()
}

// /*
// 方式二：点击移动端链接快速授权
// 第三方平台方可以生成授权链接，将链接通过移动端直接发给授权管理员，管理员确认后即授权成功
// See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Authorization_Process_Technical_Description.html
// GET https://mp.weixin.qq.com/safe/bindcomponent?action=bindcomponent&auth_type=3&no_scan=1&component_appid=xxxx&pre_auth_code=xxxxx&redirect_uri=xxxx&auth_type=xxx&biz_appid=xxxx#wechat_redirect
// */
func (api *WxOpen) GetComponentLoginH5Page(
	preAuthCode, redirectUri, authType, bizAppid string,
) (uri string) {
	params := url.Values{}
	params.Add("component_appid", api.Config.Appid)
	params.Add("pre_auth_code", preAuthCode)
	params.Add("redirect_uri", redirectUri)
	params.Add("action", "bindcomponent")
	params.Add("no_scan", "1")

	if authType != "" {
		params.Add("auth_type", authType)
	}
	if bizAppid != "" {
		params.Add("biz_appid", bizAppid)
	}

	return "https://mp.weixin.qq.com/safe/bindcomponent?" + params.Encode() + "#wechat_redirect"
}

/*
拉取所有已授权的帐号信息
使用本 API 拉取当前所有已授权的帐号基本信息
See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/api_get_authorizer_list.html
POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_list?component_access_token=COMPONENT_ACCESS_TOKEN
*/

// 拉取所有已授权的帐号信息
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/Account_Authorization/api_get_authorizer_list.html
type AuthorizationLite struct {
	AuthorizerAppid        string `json:"authorizer_appid"`
	AuthorizerRefreshToken string `json:"refresh_token"`
	AuthTime               int    `json:"auth_time"`
}

// 拉取所有已授权的帐号信息
func (api *WxOpen) GetAuthorizerList(
	ctx context.Context,
	offset, count int,
) ([]AuthorizationLite, error) {
	payload := struct {
		ComponentAppid string `json:"component_appid"`
		Offset         int    `json:"offset"`
		Count          int    `json:"count"`
	}{
		ComponentAppid: api.Config.Appid,
		Offset:         offset,
		Count:          count,
	}

	result := struct {
		utils.WeixinError
		TotalCount int                 `json:"total_count"`
		List       []AuthorizationLite `json:"list"`
	}{}

	err := api.Client.HTTPPostJson(ctx, apiApiGetAuthorizerList, payload, &result)
	if err == nil {
		return result.List, nil
	} else {
		return nil, err
	}
}

// 获取授权方选项信息
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/Account_Authorization/api_get_authorizer_option.html
func (api *WxOpen) GetAuthorizerOption(
	ctx context.Context,
	authorizerAppid, optionName string,
) (string, error) {
	payload := map[string]string{
		"component_appid":  api.Config.Appid,
		"authorizer_appid": authorizerAppid,
		"option_name":      optionName,
	}

	result := struct {
		utils.WeixinError
		AuthorizerAppid string `json:"authorizer_appid"`
		OptionName      string `json:"option_name"`
		OptionValue     string `json:"option_value"`
	}{}

	err := api.Client.HTTPPostJson(ctx, apiApiGetAuthorizerOption, payload, &result)
	if err == nil {
		return result.OptionValue, nil
	} else {
		return "", err
	}
}

// 设置授权方选项信息
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/Account_Authorization/api_set_authorizer_option.html
func (api *WxOpen) SetAuthorizerOption(
	ctx context.Context,
	authorizerAppid, optionName, optionValue string,
) error {
	payload := map[string]string{
		"component_appid":  api.Config.Appid,
		"authorizer_appid": authorizerAppid,
		"option_name":      optionName,
		"option_value":     optionValue,
	}

	return api.Client.HTTPPostJson(ctx, apiApiGetAuthorizerOption, payload, nil)
}

// 使用授权码获取授权信息
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/authorization_info.html
type AuthorizationInfo struct {
	AuthorizerAppid        string                  `json:"authorizer_appid"`
	AuthorizerRefreshToken string                  `json:"authorizer_refresh_token"`
	AuthorizerAccessToken  string                  `json:"authorizer_access_token"`
	ExpiresIn              int                     `json:"expires_in"`
	FuncInfo               []AuthorizationFuncInfo `json:"func_info"`
}

// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/product/third_party_authority_instructions.html
// + https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/product/offical_account_authority.html
// + https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/product/miniprogram_authority.html
type AuthorizationFuncInfo struct {
	FuncscopeCategory struct {
		ID int `json:"id"`
	} `json:"funcscope_category"`
}

// 使用授权码获取授权信息
func (api *WxOpen) QueryAuth(
	ctx context.Context,
	authorizationCode string,
) (*AuthorizationInfo, error) {
	payload := map[string]string{
		"component_appid":    api.Config.Appid,
		"authorization_code": authorizationCode,
	}

	result := &struct {
		utils.WeixinError
		AuthorizationInfo *AuthorizationInfo `json:"authorization_info"`
	}{}
	err := api.Client.HTTPPostJson(ctx, apiApiQueryAuth, payload, result)
	if err == nil {
		return result.AuthorizationInfo, nil
	} else {
		return nil, err
	}
}

// 获取/刷新接口调用令牌
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/api_authorizer_token.html
type AuthorizerToken struct {
	utils.WeixinError
	AccessToken  string `json:"authorizer_access_token"`
	RefreshToken string `json:"authorizer_refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// 获取/刷新接口调用令牌
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/api_authorizer_token.html
func (api *WxOpen) GetAuthorizerToken(
	ctx context.Context, authorizerAppid, authorizerRefreshToken string,
) (*AuthorizerToken, error) {
	payload := map[string]string{
		"component_appid":          api.Config.Appid,
		"authorizer_appid":         authorizerAppid,
		"authorizer_refresh_token": authorizerRefreshToken,
	}

	result := &AuthorizerToken{}
	err := api.Client.HTTPPostJson(ctx, apiApiAuthorizerToken, payload, result)
	if err == nil {
		return result, nil
	} else {
		return nil, err
	}
}

// 获取授权方的帐号基本信息
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/api_get_authorizer_info.html
type AuthorizerInfo struct {
	NickName        string `json:"nick_name"`
	HeadImg         string `json:"head_img"`
	UserName        string `json:"user_name"`
	PrincipalName   string `json:"principal_name"`
	Alias           string `json:"alias"`
	Signature       string `json:"signature"`
	QrcodeUrl       string `json:"qrcode_url"`
	ServiceTypeInfo struct {
		ID int `json:"service_type_info"`
	} `json:"service_type_info"`
	VerifyTypeInfo struct {
		ID int `json:"id"`
	} `json:"verify_type_info"`
}

type AuthorizerDetail struct {
	utils.WeixinError
	AuthorizationInfo *AuthorizationInfo `json:"authorization_info"`
	AuthorizerInfo    *AuthorizerInfo    `json:"authorizer_info"`
}

// 获取授权方的帐号基本信息
func (api *WxOpen) GetAuthorizerInfo(
	ctx context.Context, authorizerAppid string,
) (*AuthorizerDetail, error) {
	payload := map[string]string{
		"component_appid":  api.Config.Appid,
		"authorizer_appid": authorizerAppid,
	}

	result := &AuthorizerDetail{}
	err := api.Client.HTTPPostJson(ctx, apiApiGetAuthorizerInfo, payload, result)
	if err == nil {
		return result, nil
	} else {
		return nil, err
	}
}
