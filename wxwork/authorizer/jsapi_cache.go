package authorizer

import (
	"context"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

// utils.AccessTokenGetter 接口实现
type corpJsApiTicketGetterAdapter struct {
	accessTokenKey     string
	accessTokenLockKey string
	accessTokenGetter  RefreshAccessToken
}

// GetAccessToken 接口 utils.AccessTokenGetter 实现
func (adapter *corpJsApiTicketGetterAdapter) GetAccessToken(
	ctx context.Context,
) (string, int, error) {
	return adapter.accessTokenGetter(ctx)
}

// GetAccessTokenKey 接口 utils.AccessTokenGetter 实现
func (adapter *corpJsApiTicketGetterAdapter) GetAccessTokenKey() string {
	return adapter.accessTokenKey
}

// GetAccessTokenLockKey 接口 utils.AccessTokenGetter 实现
func (adapter *corpJsApiTicketGetterAdapter) GetAccessTokenLockKey() string {
	return adapter.accessTokenLockKey
}

func newCorpJsApiTicketAdapter(
	suitid, corpid string, accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &corpJsApiTicketGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"qywx.suite_agent_jsapi_ticket.%s.%s", suitid, corpid,
		),
		accessTokenLockKey: fmt.Sprintf(
			"qywx.suite_agent_jsapi_ticket.%s.%s.lock", suitid, corpid,
		),
	}
}

func (authorizer *Authorizer) EnableCorpJSApiTicketCache(
	cache utils.Cache,
	locker utils.Lock,
	tokenRefreshHandler utils.TokenRefreshHandler, // 刷新callback
) {
	if authorizer.corpJsApiTicketCache != nil {
		return
	}

	authorizer.corpJsApiTicketCache = utils.NewAccessTokenCache(
		newCorpJsApiTicketAdapter(
			authorizer.SuiteID,
			authorizer.CorpID,
			func(ctx context.Context) (string, int, error) {
				ticket, expiresIn, err := authorizer.getCorpJSApiTicket(ctx)
				return ticket, int(expiresIn), err
			}),
		cache, locker,
		utils.CacheClientTokenOptWithExpireBefore(tokenRefreshHandler),
	)
}

// utils.AccessTokenGetter 接口实现
type agentJsApiTicketGetterAdapter struct {
	accessTokenKey     string
	accessTokenLockKey string
	accessTokenGetter  RefreshAccessToken
}

// GetAccessToken 接口 utils.AccessTokenGetter 实现
func (adapter *agentJsApiTicketGetterAdapter) GetAccessToken(
	ctx context.Context,
) (string, int, error) {
	return adapter.accessTokenGetter(ctx)
}

// GetAccessTokenKey 接口 utils.AccessTokenGetter 实现
func (adapter *agentJsApiTicketGetterAdapter) GetAccessTokenKey() string {
	return adapter.accessTokenKey
}

// GetAccessTokenLockKey 接口 utils.AccessTokenGetter 实现
func (adapter *agentJsApiTicketGetterAdapter) GetAccessTokenLockKey() string {
	return adapter.accessTokenLockKey
}

func newAgentJsApiTicketAdapter(
	suitid, corpid string, accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &agentJsApiTicketGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"qywx.suite_agent_jsapi_agent_ticket.%s.%s", suitid, corpid,
		),
		accessTokenLockKey: fmt.Sprintf(
			"qywx.suite_agent_jsapi_agent_ticket.%s.%s.lock", suitid, corpid,
		),
	}
}

func (authorizer *Authorizer) EnableAgentJSApiTicketCache(
	cache utils.Cache,
	locker utils.Lock,
	tokenRefreshHandler utils.TokenRefreshHandler, // 刷新callback
) {
	if authorizer.agentJsApiTicketCache != nil {
		return
	}

	authorizer.agentJsApiTicketCache = utils.NewAccessTokenCache(
		newAgentJsApiTicketAdapter(
			authorizer.SuiteID,
			authorizer.CorpID,
			func(ctx context.Context) (string, int, error) {
				ticket, expiresIn, err := authorizer.getAgentJSApiTicket(ctx)
				return ticket, int(expiresIn), err
			}),
		cache, locker,
		utils.CacheClientTokenOptWithExpireBefore(tokenRefreshHandler),
	)
}
