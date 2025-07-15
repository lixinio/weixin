package wxwork_suite

import (
	"context"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

// ticket 缓存的adapter， 无法自行刷新，只能企业微信主动上报
// 推送suite_ticket
// https://open.work.weixin.qq.com/api/doc/90001/90143/90628
type ticketAdaptor struct {
	tokenKey  string
	lockerKey string
}

func (ta *ticketAdaptor) GetAccessToken(
	ctx context.Context,
) (accessToken string, expiresIn int, err error) {
	return "", 0, fmt.Errorf("can NOT update wx suite ticket: %w", ErrTicketUpdateForbidden)
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (ta *ticketAdaptor) GetAccessTokenKey() string {
	return ta.tokenKey
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (ta *ticketAdaptor) GetAccessTokenLockKey() string {
	return ta.lockerKey
}

func newTicketAdapter(suiteID string) *ticketAdaptor {
	return &ticketAdaptor{
		tokenKey:  fmt.Sprintf("qywx.suite_access_ticket.%s", suiteID),
		lockerKey: fmt.Sprintf("qywx.suite_access_ticket.%s.lock", suiteID),
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

// 获取第三方应用凭证
// https://open.work.weixin.qq.com/api/doc/90001/90143/90600
func (ta *accessTokenAdaptor) GetAccessToken(
	ctx context.Context,
) (accessToken string, expiresIn int, err error) {
	if ta.ticketCache == nil {
		return "", 0, fmt.Errorf(
			"wxopen appid : %s, error: %w", ta.config.SuiteID, ErrTokenUpdateForbidden,
		)
	}

	ticket, err := ta.ticketCache.GetAccessToken(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("can NOT get suite access token without ticket, %w", err)
	}

	// AccessToken 和其他地方 字段不一致
	result := struct {
		utils.WeixinError
		AccessToken string `json:"suite_access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}{}

	payload := map[string]string{
		"suite_id":     ta.config.SuiteID,
		"suite_secret": ta.config.SuiteSecret,
		"suite_ticket": ticket,
	}
	if err := ta.client.HTTPPostToken(
		ctx, apiGetSuiteToken, payload, &result,
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
		tokenKey:    fmt.Sprintf("qywx.suite_access_token.%s", config.SuiteID),
		lockerKey:   fmt.Sprintf("qywx.suite_access_token.%s.lock", config.SuiteID),
	}
	if ticketCache != nil {
		adaptor.client = utils.NewClient(WXServerUrl, utils.EmptyClientAccessTokenGetter(0))
	}
	return adaptor
}
