package wxopen

import (
	"context"
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

var (
	ErrTokenUpdateForbidden = errors.New("can NOT refresh&update token in wxopen lite mode")
)

// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Before_Develop/creat_token.html
// 出于安全考虑，在第三方平台创建审核通过后，微信服务器 每隔 10 分钟会向第三方的消息接收地址推送一次 component_verify_ticket，
// 用于获取第三方平台接口调用凭据。component_verify_ticket有效期为12h
const ticketExpiresIn = 43000 // 12h (43200)

const (
	WXServerUrl          = "https://api.weixin.qq.com" // 微信 api 服务器地址
	accessTokenKey       = "component_access_token"
	apiGetComponentToken = "/cgi-bin/component/api_component_token"
	apiStartPushTicket   = "/cgi-bin/component/api_start_push_ticket"
)

// ticket 缓存的adapter， 无法自行刷新，只能微信主动上报
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/component_verify_ticket.html
type ticketAdaptor struct {
	appid string
}

func (ta *ticketAdaptor) GetAccessToken() (accessToken string, expiresIn int, err error) {
	return "", 0, errors.New("can NOT update wxopen ticket")
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (ta *ticketAdaptor) GetAccessTokenKey() string {
	return fmt.Sprintf("access-token:wxopen_ticket:%s", ta.appid)
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (ta *ticketAdaptor) GetAccessTokenLockKey() string {
	return fmt.Sprintf("access-token:wxopen_ticket:%s.lock", ta.appid)
}

/*
开放平台配置
*/
type Config struct {
	Appid          string
	Secret         string
	Token          string
	EncodingAESKey string
}

type WxOpen struct {
	Config      *Config
	Client      *utils.Client
	ticketCache *utils.AccessTokenCache
}

func New(cache utils.Cache, locker utils.Lock, config *Config) *WxOpen {
	ticketCache := utils.NewAccessTokenCache(&ticketAdaptor{config.Appid}, cache, locker, 0)
	instance := &WxOpen{
		Config:      config,
		ticketCache: ticketCache,
	}

	client := utils.NewClient(WXServerUrl, utils.NewAccessTokenCache(instance, cache, locker, 0))
	client.UpdateAccessTokenKey(accessTokenKey) // token的名称不一样
	instance.Client = client

	return instance
}

func NewLite(
	cache utils.Cache,
	locker utils.Lock,
	appID string,
) *WxOpen {
	instance := &WxOpen{
		Config:      &Config{Appid: appID},
		ticketCache: nil,
	}
	client := utils.NewClient(WXServerUrl, utils.NewAccessTokenCache(instance, cache, locker, 0))
	client.UpdateAccessTokenKey(accessTokenKey) // token的名称不一样
	instance.Client = client
	return instance
}

// GetAccessToken 接口 weixin.AccessTokenGetter 实现
func (wxopen *WxOpen) GetAccessToken() (accessToken string, expiresIn int, err error) {
	accessToken, expiresIn, err = wxopen.refreshAccessTokenFromWXServer()
	return
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (wxopen *WxOpen) GetAccessTokenKey() string {
	return fmt.Sprintf("access-token:wxopen:%s", wxopen.Config.Appid)
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (wxopen *WxOpen) GetAccessTokenLockKey() string {
	return fmt.Sprintf("access-token:wxopen:%s.lock", wxopen.Config.Appid)
}

/*
从微信服务器获取新的 AccessToken
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/component_access_token.html
*/
func (wxopen *WxOpen) refreshAccessTokenFromWXServer() (accessToken string, expiresIn int, err error) {
	if wxopen.ticketCache == nil {
		return "", 0, fmt.Errorf(
			"wxopen appid : %s, error: %w", wxopen.Config.Appid, ErrTokenUpdateForbidden,
		)
	}

	ticket, err := wxopen.ticketCache.GetAccessToken()
	if err != nil {
		return "", 0, fmt.Errorf("can NOT get wxopen access token without ticket, %w", err)
	}

	// AccessToken 和其他地方 字段不一致
	result := struct {
		utils.WeixinError
		AccessToken string `json:"component_access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}{}

	payload := map[string]string{
		"component_appid":         wxopen.Config.Appid,
		"component_appsecret":     wxopen.Config.Secret,
		"component_verify_ticket": ticket,
	}
	if err := wxopen.Client.HTTPPostToken(
		context.TODO(), apiGetComponentToken, payload, &result,
	); err != nil {
		return "", 0, err
	}
	return result.AccessToken, result.ExpiresIn, nil
}

/*
启动ticket推送服务
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/component_verify_ticket_service.html
*/
func (wxopen *WxOpen) StartPushTicket(ctx context.Context) error {
	payload := map[string]string{
		"component_appid":  wxopen.Config.Appid,
		"component_secret": wxopen.Config.Secret,
	}

	return wxopen.Client.HTTPPostToken(
		context.TODO(), apiStartPushTicket, payload, nil,
	)
}

// 当收到EventComponentVerifyTicket时， 用于更新ticket到cache
func (wxopen *WxOpen) UpdateTicket(token string) error {
	if wxopen.ticketCache == nil {
		return fmt.Errorf("wxopen appid : %s, error: %w", wxopen.Config.Appid, ErrTokenUpdateForbidden)
	}
	_, err := wxopen.ticketCache.UpdateAccessToken(token, ticketExpiresIn)
	return err
}
