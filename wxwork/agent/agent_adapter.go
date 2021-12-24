package agent

import (
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

var (
	ErrTokenUpdateForbidden = errors.New("can NOT refresh&update token in wxwork agent lite mode")
)

// 刷新 Access token
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
type RefreshAccessToken func() (string, int, error) // 直接获取token， 不做任何缓存

// utils.AccessTokenGetter 接口实现
type agentAccessTokenGetterAdapter struct {
	accessTokenKey     string
	accessTokenLockKey string
	accessTokenGetter  RefreshAccessToken
}

// GetAccessToken 接口 utils.AccessTokenGetter 实现
func (adapter *agentAccessTokenGetterAdapter) GetAccessToken() (string, int, error) {
	return adapter.accessTokenGetter()
}

// GetAccessTokenKey 接口 utils.AccessTokenGetter 实现
func (adapter *agentAccessTokenGetterAdapter) GetAccessTokenKey() string {
	return adapter.accessTokenKey
}

// GetAccessTokenLockKey 接口 utils.AccessTokenGetter 实现
func (adapter *agentAccessTokenGetterAdapter) GetAccessTokenLockKey() string {
	return adapter.accessTokenLockKey
}

func newAdapter(
	corpID string, agentID int, accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &agentAccessTokenGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"qywx_%s_%d.access_token", corpID, agentID,
		),
		accessTokenLockKey: fmt.Sprintf(
			"qywx_%s_%d.access_token.lock", corpID, agentID,
		),
	}
}
