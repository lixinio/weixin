package authorizer

import (
	"context"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

// 需要通过wxopen对象刷新authorizer Access token
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/api_authorizer_token.html
type RefreshAccessToken func(context.Context) (string, int, error) // 直接获取token， 不做任何缓存

// utils.AccessTokenGetter 接口实现
type authorizerAccessTokenGetterAdapter struct {
	accessTokenKey     string
	accessTokenLockKey string
	accessTokenGetter  RefreshAccessToken
}

// GetAccessToken 接口 utils.AccessTokenGetter 实现
func (adapter *authorizerAccessTokenGetterAdapter) GetAccessToken(
	ctx context.Context,
) (string, int, error) {
	return adapter.accessTokenGetter(ctx)
}

// GetAccessTokenKey 接口 utils.AccessTokenGetter 实现
func (adapter *authorizerAccessTokenGetterAdapter) GetAccessTokenKey() string {
	return adapter.accessTokenKey
}

// GetAccessTokenLockKey 接口 utils.AccessTokenGetter 实现
func (adapter *authorizerAccessTokenGetterAdapter) GetAccessTokenLockKey() string {
	return adapter.accessTokenLockKey
}

func newAdapter(
	componentAppid, appid string,
	accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &authorizerAccessTokenGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"weixin.authorizer_access_token.%s.%s",
			componentAppid, appid,
		),
		accessTokenLockKey: fmt.Sprintf(
			"weixin.authorizer_access_token.%s.%s.lock",
			componentAppid, appid,
		),
	}
}
