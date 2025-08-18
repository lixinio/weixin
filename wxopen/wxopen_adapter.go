package wxopen

import (
	"context"
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

// ticket 缓存的adapter， 无法自行刷新，只能微信主动上报
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/component_verify_ticket.html
type ticketAdaptor struct {
	tokenKey  string
	lockerKey string
}

func (ta *ticketAdaptor) GetAccessToken(
	ctx context.Context,
) (accessToken string, expiresIn int, err error) {
	return "", 0, errors.New("can NOT update wxopen ticket")
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (ta *ticketAdaptor) GetAccessTokenKey() string {
	return ta.tokenKey
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (ta *ticketAdaptor) GetAccessTokenLockKey() string {
	return ta.lockerKey
}

func newTicketAdapter(appid string) *ticketAdaptor {
	return &ticketAdaptor{
		tokenKey:  fmt.Sprintf("weixin.component.wxopen_ticket.%s", appid),
		lockerKey: fmt.Sprintf("weixin.component.wxopen_ticket.%s.lock", appid),
	}
}

// access token的cache
type accessTokenAdaptor struct {
	config      *Config
	ticketCache *utils.AccessTokenCache
	client      *utils.Client
	tokenKey    string
	lockerKey   string
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (ta *accessTokenAdaptor) GetAccessTokenKey() string {
	return ta.tokenKey
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (ta *accessTokenAdaptor) GetAccessTokenLockKey() string {
	return ta.lockerKey
}

func (ta *accessTokenAdaptor) GetAccessToken(
	ctx context.Context,
) (accessToken string, expiresIn int, err error) {
	if ta.ticketCache == nil {
		return "", 0, fmt.Errorf(
			"wxopen appid : %s, error: %w", ta.config.Appid, ErrTokenUpdateForbidden,
		)
	}

	ticket, err := ta.ticketCache.GetAccessToken(ctx)
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
		"component_appid":         ta.config.Appid,
		"component_appsecret":     ta.config.Secret,
		"component_verify_ticket": ticket,
	}
	if err := ta.client.HTTPPostToken(
		ctx, apiGetComponentToken, payload, &result,
	); err != nil {
		return "", 0, err
	}
	return result.AccessToken, result.ExpiresIn, nil
}

func newAccessTokenAdaptor(
	config *Config,
	ticketCache *utils.AccessTokenCache,
) *accessTokenAdaptor {
	adaptor := &accessTokenAdaptor{
		config:      config,
		ticketCache: ticketCache,
		tokenKey:    fmt.Sprintf("weixin.component.access_token.%s", config.Appid),
		lockerKey:   fmt.Sprintf("weixin.component.access_token.%s.lock", config.Appid),
	}
	if ticketCache != nil {
		adaptor.client = utils.NewClient(WXServerUrl, utils.EmptyClientAccessTokenGetter(0))
	}
	return adaptor
}
