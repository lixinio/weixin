package wxwork_provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	WXServerUrl         = "https://qyapi.weixin.qq.com" // 微信 api 服务器地址
	apiAuthorize        = "https://open.work.weixin.qq.com/wwopen/sso/3rd_qrConnect"
	apiGetProviderToken = "/cgi-bin/service/get_provider_token"
	apiGetLoginInfo     = "/cgi-bin/service/get_login_info"
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

func New(cache utils.Cache, locker utils.Lock, config *Config) *WxWorkProvider {
	accessTokenCache := utils.NewAccessTokenCache(
		newAccessTokenAdaptor(config), cache, locker,
	)
	instance := &WxWorkProvider{
		Config:           config,
		Client:           utils.NewClient(WXServerUrl, accessTokenCache),
		accessTokenCache: accessTokenCache,
	}
	return instance
}

func NewLite(cache utils.Cache, locker utils.Lock, corpID string) *WxWorkProvider {
	return New(cache, locker, &Config{CorpID: corpID})
}

func (provider *WxWorkProvider) RefreshAccessToken(expireBefore int) (string, error) {
	if provider.Config.ProviderSecret == "" {
		return "", fmt.Errorf(
			"provider corpid : %s, error: %w", provider.Config.CorpID, ErrTokenUpdateForbidden,
		)
	}
	return provider.accessTokenCache.RefreshAccessToken(expireBefore)
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
