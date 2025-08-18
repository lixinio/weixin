package authorizer

import (
	"context"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

// https://open.work.weixin.qq.com/api/doc/90001/90143/90605
type RefreshAccessToken func(ctx context.Context) (string, int, error) // 直接获取token， 不做任何缓存

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
	suiteID, corpID string, _ int,
	accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &authorizerAccessTokenGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"qywx.suite_agent_access_token.%s.%s", suiteID, corpID,
		),
		accessTokenLockKey: fmt.Sprintf(
			"qywx.suite_agent_access_token.%s.%s.lock", suiteID, corpID,
		),
	}
}
