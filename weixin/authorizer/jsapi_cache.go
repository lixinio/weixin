package authorizer

import (
	"context"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

// utils.AccessTokenGetter 接口实现
type jsApiTicketGetterAdapter struct {
	accessTokenKey     string
	accessTokenLockKey string
	accessTokenGetter  RefreshAccessToken
}

// GetAccessToken 接口 utils.AccessTokenGetter 实现
func (adapter *jsApiTicketGetterAdapter) GetAccessToken(
	ctx context.Context,
) (string, int, error) {
	return adapter.accessTokenGetter(ctx)
}

// GetAccessTokenKey 接口 utils.AccessTokenGetter 实现
func (adapter *jsApiTicketGetterAdapter) GetAccessTokenKey() string {
	return adapter.accessTokenKey
}

// GetAccessTokenLockKey 接口 utils.AccessTokenGetter 实现
func (adapter *jsApiTicketGetterAdapter) GetAccessTokenLockKey() string {
	return adapter.accessTokenLockKey
}

func newJsApiTicketAdapter(
	componentAppid, appid string, accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &jsApiTicketGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"weixin.authorizer_jsapi_ticket.%s.%s", componentAppid, appid,
		),
		accessTokenLockKey: fmt.Sprintf(
			"weixin.authorizer_jsapi_ticket.%s.%s.lock", componentAppid, appid,
		),
	}
}

func (api *Authorizer) EnableJSApiTicketCache(
	cache utils.Cache, locker utils.Lock,
	tokenRefreshHandler utils.TokenRefreshHandler, // 刷新callback
) {
	if api.jsApiTicketCache != nil {
		return
	}

	api.jsApiTicketCache = utils.NewAccessTokenCache(
		newJsApiTicketAdapter(
			api.ComponentAppid,
			api.Appid,
			func(ctx context.Context) (string, int, error) {
				ticket, expiresIn, err := api.getJSApiTicket(ctx)
				return ticket, int(expiresIn), err
			}),
		cache, locker,
		utils.CacheClientTokenOptWithExpireBefore(tokenRefreshHandler),
	)
}

// utils.AccessTokenGetter 接口实现
type wxCardTicketGetterAdapter struct {
	accessTokenKey     string
	accessTokenLockKey string
	accessTokenGetter  RefreshAccessToken
}

// GetAccessToken 接口 utils.AccessTokenGetter 实现
func (adapter *wxCardTicketGetterAdapter) GetAccessToken(
	ctx context.Context,
) (string, int, error) {
	return adapter.accessTokenGetter(ctx)
}

// GetAccessTokenKey 接口 utils.AccessTokenGetter 实现
func (adapter *wxCardTicketGetterAdapter) GetAccessTokenKey() string {
	return adapter.accessTokenKey
}

// GetAccessTokenLockKey 接口 utils.AccessTokenGetter 实现
func (adapter *wxCardTicketGetterAdapter) GetAccessTokenLockKey() string {
	return adapter.accessTokenLockKey
}

func newWxCardTicketAdapter(
	componentAppid, appid string, accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &wxCardTicketGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"weixin.authorizer_wx_card_ticket.%s.%s", componentAppid, appid,
		),
		accessTokenLockKey: fmt.Sprintf(
			"weixin.authorizer_wx_card_ticket.%s.%s.lock", componentAppid, appid,
		),
	}
}

func (api *Authorizer) EnableWxCardTicketCache(
	cache utils.Cache, locker utils.Lock,
	tokenRefreshHandler utils.TokenRefreshHandler, // 刷新callback
) {
	if api.wxCardTicketCache != nil {
		return
	}

	api.wxCardTicketCache = utils.NewAccessTokenCache(
		newWxCardTicketAdapter(
			api.ComponentAppid,
			api.Appid,
			func(ctx context.Context) (string, int, error) {
				ticket, expiresIn, err := api.getWxCardApiTicket(ctx)
				return ticket, int(expiresIn), err
			}),
		cache, locker,
		utils.CacheClientTokenOptWithExpireBefore(tokenRefreshHandler),
	)
}
