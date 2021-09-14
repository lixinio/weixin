package authorizer

import (
	"fmt"

	"github.com/lixinio/weixin/utils"
)

// https://open.work.weixin.qq.com/api/doc/90001/90143/90605
type RefreshAccessToken func() (string, int, error) // 直接获取token， 不做任何缓存

// utils.AccessTokenGetter 接口实现
type authorizerAccessTokenGetterAdapter struct {
	accessTokenKey     string
	accessTokenLockKey string
	accessTokenGetter  RefreshAccessToken
}

// GetAccessToken 接口 utils.AccessTokenGetter 实现
func (adapter *authorizerAccessTokenGetterAdapter) GetAccessToken() (string, int, error) {
	return adapter.accessTokenGetter()
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
	suiteID, corpID string, agentID int,
	accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &authorizerAccessTokenGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"qywx.suite_agent_access_token.%s.%s.%d",
			suiteID, corpID, agentID,
		),
		accessTokenLockKey: fmt.Sprintf(
			"qywx.suite_agent_access_token.%s.%s.%d.lock",
			suiteID, corpID, agentID,
		),
	}
}
