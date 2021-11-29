package web_sso

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiAuthorize    = "https://open.weixin.qq.com/connect/qrconnect"
	apiAccessToken  = "/sns/oauth2/access_token"
	apiRefreshToken = "/sns/oauth2/refresh_token"
	WXServerUrl     = "https://api.weixin.qq.com" // 微信 api 服务器地址
)

type OauthAccessToken struct {
	utils.WeixinError
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}

type Config struct {
	Appid  string
	Secret string
}

/*
网页扫码
*/
type WebSSO struct {
	Config *Config
	Client *utils.Client
}

func New(config *Config) *WebSSO {
	instance := &WebSSO{
		Config: config,
	}
	instance.Client = utils.NewClient(
		WXServerUrl,
		utils.EmptyClientAccessTokenGetter(0),
	)
	return instance
}

// https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
// 请求CODE
func (sso *WebSSO) GetAuthorizeUrl(redirectUri string, state string) (authorizeUrl string) {
	params := url.Values{}
	params.Add("appid", sso.Config.Appid)
	params.Add("redirect_uri", redirectUri)
	params.Add("response_type", "code")
	params.Add("scope", "snsapi_login")
	params.Add("state", state)
	return apiAuthorize + "?" + params.Encode() + "#wechat_redirect"
}

// https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
// 通过code获取access_token
// https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code
func (sso *WebSSO) GetSnsAccessToken(
	ctx context.Context,
	code string,
) (*OauthAccessToken, error) {
	result := &OauthAccessToken{}
	// 无需 access token
	if err := sso.Client.HTTPGetToken(context.TODO(), apiAccessToken, func(params url.Values) {
		params.Add("appid", sso.Config.Appid)
		params.Add("secret", sso.Config.Secret)
		params.Add("code", code)
		params.Add("grant_type", "authorization_code")
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

// https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
// 刷新access_token有效期
// https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=APPID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN
func (sso *WebSSO) RefreshSnsToken(
	ctx context.Context,
	refreshToken string,
) (*OauthAccessToken, error) {
	result := &OauthAccessToken{}
	// 无需 access token
	if err := sso.Client.HTTPGetToken(context.TODO(), apiRefreshToken, func(params url.Values) {
		params.Add("appid", sso.Config.Appid)
		params.Add("grant_type", "refresh_token")
		params.Add("refresh_token", refreshToken)
	}, result); err != nil {
		return nil, err
	}

	return result, nil
}
