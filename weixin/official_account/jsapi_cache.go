package official_account

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
func (adapter *jsApiTicketGetterAdapter) GetAccessToken() (string, int, error) {
	return adapter.accessTokenGetter()
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
	appid string, accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &jsApiTicketGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"weixin.jsapi_ticket.%s", appid,
		),
		accessTokenLockKey: fmt.Sprintf(
			"weixin.jsapi_ticket.%s.lock", appid,
		),
	}
}

func (officialAccount *OfficialAccount) EnableJSApiTicketCache(cache utils.Cache, locker utils.Lock) {
	if officialAccount.jsApiTicketCache != nil {
		return
	}

	officialAccount.jsApiTicketCache = utils.NewAccessTokenCache(
		newJsApiTicketAdapter(officialAccount.Config.Appid, func() (string, int, error) {
			ticket, expiresIn, err := officialAccount.getJSApiTicket(context.TODO())
			return ticket, int(expiresIn), err
		}),
		cache, locker,
	)
}

// utils.AccessTokenGetter 接口实现
type wxCardTicketGetterAdapter struct {
	accessTokenKey     string
	accessTokenLockKey string
	accessTokenGetter  RefreshAccessToken
}

// GetAccessToken 接口 utils.AccessTokenGetter 实现
func (adapter *wxCardTicketGetterAdapter) GetAccessToken() (string, int, error) {
	return adapter.accessTokenGetter()
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
	appid string, accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &wxCardTicketGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"weixin.wx_card_ticket.%s", appid,
		),
		accessTokenLockKey: fmt.Sprintf(
			"weixin.wx_card_ticket.%s.lock", appid,
		),
	}
}

func (officialAccount *OfficialAccount) EnableWxCardTicketCache(cache utils.Cache, locker utils.Lock) {
	if officialAccount.wxCardTicketCache != nil {
		return
	}

	officialAccount.wxCardTicketCache = utils.NewAccessTokenCache(
		newWxCardTicketAdapter(officialAccount.Config.Appid, func() (string, int, error) {
			ticket, expiresIn, err := officialAccount.getWxCardApiTicket(context.TODO())
			return ticket, int(expiresIn), err
		}),
		cache, locker,
	)
}
