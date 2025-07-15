package wxwork_provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	WXServerUrl          = "https://qyapi.weixin.qq.com" // 微信 api 服务器地址
	apiAuthorize         = "https://open.work.weixin.qq.com/wwopen/sso/3rd_qrConnect"
	apiGetProviderToken  = "/cgi-bin/service/get_provider_token"
	apiGetLoginInfo      = "/cgi-bin/service/get_login_info"
	apiGetAppLicenseInfo = "/cgi-bin/license/get_app_license_info"
)

type Config struct {
	CorpID         string // 企业服务商ID
	ProviderSecret string // 企业服务商密钥
}

type WxWorkProvider struct {
	Config           *Config
	Client           *utils.Client
	accessTokenCache *utils.AccessTokenCache
}

func New(
	cache utils.Cache, locker utils.Lock, config *Config,
	tokenRefreshHandler utils.TokenRefreshHandler, // 刷新callback
) *WxWorkProvider {
	accessTokenCache := utils.NewAccessTokenCache(
		newAccessTokenAdaptor(config), cache, locker,
		utils.CacheClientTokenOptWithExpireBefore(tokenRefreshHandler),
	)

	client := utils.NewClient(WXServerUrl, accessTokenCache)
	client.UpdateAccessTokenKey("provider_access_token")

	instance := &WxWorkProvider{
		Config:           config,
		Client:           client,
		accessTokenCache: accessTokenCache,
	}
	return instance
}

func NewLite(cache utils.Cache, locker utils.Lock, corpID string) *WxWorkProvider {
	return New(cache, locker, &Config{CorpID: corpID}, nil)
}

func (provider *WxWorkProvider) RefreshAccessToken(
	ctx context.Context, expireBefore int,
) (string, error) {
	if provider.Config.ProviderSecret == "" {
		return "", fmt.Errorf(
			"provider corpid : %s, error: %w", provider.Config.CorpID, ErrTokenUpdateForbidden,
		)
	}
	return provider.accessTokenCache.RefreshAccessToken(ctx, expireBefore)
}

// https://open.work.weixin.qq.com/api/doc/90001/90143/91124
func (provider *WxWorkProvider) GetAuthorizeUrl(
	redirectUri, userType, state string,
) (authorizeUrl string) {
	params := url.Values{}
	params.Add("appid", provider.Config.CorpID)
	params.Add("redirect_uri", redirectUri)
	params.Add("usertype", userType)
	params.Add("state", state)
	params.Add("lang", "zh") // 自定义语言，支持zh、en；lang为空则从Headers读取Accept-Language，默认值为zh
	return apiAuthorize + "?" + params.Encode()
}

type UserInfo struct {
	UserID     string `json:"userid"`
	OpenUserID string `json:"open_userid"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
}

type CorpInfo struct {
	CorpID string `json:"corpid"`
}

type LoginInfo struct {
	utils.WeixinError
	UserType  int       `json:"usertype"`
	UserInfo  *UserInfo `json:"user_info"`
	CorpInfo  *CorpInfo `json:"corp_info"`
	AgentInfo []struct {
		AgentID  int `json:"agentid"`
		AuthType int `json:"auth_type"`
	} `json:"agent"`
	AuthInfo struct {
		Department []struct {
			ID       int  `json:"id"`
			Writable bool `json:"writable"`
		} `json:"department"`
	} `json:"auth_info"`
}

// 扫码登录 获取登录用户信息
// https://open.work.weixin.qq.com/api/doc/90001/90143/91125
func (provider *WxWorkProvider) GetLoginInfo(
	ctx context.Context,
	authCode string,
) (*LoginInfo, error) {
	result := &LoginInfo{}
	if err := provider.Client.HTTPPostJson(ctx, apiGetLoginInfo, map[string]string{
		"auth_code": authCode,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

type TrailInfo struct {
	utils.WeixinError
	/*
		license检查开启状态。
		0：未开启license检查状态（未迁移的历史授权的第三方应用（接入版本付费）或者未达到拦截时间的历史授权的的第三方应用（未接入版本付费）以及代开发应用）
		1：已开启license检查状态。若开启且已过试用期，则需要为企业购买license账号才可以使用
	*/
	LicenseStatus    int   `json:"license_status"`
	LicenseCheckTime int64 `json:"license_check_time"` // 接口开启拦截校验时间。开始拦截校验后，无接口许可将会被拦截，有接口许可将不会被拦截。
	TrailInfo        struct {
		StartTime int64 `json:"start_time"` // 接口许可试用开始时间
		// 若企业多次安装卸载同一个第三方应用，以第一次安装的时间为试用期开始时间，第一次安装完90天后为结束试用时间。
		EndTime int64 `json:"end_time"` // 接口许可试用到期时间。
	}
}

// 获取应用的接口许可状态
// https://developer.work.weixin.qq.com/document/path/95844
func (provider *WxWorkProvider) GetAppLicenseInfo(
	ctx context.Context, suiteid, corpid string,
) (*TrailInfo, error) {
	result := &TrailInfo{}
	if err := provider.Client.HTTPPostJson(ctx, apiGetAppLicenseInfo, map[string]interface{}{
		"corpid":   corpid,
		"suite_id": suiteid,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}
