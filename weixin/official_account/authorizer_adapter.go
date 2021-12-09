package official_account

import (
	"fmt"

	"github.com/lixinio/weixin/utils"
)

// 刷新 Access token
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
type RefreshAccessToken func() (string, int, error) // 直接获取token， 不做任何缓存

// utils.AccessTokenGetter 接口实现
type oaAccessTokenGetterAdapter struct {
	accessTokenKey     string
	accessTokenLockKey string
	accessTokenGetter  RefreshAccessToken
}

// GetAccessToken 接口 utils.AccessTokenGetter 实现
func (adapter *oaAccessTokenGetterAdapter) GetAccessToken() (string, int, error) {
	return adapter.accessTokenGetter()
}

// GetAccessTokenKey 接口 utils.AccessTokenGetter 实现
func (adapter *oaAccessTokenGetterAdapter) GetAccessTokenKey() string {
	return adapter.accessTokenKey
}

// GetAccessTokenLockKey 接口 utils.AccessTokenGetter 实现
func (adapter *oaAccessTokenGetterAdapter) GetAccessTokenLockKey() string {
	return adapter.accessTokenLockKey
}

func newAdapter(
	appid string, accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &oaAccessTokenGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"weixin.access_token.%s", appid,
		),
		accessTokenLockKey: fmt.Sprintf(
			"weixin.access_token.%s.lock", appid,
		),
	}
}
