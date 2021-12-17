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
	apiUserInfo     = "/sns/userinfo"
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

// https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Authorized_Interface_Calling_UnionID.html
// 获取用户个人信息（UnionID机制）
// https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID
func (sso *WebSSO) GetUserInfo(
	ctx context.Context, accessToken string, openid string, lang string,
) (*OauthUserInfo, error) {
	result := &OauthUserInfo{}
	if lang == "" {
		lang = LANG_zh_CN // 国家地区语言版本，zh_CN 简体，zh_TW 繁体，en 英语，默认为en
	}
	// 无需 access token
	if err := sso.Client.HTTPGetToken(context.TODO(), apiUserInfo, func(params url.Values) {
		params.Add("access_token", accessToken)
		params.Add("openid", openid)
		params.Add("lang", lang)
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}
