package wxopen

import (
	"context"
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

var (
	ErrTicketUpdateForbidden = errors.New("can NOT update ticket in wxopen lite mode")
	ErrTokenUpdateForbidden  = errors.New("can NOT refresh&update token in wxopen lite mode")
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
	Config           *Config
	Client           *utils.Client
	ticketCache      *utils.AccessTokenCache
	accessTokenCache *utils.AccessTokenCache
}

func New(cache utils.Cache, locker utils.Lock, config *Config) *WxOpen {
	ticketCache := utils.NewAccessTokenCache(newTicketAdapter(config.Appid), cache, locker)
	accessTokenCache := utils.NewAccessTokenCache(
		newAccessTokenAdaptor(config, ticketCache), cache, locker,
	)
	instance := &WxOpen{
		Config:           config,
		Client:           utils.NewClient(WXServerUrl, accessTokenCache),
		ticketCache:      ticketCache,
		accessTokenCache: accessTokenCache,
	}
	instance.Client.UpdateAccessTokenKey(accessTokenKey) // token的名称不一样
	return instance
}

func NewLite(
	cache utils.Cache,
	locker utils.Lock,
	appID string,
) *WxOpen {
	config := &Config{Appid: appID}
	instance := &WxOpen{
		Config: config,
		Client: utils.NewClient(WXServerUrl, utils.NewAccessTokenCache(
			newAccessTokenAdaptor(config, nil), cache, locker,
		)),
	}
	instance.Client.UpdateAccessTokenKey(accessTokenKey) // token的名称不一样
	return instance
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
		return fmt.Errorf(
			"wxopen appid : %s, error: %w", wxopen.Config.Appid, ErrTicketUpdateForbidden,
		)
	}
	_, err := wxopen.ticketCache.UpdateAccessToken(token, ticketExpiresIn)
	return err
}

func (wxopen *WxOpen) RefreshAccessToken(expireBefore int) (string, error) {
	if wxopen.accessTokenCache == nil {
		return "", fmt.Errorf(
			"wxopen appid : %s, error: %w", wxopen.Config.Appid, ErrTokenUpdateForbidden,
		)
	}
	return wxopen.accessTokenCache.RefreshAccessToken(expireBefore)
}
